package openapi

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

// FlagBuilder returns a new request body parsed from flag values
func (builder *cmdBuilderImpl) buildFlags(cmd *cobra.Command, resource *SubResource) (map[string]interface{}, error) {
	gvk := schema.GroupVersionKind{
		resource.apiGroupVersion.Group,
		resource.apiGroupVersion.Version,
		resource.resource.Kind,
	}

	// Build a request body from flags
	visitor := newKindVisitor(cmd, gvk)
	resource.openapiSchema.Accept(visitor)
	return visitor.resource, nil
}

func newKindVisitor(cmd *cobra.Command, gvk schema.GroupVersionKind) *kindVisitor {
	resource := map[string]interface{}{}
	resource["apiVersion"] = fmt.Sprintf("%v/%v", gvk.Group, gvk.Version)
	resource["kind"] = fmt.Sprintf("%v", gvk.Kind)

	return &kindVisitor{
		emptyVisitor{},
		cmd,
		resource,
		map[string]*string{},
	}
}

type kindVisitor struct {
	emptyVisitor
	cmd         *cobra.Command
	resource    map[string]interface{}
	stringflags map[string]*string
}

func (v *kindVisitor) getRequest() interface{} {
	return v.resource
}

func (visitor *kindVisitor) VisitKind(k *openapi.Kind) {
	visitor.stringflags["name"] = visitor.cmd.Flags().String("name", "", "name of the resource")
	visitor.stringflags["namespace"] = visitor.cmd.Flags().String("namespace", "default", "namespace of the resource")

	visitor.resource["metadata"] = map[string]interface{}{
		"name":      visitor.stringflags["name"],
		"namespace": visitor.stringflags["namespace"],
	}

	for k, v := range k.Fields {
		if blacklistedFields.Has(k) {
			continue
		}
		fv := visitor.newFieldVisitor(k)
		v.Accept(fv)
		if fv.field != nil {
			visitor.resource[k] = fv.field
		}
	}
}

func (v *kindVisitor) newFieldVisitor(name string) *fieldVisitor {
	return &fieldVisitor{
		v.emptyVisitor,
		name,
		v.cmd,
		map[string]interface{}{},
		false,
		v.stringflags,
	}
}

// fieldVisitor walks the openapi schema and registers flags for primitive fields
type fieldVisitor struct {
	emptyVisitor
	name        string
	cmd         *cobra.Command
	field       interface{}
	array       bool
	stringflags map[string]*string
}

var whitelistedFields = sets.NewString("spec", "rollbackTo")

// VisitKind recurses into certain fields to populate flags
func (visitor *fieldVisitor) VisitKind(k *openapi.Kind) {
	// Only recurse for whitelisted fields
	if !whitelistedFields.HasAny(visitor.name) {
		return
	}

	// The result for a Kind is a map
	resource := map[string]interface{}{}
	visitor.field = resource

	for k, v := range k.Fields {
		fv := visitor.newFieldVisitor(k)
		v.Accept(fv)
		if fv.field != nil {
			resource[k] = fv.field
		}
	}
}

var blacklistedFields = sets.NewString("apiVersion", "kind", "metadata", "status")

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *fieldVisitor) VisitPrimitive(p *openapi.Primitive) {
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

func (visitor *fieldVisitor) VisitArray(p *openapi.Array) {
	// Never set flags for blacklisted fields
	if blacklistedFields.HasAny(visitor.name) {
		return
	}

	fv := visitor.newFieldVisitor(visitor.name)
	fv.array = true
	p.SubType.Accept(fv)
	if fv.field != nil {
		visitor.field = fv.field
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
func (v *fieldVisitor) newFieldVisitor(name string) *fieldVisitor {
	return &fieldVisitor{
		v.emptyVisitor,
		name,
		v.cmd,
		nil,
		false,
		v.stringflags,
	}
}

// emptyVisitor is a base implementation for composition that panics for all calls.
type emptyVisitor struct{}

func (*emptyVisitor) VisitArray(a *openapi.Array) { panic(fmt.Errorf("Unexpected array call %+v", a)) }
func (*emptyVisitor) VisitMap(m *openapi.Map)     { panic(fmt.Errorf("Unexpected map call %+v", m)) }
func (*emptyVisitor) VisitPrimitive(p *openapi.Primitive) {
	panic(fmt.Errorf("Unexpected primitive call %+v", p))
}
func (*emptyVisitor) VisitKind(k *openapi.Kind) { panic(fmt.Errorf("Unexpected kind call %+v", k)) }
func (*emptyVisitor) VisitReference(r openapi.Reference) {
	panic(fmt.Errorf("Unexpected reference call %+v", r))
}
