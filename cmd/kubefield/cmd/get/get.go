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

package prefix

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
	kubefield "k8s.io/kubectl/cmd/kubefield/cmd"
	"k8s.io/kubectl/cmd/kubefield/pkg"
	"k8s.io/kubectl/pkg/framework/resource"
	resourcecmd "k8s.io/kubectl/pkg/framework/resource/cmd"
)

// getCmd represents the set command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  ``,
}

var supportedFields = []pkg.FieldDef{
	//{"container-labels", []string{"spec", "template", "metadata", "labels"}},
	//{"env", []string{"spec", "template", "spec", "containers", "env"}},
	//{"image", []string{"spec", "template", "spec", "containers", "image"}},
	//{"labels", []string{"metadata", "labels"}},
	{"name", []string{"metadata", "name"}},
	//{"selector", []string{"spec", "selector", "matchLabels"}},
}

type Buildable struct {
	field pkg.FieldDef
}

func (b *Buildable) Build(r *resource.Resource) *cobra.Command {
	cmd := &cobra.Command{
		Use: b.field.Name,
	}
	run := &Runnable{
		field:    b.field,
		output:   resourcecmd.OutputFn(cmd),
		resource: r,
	}
	cmd.Run = resourcecmd.RunFn(run.Run)
	return cmd
}

type Runnable struct {
	field    pkg.FieldDef
	output   func(io.Writer, interface{})
	resource *resource.Resource
}

func (b *Runnable) Run(obj map[string]interface{}) {
	var val interface{}
	_, err := b.resource.Field(b.field.Path, obj,
		func(i interface{}, _ openapi.BaseSchema, _ openapi.Schema) interface{} {
			val = i
			return i
		})
	if err != nil {
		log.Fatalf("%v", err)
	}
	b.output(os.Stdout, val)
}

func init() {
	builder := resourcecmd.NewBuilder()
	for _, field := range supportedFields {
		builder.BuildCmdsForResources(resource.NewFieldFilter(field.Path), &Buildable{field: field})
	}

	kubefield.RootCmd.AddCommand(getCmd)
	for _, c := range builder.Cmds() {
		getCmd.AddCommand(c)
	}
}
