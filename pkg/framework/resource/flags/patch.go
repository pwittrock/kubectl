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
	//"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fw "k8s.io/kubectl/pkg/framework/openapi"
	//"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
)

func newPatchKindVisitor(cmd *cobra.Command, gvk schema.GroupVersionKind, path []string) *patchKindVisitor {
	return &patchKindVisitor{
		fw.PanicVisitor{},
		gvk,
		cmd,
		nil,
		path,
		map[string]*string{},
	}
}

type patchKindVisitor struct {
	fw.PanicVisitor
	gvk         schema.GroupVersionKind
	cmd         *cobra.Command
	resource    func() map[string]interface{}
	path        []string
	stringflags map[string]*string
}

func (visitor *patchKindVisitor) VisitKind(k *openapi.Kind) {
	//visitor.stringflags["name"] = visitor.cmd.Flags().String("name", "", "name of the resource")
	//visitor.stringflags["namespace"] = visitor.cmd.Flags().String("namespace", "default", "namespace of the resource")

	resource := map[string]interface{}{}
	//resource["apiVersion"] = fmt.Sprintf("%v/%v", visitor.gvk.Group, visitor.gvk.Version)
	//resource["kind"] = fmt.Sprintf("%v", visitor.gvk.Kind)
	//resource["metadata"] = map[string]interface{}{
	//	"name":      visitor.stringflags["name"],
	//	"namespace": visitor.stringflags["namespace"],
	//}

	if len(visitor.path) == 0 {
		panic(fmt.Errorf("path must have length greater than 0: %v", visitor.path))
	}

	// Lookup the first field in the path
	field := visitor.path[0]
	if _, found := k.Fields[field]; !found {
		panic(fmt.Errorf("field %v not found", visitor.path[0]))
	}

	// Visit this field
	fieldVisitor := visitor.newFieldVisitor(visitor.path)
	k.Fields[field].Accept(fieldVisitor)

	visitor.resource = func() map[string]interface{} {
		resource[field] = fieldVisitor.resource()
		return resource
	}
}

func (v *patchKindVisitor) newFieldVisitor(path []string) *patchFieldVisitor {
	return &patchFieldVisitor{
		v.PanicVisitor,
		v.cmd,
		nil,
		false,
		path,
		v.stringflags,
	}
}

// patchFieldVisitor walks the openapi schema and registers flags for primitive fields
type patchFieldVisitor struct {
	fw.PanicVisitor
	cmd         *cobra.Command
	resource    func() interface{}
	array       bool
	path        []string
	stringflags map[string]*string
}

// VisitKind recurses into certain fields to populate flags
func (visitor *patchFieldVisitor) VisitKind(k *openapi.Kind) {
	resource := map[string]interface{}{}

	// If this is the last element, provide a flag
	if len(visitor.path) <= 1 {
		ov := newObjectKindVisitor(visitor.cmd, visitor.path[0])
		k.Accept(ov)
		visitor.resource = func() interface{} {
			value, _ := ov.field()
			return value
		}
		return
	}

	// Otherwise recurse down the path
	field := visitor.path[1]
	if _, found := k.Fields[field]; !found {
		panic(fmt.Errorf("field %v not found", visitor.path[1]))
	}

	// Visit this field
	fieldVisitor := visitor.newFieldVisitor(visitor.path[1:])
	k.Fields[field].Accept(fieldVisitor)

	visitor.resource = func() interface{} {
		resource[field] = fieldVisitor.resource()
		return resource
	}
}

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *patchFieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	// Create a flag reference
	var value interface{}
	if !visitor.array {
		switch p.Type {
		case "integer":
			value = visitor.cmd.Flags().Int32(visitor.path[0], 0, p.Description)
		case "boolean":
			value = visitor.cmd.Flags().Bool(visitor.path[0], false, p.Description)
		case "string":
			if _, found := visitor.stringflags[visitor.path[0]]; !found {
				visitor.stringflags[visitor.path[0]] = visitor.cmd.Flags().String(visitor.path[0], "", p.Description)
			}
			value = visitor.stringflags[visitor.path[0]]
		}
	} else {
		switch p.Type {
		case "integer":
			value = visitor.cmd.Flags().IntSlice(visitor.path[0], []int{}, p.Description)
		case "boolean":
			value = visitor.cmd.Flags().BoolSlice(visitor.path[0], []bool{}, p.Description)
		case "string":
			value = visitor.cmd.Flags().StringSlice(visitor.path[0], []string{}, p.Description)
		}
	}

	// Return the parsed value
	visitor.resource = func() interface{} {
		return value
	}
}

func (visitor *patchFieldVisitor) VisitArray(p *openapi.Array) {
	resource := map[string]interface{}{}

	if _, found := p.Extensions["x-kubernetes-patch-merge-key"]; !found {
		panic(fmt.Errorf("Cannot update items in unmergeable lists"))
	}
	mergeKey, ok := p.Extensions["x-kubernetes-patch-merge-key"].(string)
	if !ok {
		panic(fmt.Errorf("Mergekey not a string %v %T", mergeKey, mergeKey))
	}

	if len(visitor.path) > 1 {
		resource[mergeKey] = visitor.cmd.Flags().
			String(fmt.Sprintf("%s-%s", visitor.path[0], mergeKey), "", p.Description)
	}

	fv := visitor.newFieldVisitor(visitor.path)
	fv.array = true
	p.SubType.Accept(fv)

	visitor.resource = func() interface{} {
		result := fv.resource()

		casted, ok := result.(map[string]interface{})
		if ok {
			for key, value := range casted {
				resource[key] = value
			}
			return []interface{}{resource}
		}

		name := visitor.path[1]
		resource[name] = result
		return []interface{}{resource}
	}
}

func (visitor *patchFieldVisitor) VisitMap(m *openapi.Map) {
	resource := map[string]interface{}{}

	if len(visitor.path) == 1 {
		fv := visitor.newFieldVisitor(visitor.path)
		m.SubType.Accept(fv)

		// If this is the last element, provide a flag
		key := visitor.cmd.Flags().String(fmt.Sprintf("%s-key", visitor.path[0]), "", m.Description)
		visitor.resource = func() interface{} {
			resource[*key] = fv.resource()
			return resource
		}
		return
	}

	fv := visitor.newFieldVisitor(visitor.path[1:])
	m.SubType.Accept(fv)

	visitor.resource = func() interface{} {
		result := fv.resource()

		casted, ok := result.(map[string]interface{})
		if ok {
			for key, value := range casted {
				resource[key] = value
			}
			return resource
		}

		name := visitor.path[1]
		resource[name] = result
		return resource
	}
}

// VisitReference traverses references
func (visitor *patchFieldVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}

// newFieldVisitor creates a new patchFieldVisitor for recursion
func (v *patchFieldVisitor) newFieldVisitor(path []string) *patchFieldVisitor {
	return &patchFieldVisitor{
		v.PanicVisitor,
		v.cmd,
		nil,
		false,
		path,
		v.stringflags,
	}
}
