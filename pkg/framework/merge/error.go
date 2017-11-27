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

import (
	"fmt"
	"k8s.io/kubernetes/pkg/kubectl/apply"
)

type ErrorStrategy struct {
}

func (fs *ErrorStrategy) MergeList(element apply.ListElement) (apply.Result, error) {
	return apply.Result{}, fmt.Errorf("MergeList not implemented")
}

func (fs *ErrorStrategy) MergeMap(element apply.MapElement) (apply.Result, error) {
	return apply.Result{}, fmt.Errorf("MergeMap not implemented")
}

func (fs *ErrorStrategy) MergeType(element apply.TypeElement) (apply.Result, error) {
	return apply.Result{}, fmt.Errorf("MergeType not implemented")
}

func (fs *ErrorStrategy) MergePrimitive(element apply.PrimitiveElement) (apply.Result, error) {
	return apply.Result{}, fmt.Errorf("MergePrimitive not implemented")
}

func (fs *ErrorStrategy) MergeEmpty(element apply.EmptyElement) (apply.Result, error) {
	return apply.Result{}, fmt.Errorf("MergeEmpty not implemented")
}
