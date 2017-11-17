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
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubectl/pkg/framework/openapi"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "do",
	Short: "do performs write operations against Kubernetes APIs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.PersistentFlags().String("api-group", "", "")
	setCmd.Flag("api-group").Hidden = true
	setCmd.PersistentFlags().String("api-version", "", "")
	setCmd.Flag("api-version").Hidden = true
	builder := openapi.NewCmdBuilder()
	cmds, _ := builder.BuildCommands("PUT", sets.NewString("update", "create"))
	for _, cmd := range cmds {
		setCmd.AddCommand(cmd)
	}
}
