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

package resource

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

type Resource struct {
	resource         v1.APIResource
	groupVersionKind schema.GroupVersionKind
	openapiSchema    openapi.Schema
}

type SubResource struct {
	resource                 v1.APIResource
	resourceGroupVersionKind schema.GroupVersionKind
	parent                   v1.APIResource
	apiGroupVersion          schema.GroupVersion
	openapiSchema            openapi.Schema
}

type Parser struct {
	resources  openapi.Resources
	discovery  discovery.DiscoveryInterface
	rest       rest.Interface
	apiGroup   string
	apiVersion string
}

func (*Parser) Resources() []*Resource {
	return nil
}

func (*Parser) SubResources() []*SubResource {
	return nil
}

// SubResource returns a resource name, subresource name pair and true if resource is subresource.
// Returns a resource name, empty string and false if resource is not a subresource.
func (*Parser) SubResource(resource *v1.APIResource) (string, string, bool) {
	parts := strings.Split(resource.Name, "/")
	if len(parts) > 1 {
		return parts[0], parts[1], true
	}

	return parts[0], "", false
}

// CopyGroupVersion copies the group and version from src to dest if either is missing from dest
// If the src group is empty and the dest group is empty, sets the dest group to "core"
func (*Parser) CopyGroupVersion(src *v1.APIResource, dest *v1.APIResource) {
	if len(dest.Group) == 0 {
		dest.Group = src.Group
	}
	if len(dest.Group) == 0 {
		dest.Group = "core"
	}
	if len(dest.Version) == 0 {
		dest.Version = src.Version
	}
}

// SplitGroupVersion splits the groupVersion string into the group and version components
func (*Parser) SplitGroupVersion(groupVersion string) (string, string) {
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
