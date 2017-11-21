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

package kubefield

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fw "k8s.io/kubectl/pkg/framework/openapi"
	"k8s.io/kubectl/pkg/framework/resource"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

// FlagBuilder returns a new request body parsed from flag values
func (builder *CmdBuilderImpl) BuildObject(
	cmd *cobra.Command,
	r *resource.Resource,
	path []string) (func() map[string]interface{}, error) {

	visitor := newKindVisitor(cmd, r.ResourceGroupVersionKind(), path)
	r.Accept(visitor)
	return visitor.resource, nil
}

func newKindVisitor(cmd *cobra.Command, gvk schema.GroupVersionKind, path []string) *kindVisitor {
	return &kindVisitor{
		fw.PanicVisitor{},
		gvk,
		cmd,
		nil,
		path,
		map[string]*string{},
	}
}

type kindVisitor struct {
	fw.PanicVisitor
	gvk         schema.GroupVersionKind
	cmd         *cobra.Command
	resource    func() map[string]interface{}
	path        []string
	stringflags map[string]*string
}

func (visitor *kindVisitor) VisitKind(k *openapi.Kind) {
	visitor.stringflags["name"] = visitor.cmd.Flags().String("name", "", "name of the resource")
	visitor.stringflags["namespace"] = visitor.cmd.Flags().String("namespace", "default", "namespace of the resource")

	resource := map[string]interface{}{}
	resource["apiVersion"] = fmt.Sprintf("%v/%v", visitor.gvk.Group, visitor.gvk.Version)
	resource["kind"] = fmt.Sprintf("%v", visitor.gvk.Kind)
	resource["metadata"] = map[string]interface{}{
		"name":      visitor.stringflags["name"],
		"namespace": visitor.stringflags["namespace"],
	}

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

func (v *kindVisitor) newFieldVisitor(path []string) *fieldVisitor {
	return &fieldVisitor{
		v.PanicVisitor,
		v.cmd,
		nil,
		false,
		path,
		v.stringflags,
	}
}

// fieldVisitor walks the openapi schema and registers flags for primitive fields
type fieldVisitor struct {
	fw.PanicVisitor
	cmd         *cobra.Command
	resource    func() interface{}
	array       bool
	path        []string
	stringflags map[string]*string
}

// VisitKind recurses into certain fields to populate flags
func (visitor *fieldVisitor) VisitKind(k *openapi.Kind) {
	resource := map[string]interface{}{}

	// If this is the last element, provide a flag
	if len(visitor.path) <= 1 {
		value := visitor.cmd.Flags().String(visitor.path[0], "", k.Description)
		visitor.resource = func() interface{} {
			err := json.Unmarshal([]byte(*value), &resource)
			if err != nil {
				panic(err)
			}
			return resource
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
	resource[field] = fieldVisitor.resource()

	visitor.resource = func() interface{} {
		return resource
	}
}

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *fieldVisitor) VisitPrimitive(p *openapi.Primitive) {
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

func (visitor *fieldVisitor) VisitArray(p *openapi.Array) {
	resource := map[string]interface{}{}
	if len(visitor.path) <= 1 {
		value := visitor.cmd.Flags().String(visitor.path[0], "", p.Description)
		visitor.resource = func() interface{} {
			err := json.Unmarshal([]byte(*value), &resource)
			if err != nil {
				panic(err)
			}
			return []interface{}{resource}
		}
		return
	}

	if _, found := p.Extensions["x-kubernetes-patch-merge-key"]; !found {
		panic(fmt.Errorf("Cannot update items in unmergeable lists"))
	}
	mergeKey, ok := p.Extensions["x-kubernetes-patch-merge-key"].(string)
	if !ok {
		panic(fmt.Errorf("Mergekey not a string %v %T", mergeKey, mergeKey))
	}

	resource[mergeKey] = visitor.cmd.Flags().
		String(fmt.Sprintf("%s-%s", visitor.path[0], mergeKey), "", p.Description)

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

func (*fieldVisitor) VisitMap(m *openapi.Map) {
	// do nothing
}

// VisitReference traverses references
func (visitor *fieldVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}

// newFieldVisitor creates a new fieldVisitor for recursion
func (v *fieldVisitor) newFieldVisitor(path []string) *fieldVisitor {
	return &fieldVisitor{
		v.PanicVisitor,
		v.cmd,
		nil,
		false,
		path,
		v.stringflags,
	}
}
