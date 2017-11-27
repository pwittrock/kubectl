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

package merge

import "k8s.io/kubernetes/pkg/kubectl/apply"

type EmptyStrategy struct {
}

func (fs *EmptyStrategy) MergeList(element apply.ListElement) (apply.Result, error) {
	return apply.Result{}, nil
}

func (fs *EmptyStrategy) MergeMap(element apply.MapElement) (apply.Result, error) {
	return apply.Result{}, nil
}

func (fs *EmptyStrategy) MergeType(element apply.TypeElement) (apply.Result, error) {
	return apply.Result{}, nil
}

func (fs *EmptyStrategy) MergePrimitive(element apply.PrimitiveElement) (apply.Result, error) {
	return apply.Result{}, nil
}

func (fs *EmptyStrategy) MergeEmpty(element apply.EmptyElement) (apply.Result, error) {
	return apply.Result{}, nil
}
