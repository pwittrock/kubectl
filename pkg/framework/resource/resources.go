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
	"sort"
)

// Resources is the set of resources found in the API server
type Resources map[string][]*Resource

// Sort returns the resources sorted alphanumberically by their names
func (r Resources) SortKeys() []string {
	ordered := []string{}
	for resource, _ := range r {
		ordered = append(ordered, resource)
	}
	sort.Strings(ordered)
	return ordered
}

// Filter filters resources and subresources
func (r Resources) Filter(filter Filter) Resources {
	value := Resources{}
	for resource, versions := range r {
		for _, version := range versions {
			if !filter.Resource(version) {
				continue
			}

			copy := r.filterSubresources(*version, filter)
			value[resource] = append(value[resource], &copy)
		}
	}

	return value
}

// filterSubresources returns a copy of resource with the subresources filtered
func (r Resources) filterSubresources(resource Resource, filter Filter) Resource {
	original := resource.SubResources
	resource.SubResources = nil
	for _, subresource := range original {
		if !filter.SubResource(subresource) {
			continue
		}
		resource.SubResources = append(resource.SubResources, subresource)
	}
	return resource
}
