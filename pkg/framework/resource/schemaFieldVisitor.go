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
	found := false
	sch.Accept(&schemaFieldVisitor{
		fw.PanicVisitor{},
		path,
		func(openapi.BaseSchema, openapi.Schema) { found = true },
	})
	return found
}

// schemaFieldVisitor walks the openapi schema and registers flags for primitive fields
type schemaFieldVisitor struct {
	fw.PanicVisitor
	path    []string
	fieldFn func(openapi.BaseSchema, openapi.Schema)
}

// VisitKind recurses into certain fields to populate flags
func (visitor *schemaFieldVisitor) VisitKind(k *openapi.Kind) {
	if len(visitor.path) == 0 {
		// Field found
		visitor.fieldFn(k.BaseSchema, k)
		return
	}

	field := visitor.path[0]
	if _, found := k.Fields[field]; !found {
		// Field not found
		return
	}

	// Eat a remainingFields element and recurse
	visitor.path = visitor.path[1:]
	k.Fields[field].Accept(visitor)
}

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *schemaFieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	// At the leaf nodes
	if len(visitor.path) == 0 {
		visitor.fieldFn(p.BaseSchema, p)
	}
}

func (visitor *schemaFieldVisitor) VisitArray(p *openapi.Array) {
	if len(visitor.path) == 0 {
		visitor.fieldFn(p.BaseSchema, p)
		return
	}
	// List do not eat a remainingFields element
	p.SubType.Accept(visitor)

}

func (visitor *schemaFieldVisitor) VisitMap(m *openapi.Map) {
	if len(visitor.path) == 0 {
		visitor.fieldFn(m.BaseSchema, m)
		return
	}
	// Maps eat a remainingFields element
	visitor.path = visitor.path[1:]
	m.SubType.Accept(visitor)
}

// VisitReference traverses references
func (visitor *schemaFieldVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}
