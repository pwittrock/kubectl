/*
Copyright 2018 The Kubernetes Authors.

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

package cobra

import (
	"k8s.io/kubectl/cmd/cnctl/pkg/apis/cli/v1alpha1"
	"github.com/spf13/cobra"
)

// ParseCommand parses the dynamic ResourceCommand into a cobra ResourceCommand
func ParseCommand(cmd *v1alpha1.Command) (*cobra.Command, Flags) {
	values := Flags{}

	// Parse the values
	cbra := &cobra.Command{
		Use:        cmd.Use,
		Short:      cmd.Short,
		Long:       cmd.Long,
		Example:    cmd.Example,
		Version:    cmd.Version,
		Deprecated: cmd.Deprecated,
		Aliases:    cmd.Aliases,
		SuggestFor: cmd.SuggestFor,
	}

	// Setup the flags
	for _, cmdFlag := range cmd.Flags {
		switch cmdFlag.Type {
		case v1alpha1.STRING:
			if values.Strings == nil {
				values.Strings = map[string]*string{}
			}
			values.Strings[cmdFlag.Name] = cbra.Flags().String(cmdFlag.Name, cmdFlag.StringValue, cmdFlag.Description)
		case v1alpha1.STRING_SLICE:
			if values.StringSlices == nil {
				values.StringSlices = map[string]*[]string{}
			}
			values.StringSlices[cmdFlag.Name] = cbra.Flags().StringSlice(
				cmdFlag.Name, cmdFlag.StringSliceValue, cmdFlag.Description)
		case v1alpha1.INT:
			if values.Ints == nil {
				values.Ints = map[string]*int32{}
			}
			values.Ints[cmdFlag.Name] = cbra.Flags().Int32(cmdFlag.Name, cmdFlag.IntValue, cmdFlag.Description)
		case v1alpha1.FLOAT:
			if values.Floats == nil {
				values.Floats = map[string]*float64{}
			}
			values.Floats[cmdFlag.Name] = cbra.Flags().Float64(cmdFlag.Name, cmdFlag.FloatValue, cmdFlag.Description)
		case v1alpha1.BOOL:
			if values.Bools == nil {
				values.Bools = map[string]*bool{}
			}
			values.Bools[cmdFlag.Name] = cbra.Flags().Bool(cmdFlag.Name, cmdFlag.BoolValue, cmdFlag.Description)
		}
	}

	return cbra, values
}

type Values struct {
	Flags     Flags
	Responses Flags
}

// Flags contains flag values setup for the cobra command
type Flags struct {
	Strings      map[string]*string
	Ints         map[string]*int32
	Bools        map[string]*bool
	Floats       map[string]*float64
	StringSlices map[string]*[]string
}

// AddTo adds to as a sub-command of cmd.  If a Path is specified in command, AddTo will
// ensure that path exists by creating intermediary sub-commands.
func AddTo(to, cmd *cobra.Command, command v1alpha1.Command) {
	// For each element on the Path
	for _, p := range command.Path {
		// Make sure the subcommand exists
		for _, c := range to.Commands() {
			if c.Use == p {
				// Found, continue on to next part of the Path
				to = c
				continue
			}
		}
		// Missing, create the sub-command
		cbra := &cobra.Command{Use: p}
		to.AddCommand(cbra)
		to = cbra
	}

	to.AddCommand(cmd)
}
