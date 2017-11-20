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

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

type fieldDefinition struct {
	name  string
	paths []string
	value fieldValue
}

type fieldValue struct {
	valueSources []fieldSource
	valueType    string
	value        interface{}
}

type fieldSource int

const (
	flag fieldSource = iota
)

var fields = []fieldDefinition{
	{
		"image",
		[]string{"spec.template.spec.containers.image"},
		fieldValue{valueType: "string", valueSources: []fieldSource{flag}},
	},
}

type Resource struct {
	resource         v1.APIResource
	groupVersionKind schema.GroupVersionKind
	openapiSchema    openapi.Schema
}

func getResourcesWithField(def fieldDefinition) []*Resource {
	result := []*Resource{}

	return result
}

func SetCmd(def fieldDefinition) []cobra.Command {
	result := []cobra.Command{}

	// Use openapi to find all models with the given path

	// Register each

	return result
}

func Set(def fieldDefinition) {
	// Validate def.value against def.valueType using openapi

	// Parse flag value into json

	// Merge json into destination
}
