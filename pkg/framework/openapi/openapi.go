package openapi

import (
	"k8s.io/kubectl/pkg/framework/internal/inject"
	fw "k8s.io/kubectl/pkg/framework/internal/openapi"
)

// TODO: switch to dependency injection
// This must be a singleton
var factory = inject.NewFactory()

func NewCmdBuilder() fw.CmdBuilder {
	return fw.NewCmdBuilder(factory.GetResources(), factory.GetDiscovery())
}
