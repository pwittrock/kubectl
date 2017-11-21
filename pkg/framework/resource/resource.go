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
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

// Resource is an API Resource
type Resource struct {
	Resource        v1.APIResource
	ApiGroupVersion schema.GroupVersion
	OpenapiSchema   openapi.Schema
	SubResources    []*SubResource
}

func (r *Resource) HasField(fieldPath, fieldType string) bool {
	return hasField(r.OpenapiSchema, fieldPath, fieldType)
}

// SubResource is an API subresource
type SubResource struct {
	Resource        v1.APIResource
	Parent          *Resource
	ApiGroupVersion schema.GroupVersion
	OpenapiSchema   openapi.Schema
}

func (sr *SubResource) HasField(fieldPath, fieldType string) bool {
	return hasField(sr.OpenapiSchema, fieldPath, fieldType)
}

func hasField(sch openapi.Schema, fieldPath, fieldType string) bool {
	return false
}
