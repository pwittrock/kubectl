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
	"fmt"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

func newObjectKindVisitor(cmd *cobra.Command, name string) *objectFieldVisitor {
	return &objectFieldVisitor{
		name,
		name,
		cmd,
		nil,
		false,
		map[string]*string{},
	}
}

// objectFieldVisitor walks the openapi schema and registers flags for primitive fields
type objectFieldVisitor struct {
	parent      string
	name        string
	cmd         *cobra.Command
	field       fn
	array       bool
	stringflags map[string]*string
}

type fn func() (interface{}, bool)

// VisitKind recurses into certain fields to populate flags
func (visitor *objectFieldVisitor) VisitKind(k *openapi.Kind) {
	// The result for a Kind is a map
	resource := map[string]interface{}{}
	values := map[string]fn{}
	for k, v := range k.Fields {
		if blacklistedFields.HasAny(k) {
			continue
		}

		fv := visitor.newObjectVisitor(k)
		v.Accept(fv)
		values[k] = fv.field
	}

	visitor.field = func() (interface{}, bool) {
		changed := false
		for k, v := range values {
			if result, changed := v(); changed {
				changed = true
				resource[k] = result
			}
		}
		return resource, changed
	}
}

var blacklistedFields = sets.NewString("apiVersion", "kind", "metadata", "status", "valueFrom")

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *objectFieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	// Never set flags for blacklisted fields
	if blacklistedFields.HasAny(visitor.name) {
		return
	}

	var value interface{}
	var name string
	if visitor.parent != visitor.name {
		name = fmt.Sprintf("%s-%s", visitor.parent, visitor.name)
	} else {
		name = visitor.name
	}

	// Create a flag reference
	if !visitor.array {
		switch p.Type {
		case "integer":
			value = visitor.cmd.Flags().Int32(name, 0, p.Description)
		case "boolean":
			value = visitor.cmd.Flags().Bool(name, false, p.Description)
		case "string":
			if _, found := visitor.stringflags[name]; !found {
				visitor.stringflags[name] = visitor.cmd.Flags().String(name, "", p.Description)
			}
			value = visitor.stringflags[name]

		}
	} else {
		switch p.Type {
		case "integer":
			value = visitor.cmd.Flags().IntSlice(name, []int{}, p.Description)
		case "boolean":
			value = visitor.cmd.Flags().BoolSlice(name, []bool{}, p.Description)
		case "string":
			value = visitor.cmd.Flags().StringSlice(name, []string{}, p.Description)
		}
	}

	visitor.field = func() (interface{}, bool) {
		return value, visitor.cmd.Flag(name).Changed
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
		v.parent,
		name,
		v.cmd,
		nil,
		false,
		v.stringflags,
	}
}
