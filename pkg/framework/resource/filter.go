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
	"fmt"
	"k8s.io/kubectl/pkg/framework/merge"
	"k8s.io/kubernetes/pkg/kubectl/apply"
	"strings"
)

type Filter interface {
	Resource(*Resource) bool
	SubResource(*SubResource) bool
}

type EmptyFilter struct {
}

func (*EmptyFilter) Resource(*Resource) bool {
	return true
}

func (*EmptyFilter) SubResource(*SubResource) bool {
	return true
}

type SkipSubresourceFilter struct {
	EmptyFilter
}

func (*SkipSubresourceFilter) SubResource(sr *SubResource) bool {
	return !strings.HasSuffix(sr.Resource.Name, "/status")
}

type AndFilter struct {
	Filters []Filter
}

func (a *AndFilter) Resource(r *Resource) bool {
	for _, f := range a.Filters {
		if !f.Resource(r) {
			return false
		}
	}
	return true
}

func (a *AndFilter) SubResource(sr *SubResource) bool {
	for _, f := range a.Filters {
		if !f.SubResource(sr) {
			return false
		}
	}
	return true
}

type OrFilter struct {
	Filters []Filter
}

func (a *OrFilter) Resource(r *Resource) bool {
	for _, f := range a.Filters {
		if f.Resource(r) {
			return true
		}
	}
	return false
}

func (a *OrFilter) SubResource(sr *SubResource) bool {
	for _, f := range a.Filters {
		if f.SubResource(sr) {
			return true
		}
	}
	return false
}

type FieldFilter struct {
	EmptyFilter
	path []string
}

func NewFieldFilter(path []string) *FieldFilter {
	return &FieldFilter{EmptyFilter{}, path}
}

func (f *FieldFilter) Resource(r *Resource) bool {
	return r.HasField(f.path)
}

type PrefixStrategy struct {
	merge.EmptyStrategy
	prefix string
}

func (fs *PrefixStrategy) MergePrimitive(element apply.PrimitiveElement) (apply.Result, error) {
	return apply.Result{MergedResult: fmt.Sprintf("%s%v", fs.prefix, element.GetRemote())}, nil
}
