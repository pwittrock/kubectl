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

	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	if !resource.Namespaced {
		return false
	}
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
	if len(resource.Group) == 0 {
		resource.Group = "core"
	}
	if len(resource.Version) == 0 {
		resource.Version = parent.Version
	}
}

func (builder *cmdBuilderImpl) operation(resource *v1.APIResource) string {
	parts := strings.Split(resource.Name, "/")
	return parts[1]
}

func (builder *cmdBuilderImpl) resource(resource *v1.APIResource) string {
	parts := strings.Split(resource.Name, "/")
	return parts[0]
}

func (builder *cmdBuilderImpl) groupVersion(groupVersion string) (string, string) {
	parts := strings.Split(groupVersion, "/")
	var group, version string

	// Group maybe missing for apis under the "core" group
	if len(parts) > 1 {
		group = parts[0]
	} else {
		group = "core"
	}

	if len(parts) > 1 {
		version = parts[1]
	} else if len(parts) > 0 {
		version = parts[0]
	} else {
		version = "v1"
	}

	return group, version
}

func (builder *cmdBuilderImpl) resourceOperation(name string) (string, string, error) {
	parts := strings.Split(name, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("%s doesn't not match subresource name format", name)
	}
	return parts[0], parts[1], nil
}
