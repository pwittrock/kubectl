// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/kubectl/pkg/framework/resource"
	"k8s.io/kubectl/pkg/framework/resource/flags"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
}

var supportedFields = []fieldDef{
	{"image", []string{"spec", "template", "spec", "containers", "image"}},
	{"cpu-limits", []string{"spec", "template", "spec", "containers", "resources", "limits", "cpu"}},
	{"memory-limits", []string{"spec", "template", "spec", "containers", "resources", "limits", "memory"}},
	{"limits", []string{"spec", "template", "spec", "containers", "resources", "limits"}},
	{"ports", []string{"spec", "template", "spec", "containers", "ports"}},
	{"env", []string{"spec", "template", "spec", "containers", "env"}},
}

type fieldDef struct {
	name string
	path []string
}

func init() {

	p := resource.NewParser()
	resources, e := p.Resources()
	if e != nil {
		panic(e)
	}

	cmds := map[string]*cobra.Command{}

	for _, field := range supportedFields {
		resources = resources.Filter(&FieldFilter{
			resource.EmptyFilter{},
			field.path})

		for k, versions := range resources {
			if _, found := cmds[k]; !found {
				cmds[k] = &cobra.Command{
					Use: k,
				}
			}
			rcmd := cmds[k]

			fcmd := &cobra.Command{
				Use: field.name,
			}
			rcmd.AddCommand(fcmd)

			builder := flags.NewFlagBuilder()
			fn, err := builder.BuildObject(
				fcmd,
				versions[0],
				field.path)
			if err != nil {
				panic(err)
			}

			fcmd.Run = func(cmd *cobra.Command, args []string) {
				value := fn()
				out, err := yaml.Marshal(value)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", out)
			}
		}

	}

	RootCmd.AddCommand(setCmd)
	for _, c := range cmds {
		setCmd.AddCommand(c)
	}
}

type FieldFilter struct {
	resource.EmptyFilter
	path []string
}

func (f *FieldFilter) Resource(r *resource.Resource) bool {
	return r.HasField(f.path)
}
