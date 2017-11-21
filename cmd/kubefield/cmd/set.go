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
	"k8s.io/kubectl/pkg/cmd/kubefield"
	"k8s.io/kubectl/pkg/framework/resource"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
}

func init() {
	imagePath := []string{"spec", "template", "spec", "containers", "image"}

	p := resource.NewParser()
	resources, e := p.Resources()
	resources = resources.Filter(&FieldFilter{
		resource.EmptyFilter{},
		imagePath})
	if e != nil {
		panic(e)
	}

	cmds := []*cobra.Command{}
	for k, versions := range resources {
		c := &cobra.Command{
			Use: k,
		}
		cmds = append(cmds, c)

		image := &cobra.Command{
			Use: "image",
		}
		c.AddCommand(image)

		builder := kubefield.NewCmdBuilder()
		fn, err := builder.BuildObject(
			image,
			versions[0],
			imagePath)
		if err != nil {
			panic(err)
		}

		image.Run = func(cmd *cobra.Command, args []string) {
			value := fn()
			out, err := yaml.Marshal(value)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", out)
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
