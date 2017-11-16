package openapi

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
    "k8s.io/apimachinery/pkg/util/sets"
)

// NewFlagBuilder returns a new request builder
func NewCmdBuilder(resources openapi.Resources) CmdBuilder {
	return &cmdBuilderImpl{resources, map[string]sets.String{}}
}

type cmdBuilderImpl struct {
	resources openapi.Resources
    seen map[string]sets.String{}
}

func (builder *cmdBuilderImpl) BuildCmd(resource v1.APIResource) (*cobra.Command, error) {
	gvk := schema.GroupVersionKind{resource.Group, resource.Version, resource.Kind}
	if builder.resources.LookupResource(gvk) == nil {
		return nil, fmt.Errorf("No openapi definition found for %+v", gvk)
	}

	parts := strings.Split(resource.Name, "/")
	kind := parts[0]
	operation := parts[1]

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%v", kind),
		Short: fmt.Sprintf("%v command for %v/%v/%v", operation, resource.Group, resource.Version, kind),
	}
	return cmd, nil
}
