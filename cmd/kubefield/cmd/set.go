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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	//"gopkg.in/yaml.v2"
	"github.com/ghodss/yaml"
	"k8s.io/kubectl/pkg/framework"
	"k8s.io/kubectl/pkg/framework/resource"
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

var supportedFields = []fieldDef{
	{"image", []string{"spec", "template", "spec", "containers", "image"}},
	{"cpu-limits", []string{"spec", "template", "spec", "containers", "resources", "limits", "cpu"}},
	{"memory-limits", []string{"spec", "template", "spec", "containers", "resources", "limits", "memory"}},
	{"limits", []string{"spec", "template", "spec", "containers", "resources", "limits"}},
	{"ports", []string{"spec", "template", "spec", "containers", "ports"}},
	{"env", []string{"spec", "template", "spec", "containers", "env"}},
	{"replicas", []string{"spec", "replicas"}},
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

	output := setCmd.PersistentFlags().String("output-format", "yaml", "")
	dest := setCmd.PersistentFlags().String("output-destination", "", "")

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
				} else {
					out, err = json.MarshalIndent(value, "", "    ")
					if err != nil {
						panic(err)
					}
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
