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
	"k8s.io/kubectl/pkg/framework/resource"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
	Run:   Do,
}

func init() {
	RootCmd.AddCommand(setCmd)
}

func Do(cmd *cobra.Command, args []string) {
	p := resource.NewParser()
	resources, e := p.Resources()
	resources = resources.Filter(&resource.SkipSubresourceFilter{})
	if e != nil {
		panic(e)
	}

	for _, resource := range resources.SortKeys() {
		versions := resources[resource]
		version := versions[0]
		fmt.Printf("%s/%s/%s %v\n", version.ApiGroupVersion.Group, version.ApiGroupVersion.Version, resource, len(versions))
		for _, subresource := range version.SubResources {
			fmt.Printf("\t%s\n", subresource.Resource.Name)
		}
	}
}
