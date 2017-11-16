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
)

func (builder *cmdBuilderImpl) BuildCommands() ([]*cobra.Command, error) {
	list, err := builder.listResources()
	if err != nil {
		panic(err)
	}

	// Setup a sub command for each operation
	parentCmds := map[string]*cobra.Command{}

	for _, resource := range list {
		if builder.isCmd(resource) {
			// Don't expose multiple versions of the same resource
			if builder.done(resource) {
				continue
			}

			// Setup the command
			cmd, err := builder.buildCmd(resource)
			if err != nil {
				panic(err)
			}

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
			builder.buildRun(cmd, resource, request)
		}
	}

	cmds := []*cobra.Command{}
	for _, cmd := range parentCmds {
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}
