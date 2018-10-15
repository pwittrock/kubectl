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

package service

import (
	"encoding/json"
	"log"
	"reflect"

	"k8s.io/kubectl/cmd/cnctl/pkg/apis/cli/v1alpha1"
	v1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

const annotation = "cli.sigs.k8s.io/cli.v1alpha1.ServiceCommandList"

// ListCommands fetches the list of dynamic commands published as Annotations on CRDs
func ListCommands(client v1beta1client.CustomResourceDefinitionInterface) (v1alpha1.ServiceCommandList, error) {
	cmds := v1alpha1.ServiceCommandList{}

	// List all CRDs with Commands
	crds, err := client.List(v1.ListOptions{LabelSelector: annotation})
	if err != nil {
		return cmds, err
	}

	for _, crd := range crds.Items {
		// Get the ServiceCommand json
		s := crd.Annotations[annotation]
		if len(s) == 0 {
			continue
		}

		// Unmarshall the annotation value into a ServiceCommandList
		l := v1alpha1.ServiceCommandList{}
		err := json.Unmarshal([]byte(s), &l)
		if err != nil {
			log.Printf("failed to parse commands for CRD %s: %v\n", crd.Name, err)
			continue
		}

		// Verify we parsed something
		if reflect.DeepEqual(l, v1alpha1.ServiceCommandList{}) {
			log.Printf("no commands for CRD %s: %s\n", crd.Name, s)
			continue
		}

		// Add the commands to the list
		for _, i := range l.Items {
			if len(i.Requests) > 0 {
				cmds.Items = append(cmds.Items, i)
			}
		}
	}
	return cmds, err
}
