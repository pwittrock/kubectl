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
	//"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	//fw "k8s.io/kubectl/pkg/framework/openapi"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

func newObjectKindVisitor(cmd *cobra.Command, name string) *objectFieldVisitor {
	return &objectFieldVisitor{
		name,
		cmd,
		nil,
		false,
		map[string]*string{},
	}
}

// objectFieldVisitor walks the openapi schema and registers flags for primitive fields
type objectFieldVisitor struct {
	name        string
	cmd         *cobra.Command
	field       interface{}
	array       bool
	stringflags map[string]*string
}

var whitelistedFields = sets.NewString("spec", "rollbackTo")

// VisitKind recurses into certain fields to populate flags
func (visitor *objectFieldVisitor) VisitKind(k *openapi.Kind) {
	// Only recurse for whitelisted fields
	if !whitelistedFields.HasAny(visitor.name) {
		return
	}

	// The result for a Kind is a map
	resource := map[string]interface{}{}
	visitor.field = resource

	for k, v := range k.Fields {
		fv := visitor.newObjectVisitor(k)
		v.Accept(fv)
		if fv.field != nil {
			resource[k] = fv.field
		}
	}
}

var blacklistedFields = sets.NewString("apiVersion", "kind", "metadata", "status")

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *objectFieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	// Never set flags for blacklisted fields
	if blacklistedFields.HasAny(visitor.name) {
		return
	}

	// Create a flag reference
	if !visitor.array {
		switch p.Type {
		case "integer":
			visitor.field = visitor.cmd.Flags().Int32(visitor.name, 0, p.Description)
		case "boolean":
			visitor.field = visitor.cmd.Flags().Bool(visitor.name, false, p.Description)
		case "string":
			if _, found := visitor.stringflags[visitor.name]; !found {
				visitor.stringflags[visitor.name] = visitor.cmd.Flags().String(visitor.name, "", p.Description)
			}
			visitor.field = visitor.stringflags[visitor.name]
		}
	} else {
		switch p.Type {
		case "integer":
			visitor.field = visitor.cmd.Flags().IntSlice(visitor.name, []int{}, p.Description)
		case "boolean":
			visitor.field = visitor.cmd.Flags().BoolSlice(visitor.name, []bool{}, p.Description)
		case "string":
			visitor.field = visitor.cmd.Flags().StringSlice(visitor.name, []string{}, p.Description)
		}
	}
}

func (visitor *objectFieldVisitor) VisitArray(p *openapi.Array) {
	// Never set flags for blacklisted fields
	if blacklistedFields.HasAny(visitor.name) {
		return
	}

	fv := visitor.newObjectVisitor(visitor.name)
	fv.array = true
	p.SubType.Accept(fv)
	if fv.field != nil {
		visitor.field = fv.field
	}
}

func (*objectFieldVisitor) VisitMap(m *openapi.Map) {
	// do nothing
}

// VisitReference traverses references
func (visitor *objectFieldVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}

// newFieldVisitor creates a new patchFieldVisitor for recursion
func (v *objectFieldVisitor) newObjectVisitor(name string) *objectFieldVisitor {
	return &objectFieldVisitor{
		name,
		v.cmd,
		nil,
		false,
		v.stringflags,
	}
}
