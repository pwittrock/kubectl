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
	"fmt"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
	fw "k8s.io/kubectl/pkg/framework/openapi"
	"strings"
)

type ObjectFieldFn func(i interface{}, b openapi.BaseSchema, s openapi.Schema) interface{}

func setField(sch openapi.Schema, path []string, obj interface{}, fn ObjectFieldFn) (interface{}, error) {
	v := &objectFieldVisitor{
		fw.PanicVisitor{},
		path,
		obj,
		obj,
		fn,
		nil,
	}
	sch.Accept(v)
	return v.currentFieldValue, v.err
}

// schemaFieldVisitor walks the openapi schema and registers flags for primitive fields
type objectFieldVisitor struct {
	fw.PanicVisitor
	remainingFields   []string
	obj               interface{}
	currentFieldValue interface{}
	fieldValueFn      ObjectFieldFn
	err               error
}

// VisitPrimitive creates a new flag to populate the primitive value
func (visitor *objectFieldVisitor) VisitPrimitive(p *openapi.Primitive) {
	// Make sure we are at the end of the line
	if len(visitor.remainingFields) != 0 {
		visitor.err = fmt.Errorf(
			"primitive field value cannot have additional path elements: %v", visitor.remainingFields)
		return
	}

	// Set the value
	visitor.currentFieldValue = visitor.fieldValueFn(visitor.currentFieldValue, p.BaseSchema, p)
}

// VisitReference traverses references
func (visitor *objectFieldVisitor) VisitReference(r openapi.Reference) {
	r.SubSchema().Accept(visitor)
}

// VisitKind recurses into certain fields to populate flags
func (visitor *objectFieldVisitor) VisitKind(k *openapi.Kind) {
	// End recursion.  Call the function to set the currentFieldValue
	if len(visitor.remainingFields) == 0 {
		visitor.currentFieldValue = visitor.fieldValueFn(visitor.currentFieldValue, k.BaseSchema, k)
		return
	}

	// Lookup the next field from the schema
	fieldName := visitor.remainingFields[0]
	if _, found := k.Fields[fieldName]; !found {
		// Field does not exist
		visitor.err = fmt.Errorf("no field named %s defined for type", fieldName)
		return
	}

	// Initialize a new value for the field if empty
	if visitor.currentFieldValue == nil {
		visitor.currentFieldValue = map[string]interface{}{}
	}

	// Cast the field value to a map
	field, ok := visitor.currentFieldValue.(map[string]interface{})
	if !ok {
		visitor.err = fmt.Errorf("field expected map field type, was %T", field)
		return
	}

	// Prepare for recursion by setting the new field value and remaining fields
	visitor.currentFieldValue = field[fieldName]
	visitor.remainingFields = visitor.remainingFields[1:]

	// Recurse
	k.Fields[fieldName].Accept(visitor)

	// Update the currentFieldValue to the recursed value
	if visitor.currentFieldValue != nil {
		field[fieldName] = visitor.currentFieldValue
	} else {
		delete(field, fieldName)
	}
	visitor.currentFieldValue = field
}

// VisitKind recurses into certain fields to populate flags
func (visitor *objectFieldVisitor) VisitMap(m *openapi.Map) {
	// End recursion. Call the function to set the currentFieldValue
	if len(visitor.remainingFields) == 0 {
		visitor.currentFieldValue = visitor.fieldValueFn(visitor.currentFieldValue, m.BaseSchema, m)
		return
	}

	// Parse map element from the path
	mapElement := visitor.remainingFields[0]
	if !strings.HasSuffix(mapElement, "]") || !strings.HasPrefix(mapElement, "[") {
		// Map index was not in brackets
		visitor.err = fmt.Errorf("map index for path must be in brackets, was %s", mapElement)
		return
	}
	mapElement = strings.TrimSuffix(strings.TrimPrefix(mapElement, "["), "]")

	// Initialize a new value for the field if empty
	if visitor.currentFieldValue == nil {
		visitor.currentFieldValue = map[string]interface{}{}
	}

	// Cast the field value to a map
	field, ok := visitor.currentFieldValue.(map[string]interface{})
	if !ok {
		visitor.err = fmt.Errorf("field expected map field type, was %T", field)
		return
	}

	// Prepare for recursion by setting the new field value and remaining fields
	visitor.currentFieldValue = field[mapElement]
	visitor.remainingFields = visitor.remainingFields[1:]

	// Recurse
	m.SubType.Accept(visitor)

	// Update the currentFieldValue to the recursed value
	if visitor.currentFieldValue != nil {
		field[mapElement] = visitor.currentFieldValue
	} else {
		delete(field, mapElement)
	}
	visitor.currentFieldValue = field
}

func (visitor *objectFieldVisitor) VisitArray(p *openapi.Array) {
	// At the currentFieldValue, end recursion, set the currentFieldValue and exit
	if len(visitor.remainingFields) == 0 {
		visitor.currentFieldValue = visitor.fieldValueFn(visitor.currentFieldValue, p.BaseSchema, p)
		return
	}

	arrayElement := visitor.remainingFields[0]
	if !strings.HasSuffix(arrayElement, "]") || !strings.HasPrefix(arrayElement, "[") {
		// Array index was not in brackets
		visitor.err = fmt.Errorf("array index for path must be in brackets, was %s", arrayElement)
		return
	}
	// Strip the brackets
	arrayElement = strings.TrimSuffix(strings.TrimPrefix(arrayElement, "["), "]")

	if _, found := p.Extensions["x-kubernetes-patch-merge-key"]; !found {
		visitor.err = fmt.Errorf("Cannot update items in unmergeable lists")
		return
	}
	mergeKey, ok := p.Extensions["x-kubernetes-patch-merge-key"].(string)
	if !ok {
		visitor.err = fmt.Errorf("Mergekey not a string %v %T", mergeKey, mergeKey)
		return
	}

	// Instantiate the currentFieldValue if it is nil
	if visitor.currentFieldValue == nil {
		visitor.currentFieldValue = []interface{}{}
	}

	// Make sure it is the right type
	field, ok := visitor.currentFieldValue.([]interface{})
	if !ok {
		visitor.err = fmt.Errorf("field expected slice field type, was %T", field)
		return
	}

	// Find the matching element in the list
	var element interface{}
	for i := range field {
		v := field[i]

		// If it is a map compare the mergeKey
		e, isMap := v.(map[string]interface{})
		if isMap {
			// Found the element we are looking for
			if fmt.Sprintf("%v", e[mergeKey]) == arrayElement {
				element = e
				break
			}
		} else {
			// Otherwise compare the value
			if fmt.Sprintf("%v", v) == arrayElement {
				element = v
				break
			}
		}
	}

	visitor.currentFieldValue = element
	visitor.remainingFields = visitor.remainingFields[1:]
	p.SubType.Accept(visitor)

	// If the element was not found before, add it to the list
	if element == nil && visitor.currentFieldValue != nil {
		e, isMap := visitor.currentFieldValue.(map[string]interface{})
		// Set the mergeKey value if the element was nil before
		if isMap {
			e[mergeKey] = arrayElement
		}
		field = append(field, visitor.currentFieldValue)
	}

	visitor.currentFieldValue = field
}
