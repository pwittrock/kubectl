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
	p := resource.NewParser()
	resources, e := p.Resources()
	resources = resources.Filter(&resource.SkipSubresourceFilter{})
	if e != nil {
		panic(e)
	}

	builder := kubefield.NewCmdBuilder()
	fn, err := builder.BuildObject(
		setCmd,
		resources["deployments"][0],
		[]string{"spec", "template", "spec", "containers", "image"})

	if err != nil {
		panic(err)
	}
	setCmd.Run = func(cmd *cobra.Command, args []string) {
		value := fn()
		out, err := yaml.Marshal(value)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", out)
	}

	RootCmd.AddCommand(setCmd)
}
