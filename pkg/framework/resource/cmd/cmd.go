/*
Copyright 2017 The Kubernetes Authors.

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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/kubectl/pkg/framework/resource"
	"log"
	"sync"
)

type Runnable interface {
	Run(map[string]interface{})
}

type Buildable interface {
	Build(*resource.Resource) *cobra.Command
}

type CmdBuilder struct {
	cmds map[string]*cobra.Command
	r    resource.Resources
	err  error
	init sync.Once
}

func NewBuilder() *CmdBuilder {
	return &CmdBuilder{
		map[string]*cobra.Command{},
		nil,
		nil,
		sync.Once{},
	}
}

type ResourceCmd struct {
	Cmd      *cobra.Command
	Resource *resource.Resource
}

func (b *CmdBuilder) cmd(resource string) *cobra.Command {
	if b.cmds[resource] == nil {
		b.cmds[resource] = &cobra.Command{Use: resource}
	}
	return b.cmds[resource]
}

func (b *CmdBuilder) version(versions []*resource.Resource) *resource.Resource {
	return versions[0]
}

func (b *CmdBuilder) resources() (resource.Resources, error) {
	b.init.Do(func() {
		p := resource.NewParser()
		resources, e := p.Resources()
		b.err = e
		b.r = resources
	})
	return b.r, b.err
}

// CmdForResources creates a new command and attaches subcommands for each resource matching the filter
func (b *CmdBuilder) BuildCmdsForResources(filter resource.Filter, build Buildable) error {
	resources, err := b.resources()
	if err != nil {
		return err
	}
	resources = resources.Filter(filter)

	for k, versions := range resources {
		// Lookup the command for this resource and add the child to it
		parent := b.cmd(k)
		parent.AddCommand(build.Build(b.version(versions)))
	}
	return nil
}

func (b *CmdBuilder) Cmds() []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, c := range b.cmds {
		cmds = append(cmds, c)
	}
	return cmds
}

func RunFn(fn func(map[string]interface{})) func(cmd *cobra.Command, args []string) {

	// Read objects from stdin
	if stat, err := os.Stdin.Stat(); err != nil && stat.Size() > 0 {
		return func(cmd *cobra.Command, args []string) {
			in, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("%v", err)
			}
			remote := map[string]interface{}{}
			yaml.Unmarshal(in, &remote)
			fn(remote)
		}
	}

	// Read objects from file specified as args
	return func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			in, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatalf("%v", err)
			}
			remote := map[string]interface{}{}
			yaml.Unmarshal(in, &remote)
			fn(remote)
		}
	}
}

func OutputFn(cmd *cobra.Command) func(io.Writer, interface{}) {
	output := cmd.PersistentFlags().String("output-format", "yaml", "maybe [yaml, json, patch]")
	dest := cmd.PersistentFlags().String("output-destination", "", "destination to write output to")

	var out []byte
	var err error
	return func(writer io.Writer, value interface{}) {
		if *output == "yaml" {
			out, err = yaml.Marshal(value)
			if err != nil {
				log.Fatalf("%v", err)
			}
		} else if *output == "json" {
			out, err = json.MarshalIndent(value, "", "    ")
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		} else if *output == "patch" {
			out, err = json.Marshal(value)
			if err != nil {
				log.Fatalf("%v", err)
			}
		} else {
			log.Fatalf("%v", err)
		}

		if len(*dest) > 0 {
			err := ioutil.WriteFile(*dest, out, 0600)
			if err != nil {
				log.Fatalf("%v", err)
			}
		} else {
			fmt.Fprintf(writer, "%s", out)
		}
	}
}
