package openapi

import (
	fw "k8s.io/kubectl/pkg/framework/internal/openapi"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

func NewCmdBuilder(resources openapi.Resources) fw.CmdBuilder {
	return fw.NewCmdBuilder(resources)
}
