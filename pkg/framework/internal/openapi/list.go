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
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func getGVR(group, version, resource string) schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
}

func (builder *cmdBuilderImpl) getSubResources() (map[string][]*SubResource, error) {
	list := map[string][]*SubResource{}
	gvs, err := builder.discovery.ServerResources()
	if err != nil {
		return nil, err
	}

	parentResources := map[schema.GroupVersionResource]v1.APIResource{}

	// Index the parent resources by group,version,resource
	for _, gv := range gvs {
		group, version := builder.groupVersion(gv.GroupVersion)
		for _, r := range gv.APIResources {
			if builder.isResource(&r) {
				parentResources[schema.GroupVersionResource{
					Group:    group,
					Version:  version,
					Resource: r.Name,
				}] = r
			}
		}
	}

	// Map subresources to resource names
	for _, gv := range gvs {
		group, version := builder.groupVersion(gv.GroupVersion)
		if len(builder.apiGroup) > 0 && builder.apiGroup != group {
			continue
		}

		if len(builder.apiVersion) > 0 && builder.apiVersion != version {
			continue
		}

		for _, r := range gv.APIResources {
			if !builder.isSubResource(&r) {
				continue
			}

			// Sanity check - this shouldn't happen in practice
			if _, found := parentResources[getGVR(group, version, builder.resource(&r))]; !found {
				return nil, fmt.Errorf("Missing parent for subresource %s", r.Name)
			}

			// Set the group and version to the API groupVersion if missing
			if len(r.Group) == 0 {
				r.Group = group
			}
			if len(r.Version) == 0 {
				r.Version = version
			}

			gvk := schema.GroupVersionKind{
				r.Group,
				r.Version,
				r.Kind,
			}
			openapiSchema := builder.resources.LookupResource(gvk)
			if openapiSchema == nil {
				continue
			}

			// reassign variable so we can get a pointer to it
			sub := &SubResource{
				resource:                 r,
				resourceGroupVersionKind: schema.GroupVersionKind{r.Group, r.Version, r.Kind},
				parent:          parentResources[getGVR(group, version, builder.resource(&r))],
				apiGroupVersion: schema.GroupVersion{Group: group, Version: version},
				openapiSchema:   openapiSchema,
			}
			list[r.Name] = append(list[r.Name], sub)
		}
	}
	return list, nil
}
