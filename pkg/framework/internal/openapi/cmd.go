package openapi

import (
	"fmt"

	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

// NewFlagBuilder returns a new request builder
func NewCmdBuilder(resources openapi.Resources,
	discovery discovery.DiscoveryInterface,
	rest rest.Interface,
	apiGroup, apiVersion string) CmdBuilder {
	return &cmdBuilderImpl{
		resources,
		discovery,
		rest,
		map[string]sets.String{},
		apiGroup,
		apiVersion,
	}
}

type cmdBuilderImpl struct {
	resources  openapi.Resources
	discovery  discovery.DiscoveryInterface
	rest       rest.Interface
	seen       map[string]sets.String
	apiGroup   string
	apiVersion string
}

func (builder *cmdBuilderImpl) buildCmd(name string, resource *v1.APIResource, versions []schema.GroupVersion) (*cobra.Command, error) {
	gvk := schema.GroupVersionKind{resource.Group, resource.Version, resource.Kind}
	if builder.resources.LookupResource(gvk) == nil {
		return nil, fmt.Errorf("No openapi definition found for %+v", gvk)
	}

	kind, operation, err := builder.resourceOperation(resource.Name)
	if err != nil {
		return nil, err
	}

	versionsList := []string{}
	d := false
	var group, version string
	for _, v := range versions {
		if !d {
			group = v.Group
			version = v.Version
			versionsList = append(versionsList, fmt.Sprintf("\t%s/%s (default)", v.Group, v.Version))
			d = true
		} else {
			versionsList = append(versionsList, fmt.Sprintf("\t%s/%s", v.Group, v.Version))
		}
	}

	cmd := &cobra.Command{
		Use: fmt.Sprintf("%v", kind),
		Example: fmt.Sprintf("kubecurl --api-group %s --api-version %s %s %s %s --name foo",
			group, version, name, operation, builder.resource(resource)),
		Short: fmt.Sprintf("%s %s", operation, kind),
		Long: fmt.Sprintf(`Supported group/versions:
%s

(set the group and version to use with with --api-group and --api-version *must be provided before any subcommands*)`,
			strings.Join(versionsList, "\n")),
	}
	return cmd, nil
}
