package openapi

import (
	"k8s.io/kubectl/pkg/framework/internal/inject"
	fw "k8s.io/kubectl/pkg/framework/internal/openapi"
)

// TODO: switch to dependency injection
// This must be a singleton

func NewCmdBuilder() fw.CmdBuilder {
	return fw.NewCmdBuilder(
		inject.FactorySingleton.GetResources(),
		inject.FactorySingleton.GetDiscovery(),
		inject.FactorySingleton.GetRest(),
		inject.FactorySingleton.GetApiGroup(),
		inject.FactorySingleton.GetApiVersion(),
	)
}
