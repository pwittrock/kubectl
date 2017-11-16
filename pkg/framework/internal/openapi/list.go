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

import "k8s.io/apimachinery/pkg/apis/meta/v1"

func (builder *cmdBuilderImpl) ListResources() ([]v1.APIResource, error) {
	list := []v1.APIResource{}
	gvs, err := builder.discovery.ServerResources()
	if err != nil {
		return nil, err
	}
	for _, gv := range gvs {
		for _, r := range gv.APIResources {
			list = append(list, r)
		}
	}
	return list, nil
}
