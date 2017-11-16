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

func (builder *cmdBuilderImpl) IsCmd(resource v1.APIResource) bool {
	if !strings.Contains(resource.Name, "/") {
		return false
	}
	if strings.HasSuffix(resource.Name, "/status") {
		return false
	}
	gvk := schema.GroupVersionKind{resource.Group, resource.Version, resource.Kind}
	if builder.resources.LookupResource(gvk) == nil {
		return false
	}
	return true
}

func (builder *cmdBuilderImpl) Seen(resource v1.APIResource) bool {
	parts := strings.Split(resource.Name, "/")
	kind := parts[0]
	operation := parts[1]

	return builder.seen[operation].HasAny(kind)
}

func (builder *cmdBuilderImpl) operation(resource v1.APIResource) string {
	parts := strings.Split(resource.Name, "/")
	return parts[1]
}

func (builder *cmdBuilderImpl) resource(resource v1.APIResource) string {
	parts := strings.Split(resource.Name, "/")
	return parts[0]
}

func (builder *cmdBuilderImpl) add(resource v1.APIResource) {
	parts := strings.Split(resource.Name, "/")
	kind := parts[0]
	operation := parts[1]
	if _, found := builder.seen[operation]; !found {
		builder.seen[operation] = sets.String{}
	}

	builder.seen[operation].Insert(kind)
}
