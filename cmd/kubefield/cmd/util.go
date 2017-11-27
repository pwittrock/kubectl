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
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
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

func groupVersion(groupVersion string) (string, string) {
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
