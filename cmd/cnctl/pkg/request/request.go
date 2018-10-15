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

package request

import (
	"bytes"
	"log"
	"text/template"

	"github.com/ghodss/yaml"
	"k8s.io/kubectl/cmd/cnctl/pkg/apis/cli/v1alpha1"
	pkgcobra "k8s.io/kubectl/cmd/cnctl/pkg/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
)

func GenerateUnstructured(t, name string, values pkgcobra.Values) *unstructured.Unstructured {
	s := GenerateBytes(t, name, values)

	// Transform the request body into an Unstructure object
	obj := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal(s, &obj.Object); err != nil {
		log.Fatalf("could not create request body", err)
	}

	return obj
}

func GenerateBytes(t, name string, values pkgcobra.Values) []byte {
	// GenerateUnstructured the request body string
	temp, err := template.New(name).Parse(t)
	if err != nil {
		log.Fatalf("could not create request body", err)
	}

	body := &bytes.Buffer{}
	if err := temp.Execute(body, values); err != nil {
		log.Fatalf("could not create request body", err)
	}

	return body.Bytes()
}

// SaveResponseValues parses the items specified by JsonPath from the response object back into the Flags struct
// so that the response is available in the output template
func SaveResponseValues(resp map[string]interface{}, values []v1alpha1.ResponseValue, res *pkgcobra.Values) error {
	if res.Responses.Strings == nil {
		res.Responses.Strings = map[string]*string{}
	}
	for _, v := range values {
		j := jsonpath.New(v.Name)
		buf := &bytes.Buffer{}
		if err := j.Parse(v.JsonPath); err != nil {
			return err
		}
		if err := j.Execute(buf, resp); err != nil {
			return err
		}
		s := buf.String()
		res.Responses.Strings[v.Name] = &s
	}
	return nil
}
