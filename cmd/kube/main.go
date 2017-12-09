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

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/cmd/kube/version17"
	"k8s.io/kubectl/cmd/kube/version18"
	kubecurl "k8s.io/kubectl/cmd/kubecurl/cmd"
	kubefield "k8s.io/kubectl/cmd/kubefield/cmd"
	_ "k8s.io/kubectl/cmd/kubefield/cmd/clear"
	_ "k8s.io/kubectl/cmd/kubefield/cmd/get"
	_ "k8s.io/kubectl/cmd/kubefield/cmd/patch"
	_ "k8s.io/kubectl/cmd/kubefield/cmd/prefix"
	_ "k8s.io/kubectl/cmd/kubefield/cmd/set"
	"k8s.io/kubectl/pkg/framework"
	"k8s.io/kubernetes/pkg/kubectl/util/logs"
)

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kube",
	Short: "",
	Long:  ``,
}

func Run() error {
	logs.InitLogs()
	defer logs.FlushLogs()

	factory := framework.Factory()
	version, err := factory.Version()
	if err != nil {
		version = "1.8"
	}

	var cmd *cobra.Command

	switch version {
	case "1.7":
		cmd = version17.Cmd()
	case "1.8":
		cmd = version18.Cmd()
	default:
		panic(fmt.Errorf("Version %s not supported", version))
	}

	cmd.Use = "ctl"
	cmd.Short = "new home for kubectl"
	RootCmd.AddCommand(cmd)

	kubefield.RootCmd.Use = "field"
	RootCmd.AddCommand(kubefield.RootCmd)

	kubecurl.RootCmd.Use = "subresource"
	RootCmd.AddCommand(kubecurl.RootCmd)

	return RootCmd.Execute()
}
