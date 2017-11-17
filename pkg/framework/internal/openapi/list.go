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
)

func (builder *cmdBuilderImpl) listResources() ([]*v1.APIResource, error) {
	list := []*v1.APIResource{}
	gvs, err := builder.discovery.ServerResources()
	if err != nil {
		return nil, err
	}
	for _, gv := range gvs {
		for _, r := range gv.APIResources {
			// reassign variable so we can get a pointer to it
			resource := r

			list = append(list, &resource)

			// Set the group and version on the resource from the API groupversion if it is missing
			parts := strings.Split(gv.GroupVersion, "/")

			// Group maybe missing for apis under the "core" group
			if len(resource.Group) == 0 && len(parts) > 1 {
				resource.Group = parts[0]
			} else if len(resource.Group) == 0 {
				resource.Group = "core"
			}

			if len(resource.Version) == 0 && len(parts) > 1 {
				resource.Version = parts[1]
			} else if len(resource.Version) == 0 && len(parts) > 0 {
				resource.Version = parts[0]
			}

		}
	}
	return list, nil
}
