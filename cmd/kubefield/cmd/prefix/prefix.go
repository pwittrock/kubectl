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
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
	kubefield "k8s.io/kubectl/cmd/kubefield/cmd"
	"k8s.io/kubectl/pkg/framework/resource"
	resourcecmd "k8s.io/kubectl/pkg/framework/resource/cmd"
)

// prefixCmd represents the set command
var prefixCmd = &cobra.Command{
	Use:   "prefix",
	Short: "",
	Long:  ``,
}

var supportedFields = []fieldDef{
	{"container-labels", []string{"spec", "template", "metadata", "labels"}},
	{"env", []string{"spec", "template", "spec", "containers", "env"}},
	{"image", []string{"spec", "template", "spec", "containers", "image"}},
	{"labels", []string{"metadata", "labels"}},
	{"name", []string{"metadata", "name"}},
	{"selector", []string{"spec", "selector", "matchLabels"}},
}

type fieldDef struct {
	name string
	path []string
}

type Cmd struct {
	resource *resource.Resource
	field    fieldDef
	prefix   *string
	output   func(io.Writer, interface{})
}

func (b *Cmd) Build(r *resource.Resource) *cobra.Command {
	cmd := &cobra.Command{
		Use: b.field.name,
		Run: resourcecmd.RunFn(b.Run),
	}
	b.prefix = cmd.Flags().String("value", "", "")
	b.resource = r
	b.output = resourcecmd.OutputFn(cmd)
	return cmd
}

func (b *Cmd) Run(obj map[string]interface{}) {
	value, err := b.resource.Field(b.field.path, obj,
		func(i interface{}, _ openapi.BaseSchema, _ openapi.Schema) interface{} {
			if i == nil {
				return fmt.Sprintf("%v", *b.prefix)
			}
			return fmt.Sprintf("%v%v", *b.prefix, i)
		})
	if err != nil {
		log.Fatalf("%v", err)
	}
	b.output(os.Stdout, value)
}

func init() {
	builder := resourcecmd.NewBuilder()
	for _, field := range supportedFields {
		builder.BuildCmdsForResources(resource.NewFieldFilter(field.path), &Cmd{field: field})
	}

	kubefield.RootCmd.AddCommand(prefixCmd)
	for _, c := range builder.Cmds() {
		prefixCmd.AddCommand(c)
	}
}
