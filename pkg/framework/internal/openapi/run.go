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

func (builder *cmdBuilderImpl) buildRun(cmd *cobra.Command, resource v1.APIResource, request map[string]interface{}) {
	cmd.Run = func(cmd *cobra.Command, args []string) {
		out, _ := json.Marshal(request)
		// Pull the name and namespace from the request so they are added to the url path
		meta := request["metadata"].(map[string]interface{})
		name := meta["name"].(*string)
		namespace := meta["namespace"].(*string)

		// Create the request
		result := builder.rest.Put().
			Prefix("apis", resource.Group, resource.Version).
			Namespace(*namespace).
			Resource(builder.resource(resource)).
			SubResource(builder.operation(resource)).
			Name(*name).
			Body(out)

		resp, err := result.DoRaw()

		fmt.Printf("URL: %v\n", result.URL().Path)
		fmt.Printf("RequestBody: %s\n", out)
		fmt.Printf("ResponseBody: %s", resp)
		if err != nil {
			fmt.Printf("Error: %v\n", resp, err)
		}
	}
}
