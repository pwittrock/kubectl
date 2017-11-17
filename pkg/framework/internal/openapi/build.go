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

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

func (builder *cmdBuilderImpl) BuildCommands(
	name string,
	requestType string,
	verbs sets.String) ([]*cobra.Command, error) {
	list, err := builder.getSubResources()
	if err != nil {
		panic(err)
	}

	// Setup a sub command for each operation
	parentCmds := map[string]*cobra.Command{}

	//fmt.Printf("Parents %+v\n", parentResources)
	for _, subResourceList := range list {
		resource := subResourceList[0]

		// If this subresource cannot be used as a cmd, continue
		if !builder.isCmd(&resource.resource) {
			continue
		}

		// Make sure it supports the verbs required for this command
		actualVerbs := sets.NewString(resource.resource.Verbs...)
		if len(actualVerbs.Intersection(verbs).List()) == 0 {
			continue
		}

		// Setup the command
		versions := []schema.GroupVersion{}
		for _, v := range subResourceList {
			versions = append(versions, v.apiGroupVersion)
		}
		cmd, err := builder.buildCmd(name, &resource.resource, versions)
		if err != nil {
			panic(err)
		}

		operation := builder.operation(&resource.resource)
		if _, found := parentCmds[operation]; !found {
			parentCmds[operation] = &cobra.Command{
				Use: fmt.Sprintf("%v", operation),
			}
		}
		parent := parentCmds[operation]
		parent.AddCommand(cmd)

		// Build the flags
		request, err := builder.buildFlags(cmd, resource)
		if err != nil {
			panic(err)
		}

		// Build the run function
		builder.buildRun(cmd, resource, request, requestType)
	}

	cmds := []*cobra.Command{}
	for _, cmd := range parentCmds {
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}
