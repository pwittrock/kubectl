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

package resource

import (
	"fmt"
	"log"
	"os"

	"k8s.io/kubectl/cmd/cnctl/pkg/apis/cli/v1alpha1"
	pkgcobra "k8s.io/kubectl/cmd/cnctl/pkg/cobra"
	"k8s.io/kubectl/cmd/cnctl/pkg/output"
	"k8s.io/kubectl/cmd/cnctl/pkg/request"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// ParseCommand parses the dynamic command into a cobra command
func ParseCommand(cmd *v1alpha1.ResourceCommand, i dynamic.Interface) *cobra.Command {
	cbra, f := pkgcobra.ParseCommand(&cmd.Command)
	values := pkgcobra.Values{Flags: f}

	cbra.Run = func(c *cobra.Command, args []string) {
		for _, req := range cmd.Requests {
			obj := request.GenerateUnstructured(req.BodyTemplate, cmd.Command.Use+"-resource-request", values)

			// Make the request
			gvr := schema.GroupVersionResource{Resource: req.Resource, Version: req.Version, Group: req.Group}
			resp, err := doRequest(obj, gvr, req.Operation, i)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			// Save the response values
			if err := request.SaveResponseValues(resp.Object, req.SaveResponseValues, &values); err != nil {
				log.Fatalf("could not parse resource response", err)
			}
		}

		output.Write(cmd.Output, cmd.Command.Use+"-resource-response", values, os.Stdout)
	}
	return cbra
}

// doRequest makes a request to the apiserver with the object (request body), operation (request http operation),
// group,version,resource (request url).
func doRequest(
	obj *unstructured.Unstructured,
	gvr schema.GroupVersionResource,
	op v1alpha1.ResourceOperation,
	i dynamic.Interface) (*unstructured.Unstructured, error) {

	req := i.Resource(gvr).Namespace(obj.GetNamespace())
	var resp = &unstructured.Unstructured{}
	var err error

	// TODO: Add support for specifying options
	switch op {
	case v1alpha1.CREATE_RESOURCE:
		resp, err = req.Create(obj, metav1.CreateOptions{})
	case v1alpha1.DELETE_RESOURCE:
		err = req.Delete(obj.GetName(), &metav1.DeleteOptions{})
	case v1alpha1.UPDATE_RESOURCE:
		resp, err = req.Update(obj, metav1.UpdateOptions{})
	case v1alpha1.GET_RESOURCE:
		resp, err = req.Get(obj.GetName(), metav1.GetOptions{})
	case v1alpha1.PATCH_RESOURCE:
	}
	return resp, err
}
