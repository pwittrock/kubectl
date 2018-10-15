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

package v1alpha1

// ServicCommand defines a command that is dynamically defined as an annotation on a Service
type ServicCommand struct {
	Command Command `json:"command"
`
	// Requests are the requests the command will send to the services.
	// +optional
	Requests []ServiceRequest `json:"requests,omitempty"`

	// Output is a go-template used write the command output.  It may reference values specified as flags using
	// {{index .Flags.Strings "flag-name"}}, {{index .Flags.Ints "flag-name"}}, {{index .Flags.Bools "flag-name"}},
	// {{index .Flags.Floats "flag-name"}}.
	// It may referene values from the request responses using {{index .Responses.Strings "response-value-name"}}.
	//
	// Example:
	// 		deployment.apps/{{index .Responses.Strings "responsename"}} created
	//
	// +optional
	Output string `json:"output,omitempty"`
}

type ServiceOperation string

const (
	HTTP_POST   ServiceOperation = "Post"
	HTTP_PUT                     = "Put"
	HTTP_DELETE                  = "Delete"
	HTTP_GET                     = "Get"
)

type ServiceRequest struct {
	ServiceName      string `json:"serviceName"`
	ServiceNamespace string `json:"serviceNamespace"`

	// Protocol is either "http" or "https"
	Protocol string `json:"protocol"`
	Port     string `json:"port"`

	// Operation is the type of operation to perform for the request.  One of: Create, Update, Delete, Get, Patch
	Operation ServiceOperation `json:"operation"`

	// BodyTemplate is a go-template for the request Body.  It may reference values specified as flags using
	// {{index .Flags.Strings "flag-name"}}, {{index .Flags.Ints "flag-name"}}, {{index .Flags.Bools "flag-name"}},
	// {{index .Flags.Floats "flag-name"}}
	//
	// Example:
	//      apiVersion: apps/v1
	//      kind: Deployment
	//      metadata:
	//        name: {{index .Flags.Strings "name"}}
	//        namespace: {{index .Flags.Strings "namespace"}}
	//        labels:
	//          app: nginx
	//      spec:
	//        replicas: {{index .Flags.Ints "replicas"}}
	//        selector:
	//          matchLabels:
	//            app: {{index .Flags.Strings "name"}}
	//        template:
	//          metadata:
	//            labels:
	//              app: {{index .Flags.Strings "name"}}
	//          spec:
	//            containers:
	//            - name: {{index .Flags.Strings "name"}}
	//              image: {{index .Flags.Strings "image"}}
	//
	// +optional
	BodyTemplate string `json:"bodyTemplate,omitempty"`

	// +optional
	UrlTemplate string `json:"urlTemplate,omitempty"`

	// +optional
	ParamsTemplates map[string]string `json:"paramsTemplates,omitempty"`

	// SaveResponseValues are values read from the response and saved in {{index .Responses.Strings "flag-name"}}.
	// They may be used in the ResourceCommand.Output go-template.
	//
	// Example:
	//		- name: responsename
	//        jsonPath: "{.metadata.name}"
	//
	// +optional
	SaveResponseValues []ResponseValue `json:"saveResponseValues,omitempty"`
}

// ServiceCommandList contains a list of Commands
type ServiceCommandList struct {
	Items []ServicCommand `json:"items"`
}
