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
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CmdBuilder interface {
	// FlagBuilder returns a new request body parsed from flag values
	BuildFlags(cmd *cobra.Command, resource v1.APIResource) (map[string]interface{}, error)
	BuildCmd(resource v1.APIResource) (*cobra.Command, error)
	BuildRun(command *cobra.Command, resource v1.APIResource, request map[string]interface{})
	IsCmd(resource v1.APIResource) bool
	Seen(resource v1.APIResource) bool
}
