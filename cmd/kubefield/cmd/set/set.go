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

package set

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	kubefield "k8s.io/kubectl/cmd/kubefield/cmd"
	"k8s.io/kubectl/cmd/kubefield/pkg"
	"k8s.io/kubectl/pkg/framework"
	"k8s.io/kubectl/pkg/framework/resource"
	resourcecmd "k8s.io/kubectl/pkg/framework/resource/cmd"
	"k8s.io/kubectl/pkg/framework/resource/flags"
	"k8s.io/kubernetes/pkg/kubectl/apply/parse"
	"k8s.io/kubernetes/pkg/kubectl/apply/strategy"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
}

var supportedFields = []pkg.FieldDef{
	{"container-labels", []string{"spec", "template", "metadata", "labels"}},
	{"cpu-limits", []string{"spec", "template", "spec", "containers", "resources", "limits", "cpu"}},
	{"cpu-requests", []string{"spec", "template", "spec", "containers", "resources", "requests", "cpu"}},
	{"env", []string{"spec", "template", "spec", "containers", "env"}},
	{"image", []string{"spec", "template", "spec", "containers", "image"}},
	{"labels", []string{"metadata", "labels"}},
	{"memory-limits", []string{"spec", "template", "spec", "containers", "resources", "limits", "memory"}},
	{"memory-requests", []string{"spec", "template", "spec", "containers", "resources", "requests", "memory"}},
	{"ports", []string{"spec", "template", "spec", "containers", "ports"}},
	{"replicas", []string{"spec", "replicas"}},
	{"selector", []string{"spec", "selector", "matchLabels"}},
	{"name", []string{"metadata", "name"}},
}

func init() {
	builder := resourcecmd.NewBuilder()
	for _, field := range supportedFields {
		builder.BuildCmdsForResources(resource.NewFieldFilter(field.Path), &Buildable{field: field})
	}

	kubefield.RootCmd.AddCommand(setCmd)
	for _, c := range builder.Cmds() {
		setCmd.AddCommand(c)
	}
}

type Buildable struct {
	field pkg.FieldDef
}

func (b *Buildable) Build(rs *resource.Resource) *cobra.Command {
	cmd := &cobra.Command{
		Use: b.field.Name,
	}
	builder := flags.NewFlagBuilder()
	fn, err := builder.BuildObject(cmd, rs, b.field.Path)
	if err != nil {
		log.Fatalf("%v", err)
	}
	runnable := &Runnable{
		output: resourcecmd.OutputFn(cmd),
		local:  fn,
	}
	cmd.Run = resourcecmd.RunFn(runnable.Run)
	return cmd
}

type Runnable struct {
	output func(io.Writer, interface{})
	local  func() map[string]interface{}
}

func (b *Runnable) Run(obj map[string]interface{}) {
	local := b.local()
	local["apiVersion"] = obj["apiVersion"]
	local["kind"] = obj["kind"]
	local = roundtrip(local)

	result := merge(local, obj)
	b.output(os.Stdout, result)
}

// roundtrip resolves pointers by roundtripping through serialization
func roundtrip(obj map[string]interface{}) map[string]interface{} {
	tmp, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("%v", err)
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(tmp, &result)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return result
}

func merge(local, obj map[string]interface{}) interface{} {
	factory := framework.Factory()
	elementParse := parse.Factory{factory.GetResources()}
	elem, err := elementParse.CreateElement(local, local, obj)
	if err != nil {
		log.Fatalf("%v", err)
	}
	result, err := elem.Merge(strategy.Create(strategy.Options{}))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return result.MergedResult
}
