package openapi

import (
	"fmt"
	"k8s.io/kubectl/pkg/framework/internal/inject"
	fw "k8s.io/kubectl/pkg/framework/internal/openapi"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
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

type PanicVisitor struct{}

func (p *PanicVisitor) VisitArray(a *openapi.Array)         { p.panic(a) }
func (p *PanicVisitor) VisitMap(a *openapi.Map)             { p.panic(a) }
func (p *PanicVisitor) VisitPrimitive(a *openapi.Primitive) { p.panic(a) }
func (p *PanicVisitor) VisitKind(a *openapi.Kind)           { p.panic(a) }
func (p *PanicVisitor) VisitReference(a openapi.Reference)  { p.panic(a) }

func (*PanicVisitor) panic(t interface{}) { panic(fmt.Errorf("Unexpected visitor call %T", t)) }

type NoOpVisitor struct{}

func (*NoOpVisitor) VisitArray(a *openapi.Array)         {}
func (*NoOpVisitor) VisitMap(m *openapi.Map)             {}
func (*NoOpVisitor) VisitPrimitive(p *openapi.Primitive) {}
func (*NoOpVisitor) VisitKind(k *openapi.Kind)           {}
func (*NoOpVisitor) VisitReference(r openapi.Reference)  {}
