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

package output

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"text/template"

	pkgcobra "k8s.io/kubectl/cmd/cnctl/pkg/cobra"
)

// Write parses the outputTemplate and executes it with values, writing the output to writer.
func Write(outputTemplate, name string, values pkgcobra.Values, writer io.Writer) {
	temp, err := template.New(name).Parse(outputTemplate)
	if err != nil {
		log.Fatalf("could not parse output", err)
	}
	buff := &bytes.Buffer{}
	if err := temp.Execute(buff, values); err != nil {
		log.Fatalf("could not parse output", err)
	}

	// Print the output
	fmt.Fprintf(writer, buff.String())
}
