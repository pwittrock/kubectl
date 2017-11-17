/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openapi

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

func (builder *cmdBuilderImpl) isSubResource(resource *v1.APIResource) bool {
	if !strings.Contains(resource.Name, "/") {
		return false
	}
	if strings.HasSuffix(resource.Name, "/status") {
		return false
	}
	return true
}

func (builder *cmdBuilderImpl) isCmd(resource *v1.APIResource) bool {
	gvk := schema.GroupVersionKind{resource.Group, resource.Version, resource.Kind}
	return builder.resources.LookupResource(gvk) != nil
}

func (builder *cmdBuilderImpl) isResource(resource *v1.APIResource) bool {
	return !strings.Contains(resource.Name, "/")
}

func (builder *cmdBuilderImpl) setGroupVersionFromParentIfMissing(resource *v1.APIResource, parent *v1.APIResource) {
	if len(resource.Group) == 0 {
		resource.Group = parent.Group
	}
	if len(resource.Version) == 0 {
		resource.Version = parent.Version
	}
}

func (builder *cmdBuilderImpl) done(resource *v1.APIResource) bool {
	parts := strings.Split(resource.Name, "/")
	kind := parts[0]
	operation := parts[1]

	return builder.seen[operation].HasAny(kind)
}

func (builder *cmdBuilderImpl) operation(resource *v1.APIResource) string {
	parts := strings.Split(resource.Name, "/")
	return parts[1]
}

func (builder *cmdBuilderImpl) resource(resource *v1.APIResource) string {
	parts := strings.Split(resource.Name, "/")
	return parts[0]
}

func (builder *cmdBuilderImpl) add(resource *v1.APIResource) {
	parts := strings.Split(resource.Name, "/")
	kind := parts[0]
	operation := parts[1]
	if _, found := builder.seen[operation]; !found {
		builder.seen[operation] = sets.String{}
	}

	builder.seen[operation].Insert(kind)
}
