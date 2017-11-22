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

package flags

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/framework"
	"k8s.io/kubectl/pkg/framework/resource"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

func NewFlagBuilder() *FlagBuilderImpl {
	return &FlagBuilderImpl{
		framework.Factory().GetResources(),
		framework.Factory().GetDiscovery(),
		framework.Factory().GetRest(),
		map[string]sets.String{},
		framework.Factory().GetApiGroup(),
		framework.Factory().GetApiVersion(),
	}
}

type FlagBuilderImpl struct {
	resources  openapi.Resources
	discovery  discovery.DiscoveryInterface
	rest       rest.Interface
	seen       map[string]sets.String
	apiGroup   string
	apiVersion string
}

// FlagBuilder returns a new request body parsed from flag values
func (builder *FlagBuilderImpl) BuildObject(
	cmd *cobra.Command,
	r *resource.Resource,
	path []string) (func() map[string]interface{}, error) {

	visitor := newPatchKindVisitor(cmd, r.ResourceGroupVersionKind(), path)
	r.Accept(visitor)
	return visitor.resource, nil
}
