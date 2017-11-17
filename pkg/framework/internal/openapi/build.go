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
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (builder *cmdBuilderImpl) BuildCommands(requestType string) ([]*cobra.Command, error) {
	list, err := builder.listResources()
	if err != nil {
		panic(err)
	}

	// Setup a sub command for each operation
	parentCmds := map[string]*cobra.Command{}
	parentResources := map[string]*v1.APIResource{}

	for _, resource := range list {
		if builder.isResource(resource) {
			parentResources[resource.Name] = resource
		}
	}

	//fmt.Printf("Parents %+v\n", parentResources)
	for _, resource := range list {
		// Only operate on subresources
		if !builder.isSubResource(resource) {
			continue
		}

		// Set the gvk from the parent if it is missing and the parent exists
		if parent, found := parentResources[builder.resource(resource)]; found {
			builder.setGroupVersionFromParentIfMissing(resource, parent)
		}

		// If this subresource cannot be used as a cmd, continue
		if !builder.isCmd(resource) {
			continue
		}

		// Don't expose multiple versions of the same resource
		if builder.done(resource) {
			continue
		}

		// Setup the command
		cmd, err := builder.buildCmd(resource)
		if err != nil {
			panic(err)
		}

		// Mark this resource as done for this operation
		builder.add(resource)

		operation := builder.operation(resource)
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
