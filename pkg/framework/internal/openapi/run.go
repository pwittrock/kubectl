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

package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (builder *cmdBuilderImpl) BuildRun(cmd *cobra.Command, resource v1.APIResource, request map[string]interface{}) {
	cmd.Run = func(cmd *cobra.Command, args []string) {
		out, _ := json.Marshal(request)
		fmt.Printf("Request:\n%s\n", out)

		meta := request["metadata"].(map[string]interface{})
		name := meta["name"].(*string)
		namespace := meta["namespace"].(*string)
		fmt.Printf("P%s %s\n", *name, *namespace)

		result := builder.rest.Put().
			Prefix("apis", resource.Group, resource.Version).
			Namespace(*namespace).
			Resource(builder.resource(resource)).
			SubResource(builder.operation(resource)).
			Name(*name).
			Body(out)

		fmt.Printf("URL: %v\n", result.URL().Path)
		resp, err := result.DoRaw()
		fmt.Printf("Response:\n%s\nError: %v\n", resp, err)
	}
}
