package openapi

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
)

// FlagBuilder returns a new request body parsed from flag values
func (builder *cmdBuilderImpl) BuildFlags(cmd *cobra.Command, resource v1.APIResource) (map[string]interface{}, error) {
	gvk := schema.GroupVersionKind{resource.Group, resource.Version, resource.Kind}

	apiSchema := builder.resources.LookupResource(gvk)
	if apiSchema == nil {
		return nil, fmt.Errorf("No openapi definition found for %+v", gvk)
	}

	// Build a request body from flags
	visitor := newKindVisitor(cmd, gvk)
	apiSchema.Accept(visitor)
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
	}
}

type kindVisitor struct {
	emptyVisitor
	cmd      *cobra.Command
	resource map[string]interface{}
}

func (v *kindVisitor) getRequest() interface{} {
	return v.resource
}

func (visitor *kindVisitor) VisitKind(k *openapi.Kind) {
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
	}
}

// fieldVisitor walks the openapi schema and registers flags for primitive fields
type fieldVisitor struct {
	emptyVisitor
	name  string
	cmd   *cobra.Command
	field interface{}
}

var whitelistedFields = sets.NewString("spec")

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
	switch p.Type {
	case "integer":
		visitor.field = visitor.cmd.Flags().Int32(visitor.name, 0, p.Description)
	case "boolean":
		visitor.field = visitor.cmd.Flags().Bool(visitor.name, false, p.Description)
	case "string":
		visitor.field = visitor.cmd.Flags().String(visitor.name, "", p.Description)
	}
}

func (visitor *fieldVisitor) VisitArray(p *openapi.Array) {
	// Never set flags for blacklisted fields
	if blacklistedFields.HasAny(visitor.name) {
		return
	}

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
