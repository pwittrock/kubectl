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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubectl/pkg/framework/openapi"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "read",
	Short: "read performs read operations against Kubernetes APIs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	builder := openapi.NewCmdBuilder()
	getCmd.PersistentFlags().String("api-group", "", "")
	getCmd.PersistentFlags().String("api-version", "", "")

	cmds, _ := builder.BuildCommands("read", "GET", sets.NewString("get"))
	for _, cmd := range cmds {
		getCmd.AddCommand(cmd)
	}
}
