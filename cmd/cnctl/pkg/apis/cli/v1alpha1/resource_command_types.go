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

// ResourceCommand defines a command that is dynamically defined as an annotation on a CRD
type ResourceCommand struct {
	// Command is the cli Command
	Command Command `json:"command"`

	// Requests are the requests the command will send to the apiserver.
	// +optional
	Requests []ResourceRequest `json:"requests,omitempty"`

	// Output is a go-template used write the command output.  It may reference values specified as flags using
	// {{index .Flags.Strings "flag-name"}}, {{index .Flags.Ints "flag-name"}}, {{index .Flags.Bools "flag-name"}},
	// {{index .Flags.Floats "flag-name"}}.
	//
	// It may also reference values from the responses that were saved using saveResponseValues
	// - {{index .Responses.Strings "response-value-name"}}.
	//
	// Example:
	// 		deployment.apps/{{index .Responses.Strings "responsename"}} created
	//
	// +optional
	Output string `json:"output,omitempty"`
}

type ResourceOperation string

const (
	CREATE_RESOURCE ResourceOperation = "Create"
	UPDATE_RESOURCE                   = "Update"
	DELETE_RESOURCE                   = "Delete"
	GET_RESOURCE                      = "Get"
	PATCH_RESOURCE                    = "Patch"
)

type ResourceRequest struct {
	// Group is the API group of the request endpoint
	//
	// Example: apps
	Group string `json:"group"`

	// Version is the API version of the request endpoint
	//
	// Example: v1
	Version string `json:"version"`

	// Resource is the API resource of the request endpoint
	//
	// Example: deployments
	Resource string `json:"resource"`

	// Operation is the type of operation to perform for the request.  One of: Create, Update, Delete, Get, Patch
	Operation ResourceOperation `json:"operation"`

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

// ResourceCommandList contains a list of Commands
type ResourceCommandList struct {
	Items []ResourceCommand `json:"items"`
}
