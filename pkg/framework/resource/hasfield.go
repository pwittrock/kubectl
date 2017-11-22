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

package resource

import (
	openapi "k8s.io/kube-openapi/pkg/util/proto"
	fw "k8s.io/kubectl/pkg/framework/openapi"
)

func hasField(sch openapi.Schema, path []string) bool {
	v := &fieldVisitor{fw.PanicVisitor{}, path, false}
	sch.Accept(v)
	return v.found
}

// fieldVisitor walks the openapi schema and registers flags for primitive fields
type fieldVisitor struct {
	fw.PanicVisitor
	path  []string
	found bool
}

// VisitKind recurses into certain fields to populate flags
func (visitor *fieldVisitor) VisitKind(k *openapi.Kind) {
	if len(visitor.path) == 0 {
		// Field found
		visitor.found = true
		return
	}

	field := visitor.path[0]
	if _, found := k.Fields[field]; !found {
		// Field not found
		return
	}

	// Eat a path element and recurse
	visitor.path = visitor.path[1:]
	k.Fields[field].Accept(visitor)
}

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *fieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	// At the leaf nodes
	visitor.found = len(visitor.path) == 0
}

func (visitor *fieldVisitor) VisitArray(p *openapi.Array) {
	if len(visitor.path) == 0 {
		visitor.found = true
		return
	}
	// List do not eat a path element
	p.SubType.Accept(visitor)

}

func (visitor *fieldVisitor) VisitMap(m *openapi.Map) {
	if len(visitor.path) == 0 {
		visitor.found = true
		return
	}
	// Maps eat a path element
	visitor.path = visitor.path[1:]
	m.SubType.Accept(visitor)
}

// VisitReference traverses references
func (visitor *fieldVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}
