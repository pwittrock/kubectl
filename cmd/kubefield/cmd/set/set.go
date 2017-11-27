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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	kubefield "k8s.io/kubectl/cmd/kubefield/cmd"
	"k8s.io/kubectl/pkg/framework"
	"k8s.io/kubectl/pkg/framework/merge"
	"k8s.io/kubectl/pkg/framework/resource"
	"k8s.io/kubectl/pkg/framework/resource/flags"
	"k8s.io/kubernetes/pkg/kubectl/apply"
	"k8s.io/kubernetes/pkg/kubectl/apply/parse"
	"k8s.io/kubernetes/pkg/kubectl/apply/strategy"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
}

var supportedFields = []fieldDef{
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

	output := setCmd.PersistentFlags().String("output-format", "yaml", "maybe [yaml, json, patch]")
	dest := setCmd.PersistentFlags().String("output-destination", "", "destination to write output to")

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
				var value interface{}
				var out []byte

				file := os.Stdin
				fi, err := file.Stat()
				if err != nil {
					fmt.Println("file.Stat()", err)
				}
				size := fi.Size()

				if size > 0 {
					// Read the file to be patched from the stdin
					in, err := ioutil.ReadAll(os.Stdin)
					if err != nil {
						panic(err)
					}
					remote := map[string]interface{}{}
					yaml.Unmarshal(in, &remote)

					// Copy the apiVersion and kind since merge expects them
					local := fn()
					local["apiVersion"] = remote["apiVersion"]
					local["kind"] = remote["kind"]

					// Roundtrip to get rid of points, which merge doesn't handle
					tmp, err := json.Marshal(local)
					if err != nil {
						panic(err)
					}
					local = map[string]interface{}{}
					err = json.Unmarshal(tmp, &local)
					if err != nil {
						panic(err)
					}

					// Merge the patch into the file read from stdin
					factory := framework.Factory()
					elementParse := parse.Factory{factory.GetResources()}
					elem, err := elementParse.CreateElement(local, local, remote)
					if err != nil {
						panic(err)
					}
					result, err := elem.Merge(strategy.Create(strategy.Options{}))
					if err != nil {
						panic(err)
					}
					value = result.MergedResult
				} else {
					value = fn()
				}

				if *output == "yaml" {
					out, err = yaml.Marshal(value)
					if err != nil {
						panic(err)
					}
				} else if *output == "json" {
					out, err = json.MarshalIndent(value, "", "    ")
					if err != nil {
						panic(err)
					}
				} else if *output == "patch" {
					out, err = json.Marshal(value)
					if err != nil {
						panic(err)
					}
				} else {
					panic(fmt.Errorf("Unknown format %s", *output))
				}

				if len(*dest) > 0 {
					err := ioutil.WriteFile(*dest, out, 0600)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Printf("%s", out)
				}
			}
		}

	}

	kubefield.RootCmd.AddCommand(setCmd)
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

type PrefixStrategy struct {
	merge.EmptyStrategy
	prefix string
}

func (fs *PrefixStrategy) MergePrimitive(element apply.PrimitiveElement) (apply.Result, error) {
	return apply.Result{MergedResult: fmt.Sprintf("%s%v", fs.prefix, element.GetRemote())}, nil
}
