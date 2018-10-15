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

package main

import (
	"log"

	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"

	"k8s.io/kubectl/cmd/cnctl/pkg/client"
	pkgcobra "k8s.io/kubectl/cmd/cnctl/pkg/cobra"
	"k8s.io/kubectl/cmd/cnctl/pkg/resource"
	"k8s.io/kubectl/cmd/cnctl/pkg/service"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var rootCmd = &cobra.Command{
	Use:   "cnctl",
	Short: "cnctl is cli for working with Cloud Native Resources",
}

func main() {
	// Create clients
	config := client.GetConfigOrDie()
	ccs := clientset.NewForConfigOrDie(config).ApiextensionsV1beta1().CustomResourceDefinitions()
	i := dynamic.NewForConfigOrDie(config)
	k := kubernetes.NewForConfigOrDie(config)

	// ResourceCommands
	rcmds, err := resource.ListCommands(ccs)
	if err != nil {
		log.Fatalf("could not list resource commands", err)
	}
	for _, cmd := range rcmds.Items {
		pkgcobra.AddTo(rootCmd, resource.ParseCommand(&cmd, i), cmd.Command)
	}

	// ServiceCommands
	scmds, err := service.ListCommands(ccs)
	if err != nil {
		log.Fatalf("could not list service commands", err)
	}
	for _, cmd := range scmds.Items {
		pkgcobra.AddTo(rootCmd, service.ParseCommand(config, &cmd, k), cmd.Command)
	}

	// Execute Cobra
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
