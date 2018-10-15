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
	"fmt"
	"log"
	"os"

	"k8s.io/kubectl/cmd/cnctl/pkg/apis/cli/v1alpha1"
	pkgcobra "k8s.io/kubectl/cmd/cnctl/pkg/cobra"
	"k8s.io/kubectl/cmd/cnctl/pkg/output"
	"k8s.io/kubectl/cmd/cnctl/pkg/request"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ParseCommand parses the dynamic command into a cobra command
func ParseCommand(config *rest.Config,
	cmd *v1alpha1.ServicCommand, k *kubernetes.Clientset) *cobra.Command {
	cbra, f := pkgcobra.ParseCommand(&cmd.Command)
	values := pkgcobra.Values{Flags: f}

	cbra.Run = func(c *cobra.Command, args []string) {
		for _, req := range cmd.Requests {
			body := request.GenerateBytes(req.BodyTemplate, cmd.Command.Use+"-service-request", values)
			url := request.GenerateBytes(req.UrlTemplate, cmd.Command.Use+"-service-request", values)
			params := map[string]string{}
			for k, value := range req.ParamsTemplates {
				params[k] = string(request.GenerateBytes(value, cmd.Command.Use+"param-"+k, values))
			}

			// Make the request
			resp, err := doRequest(config, body, url, params, req, k)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			// Save the response values
			if err := request.SaveResponseValues(resp, req.SaveResponseValues, &values); err != nil {
				log.Fatalf("could not parse service response", err)
			}
		}

		output.Write(cmd.Output, cmd.Command.Use+"-service-response", values, os.Stdout)
	}
	return cbra
}

// doRequest makes a request to the apiserver with the object (request body), operation (request http operation),
// group,version,resource (request url).
func doRequest(
	config *rest.Config,
	body []byte,
	url []byte,
	params map[string]string,
	serviceRequest v1alpha1.ServiceRequest,
	k *kubernetes.Clientset) (map[string]interface{}, error) {

	name := fmt.Sprintf("%s:%s:%s/proxy/",
		serviceRequest.Protocol,
		serviceRequest.ServiceName,
		serviceRequest.Port)

	var request *rest.Request
	var err error

	switch serviceRequest.Operation {
	case v1alpha1.HTTP_POST:
		request = k.RESTClient().Post()
	case v1alpha1.HTTP_PUT:
		request = k.RESTClient().Put()
	case v1alpha1.HTTP_DELETE:
		request = k.RESTClient().Delete()
	case v1alpha1.HTTP_GET:
		request = k.RESTClient().Get()
	}

	r := request.
		Prefix("api", "v1").
		Namespace(serviceRequest.ServiceNamespace).
		Resource("services").
		Suffix(name).
		Body(body)

	// TODO: Add RequestURI

	// Add Params
	for k, v := range params {
		r.Param(k, v)
	}

	result := r.Do()
	if result.Error() != nil {
		fmt.Printf("err: %v\n", result.Error())
	}
	u := &unstructured.Unstructured{}
	result.Into(u)

	return u.Object, err
}
