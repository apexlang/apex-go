/*
Copyright 2024 The Apex Authors.

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

package model

import (
	"errors"

	"github.com/apexlang/apex-go/ast"
)

func Convert(doc *ast.Document) (*Namespace, []error) {
	var converter Converter
	return converter.Convert(doc)
}

var scalars = map[string]Scalar{
	"string":   ScalarString,
	"bool":     ScalarBool,
	"i8":       ScalarI8,
	"i16":      ScalarI16,
	"i32":      ScalarI32,
	"i64":      ScalarI64,
	"u8":       ScalarU8,
	"u16":      ScalarU16,
	"u32":      ScalarU32,
	"u64":      ScalarU64,
	"f32":      ScalarF32,
	"f64":      ScalarF64,
	"bytes":    ScalarBytes,
	"datetime": ScalarDatetime,
	"any":      ScalarAny,
	"raw":      ScalarRaw,
}

type Converter struct {
	_ns         *ast.NamespaceDefinition
	_imports    []*ast.ImportDefinition
	_directives []*ast.DirectiveDefinition
	_aliases    []*ast.AliasDefinition
	_unions     []*ast.UnionDefinition
	_functions  []*ast.OperationDefinition
	_types      []*ast.TypeDefinition
	_enums      []*ast.EnumDefinition
	_interfaces []*ast.InterfaceDefinition

	named map[string]Named

	errors []error
}

func (c *Converter) Convert(doc *ast.Document) (*Namespace, []error) {
	c.named = make(map[string]Named)
	for _, def := range doc.Definitions {
		switch t := def.(type) {
		case *ast.NamespaceDefinition:
			c._ns = t
		case *ast.ImportDefinition:
			c._imports = append(c._imports, t)
		case *ast.DirectiveDefinition:
			c._directives = append(c._directives, t)
		case *ast.AliasDefinition:
			c._aliases = append(c._aliases, t)
			c.named[t.Name.Value] = Named{
				Kind: KindAlias,
				Name: t.Name.Value,
			}
		case *ast.UnionDefinition:
			c._unions = append(c._unions, t)
			c.named[t.Name.Value] = Named{
				Kind: KindUnion,
				Name: t.Name.Value,
			}
		case *ast.EnumDefinition:
			c._enums = append(c._enums, t)
			c.named[t.Name.Value] = Named{
				Kind: KindEnum,
				Name: t.Name.Value,
			}
		case *ast.OperationDefinition:
			c._functions = append(c._functions, t)
		case *ast.TypeDefinition:
			c._types = append(c._types, t)
			c.named[t.Name.Value] = Named{
				Kind: KindType,
				Name: t.Name.Value,
			}
		case *ast.InterfaceDefinition:
			c._interfaces = append(c._interfaces, t)
		}
	}

	if c._ns == nil {
		return nil, []error{errors.New("no namespace found")}
	}

	ns := Namespace{
		Description: stringValuePtr(c._ns.Description),
		Name:        c._ns.Name.Value,
		Annotations: c.convertAnnotations(c._ns.Annotations),
		Imports:     c.convertImports(c._imports),
		Directives:  c.convertDirectives(c._directives),
		Aliases:     c.convertAliases(c._aliases),
		Unions:      c.convertUnions(c._unions),
		Functions:   c.convertOperations(c._functions),
		Types:       c.convertTypes(c._types),
		Interfaces:  c.convertInterfaces(c._interfaces),
	}

	if len(c.errors) > 0 {
		return nil, c.errors
	}

	return &ns, nil
}

func (c *Converter) convertInterfaces(items []*ast.InterfaceDefinition) []Interface {
	if len(items) == 0 {
		return nil
	}
	s := make([]Interface, len(items))
	for i, item := range items {
		s[i] = Interface{
			Description: stringValuePtr(item.Description),
			Name:        item.Name.Value,
			Operations:  c.convertOperations(item.Operations),
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertTypes(items []*ast.TypeDefinition) []Type {
	if len(items) == 0 {
		return nil
	}
	s := make([]Type, len(items))
	for i, item := range items {
		s[i] = Type{
			Description: stringValuePtr(item.Description),
			Name:        item.Name.Value,
			Fields:      c.convertFields(item.Fields),
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertFields(items []*ast.FieldDefinition) []Field {
	if len(items) == 0 {
		return nil
	}
	s := make([]Field, len(items))
	for i, item := range items {
		s[i] = Field{
			Description:  stringValuePtr(item.Description),
			Name:         item.Name.Value,
			Type:         c.convertTypeRef(item.Type),
			DefaultValue: c.convertValuePtr(item.Default),
			Annotations:  c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertOperations(items []*ast.OperationDefinition) []Operation {
	if len(items) == 0 {
		return nil
	}
	s := make([]Operation, len(items))
	for i, item := range items {
		parameters := c.convertParameters(item.Parameters)
		var unary *Parameter
		if item.Unary && len(parameters) == 1 {
			unary = &parameters[0]
			parameters = nil
		}
		s[i] = Operation{
			Description: stringValuePtr(item.Description),
			Name:        item.Name.Value,
			Unary:       unary,
			Parameters:  parameters,
			Returns:     c.convertTypeRefPtr(item.Type),
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertAliases(items []*ast.AliasDefinition) []Alias {
	if len(items) == 0 {
		return nil
	}
	s := make([]Alias, len(items))
	for i, item := range items {
		s[i] = Alias{
			Description: stringValuePtr(item.Description),
			Name:        item.Name.Value,
			Type:        c.convertTypeRef(item.Type),
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertUnions(items []*ast.UnionDefinition) []Union {
	if len(items) == 0 {
		return nil
	}
	s := make([]Union, len(items))
	for i, item := range items {
		s[i] = Union{
			Description: stringValuePtr(item.Description),
			Name:        item.Name.Value,
			Members:     c.convertUnionMembers(item.Members),
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertUnionMembers(items []*ast.UnionMemberDefinition) []UnionMember {
	if len(items) == 0 {
		return nil
	}
	s := make([]UnionMember, len(items))
	for i, item := range items {
		s[i] = UnionMember{
			Description: stringValuePtr(item.Description),
			Type:        c.convertTypeRef(item.Type),
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertTypeRefPtr(t ast.Type) *TypeRef {
	if named, ok := t.(*ast.Named); ok {
		if named.Name.Value == "void" {
			return nil
		}
	}
	ref := c.convertTypeRef(t)
	return &ref
}

func (c *Converter) convertTypeRef(t ast.Type) TypeRef {
	var m TypeRef
	switch v := t.(type) {
	case *ast.Named:
		if s, ok := scalars[v.Name.Value]; ok {
			m.Scalar = &s
		} else if named, ok := c.named[v.Name.Value]; ok {
			m.Named = &named
		} else {
			c.errors = append(c.errors, errors.New("unknown type "+v.Name.Value))
		}
	case *ast.ListType:
		m.List = &List{
			Type: c.convertTypeRef(v.Type),
		}
	case *ast.MapType:
		m.Map = &Map{
			KeyType:   c.convertTypeRef(v.KeyType),
			ValueType: c.convertTypeRef(v.ValueType),
		}
	case *ast.Optional:
		m.Optional = &Optional{
			Type: c.convertTypeRef(v.Type),
		}
	case *ast.Stream:
		m.Stream = &Stream{
			Type: c.convertTypeRef(v.Type),
		}
	}
	return m
}

func (c *Converter) convertDirectives(items []*ast.DirectiveDefinition) []Directive {
	if len(items) == 0 {
		return nil
	}
	s := make([]Directive, len(items))
	for i, item := range items {
		s[i] = Directive{
			Description: stringValuePtr(item.Description),
			Name:        item.Name.Value,
			Parameters:  c.convertParameters(item.Parameters),
			Locations:   c.convertDirectiveLocations(item.Locations),
			Require:     c.convertRequires(item.Requires),
		}
	}
	return s
}

func (c *Converter) convertParameters(items []*ast.ParameterDefinition) []Parameter {
	if len(items) == 0 {
		return nil
	}
	s := make([]Parameter, len(items))
	for i, item := range items {
		s[i] = Parameter{
			Description:  stringValuePtr(item.Description),
			Name:         item.Name.Value,
			Type:         c.convertTypeRef(item.Type),
			DefaultValue: c.convertValuePtr(item.Default),
			Annotations:  c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertDirectiveLocations(items []*ast.Name) []DirectiveLocation {
	if len(items) == 0 {
		return nil
	}
	s := make([]DirectiveLocation, len(items))
	for i, item := range items {
		s[i].FromString(item.Value)
	}
	return s
}

func (c *Converter) convertRequires(items []*ast.DirectiveRequire) []DirectiveRequire {
	if len(items) == 0 {
		return nil
	}
	s := make([]DirectiveRequire, len(items))
	for i, item := range items {
		s[i] = DirectiveRequire{
			Directive: item.Directive.Value,
			Locations: c.convertDirectiveLocations(item.Locations),
		}
	}
	return s
}

func (c *Converter) convertImports(items []*ast.ImportDefinition) []Import {
	if len(items) == 0 {
		return nil
	}
	s := make([]Import, len(items))
	for i, item := range items {
		s[i] = Import{
			Description: stringValuePtr(item.Description),
			All:         item.All,
			Names:       c.convertImportNames(item.Names),
			From:        item.From.Value,
			Annotations: c.convertAnnotations(item.Annotations),
		}
	}
	return s
}

func (c *Converter) convertImportNames(items []*ast.ImportName) []ImportRef {
	if len(items) == 0 {
		return nil
	}
	s := make([]ImportRef, len(items))
	for i, item := range items {
		s[i] = ImportRef{
			Name: item.Name.Value,
			As:   nameValuePtr(item.Alias),
		}
	}
	return s
}

func (c *Converter) convertAnnotations(items []*ast.Annotation) []Annotation {
	if len(items) == 0 {
		return nil
	}
	s := make([]Annotation, len(items))
	for i, item := range items {
		s[i] = Annotation{
			Name:      item.Name.Value,
			Arguments: c.convertArguments(item.Arguments),
		}
	}
	return s
}

func (c *Converter) convertArguments(items []*ast.Argument) []Argument {
	if len(items) == 0 {
		return nil
	}
	s := make([]Argument, len(items))
	for i, item := range items {
		s[i] = Argument{
			Name:  item.Name.Value,
			Value: c.convertValue(item.Value),
		}
	}
	return s
}

func (c *Converter) convertValuePtr(v ast.Value) *Value {
	if v == nil {
		return nil
	}
	r := c.convertValue(v)
	return &r
}

func (c *Converter) convertValue(v ast.Value) Value {
	switch t := v.(type) {
	case *ast.BooleanValue:
		return Value{
			Bool: &t.Value,
		}
	case *ast.IntValue:
		i64 := int64(t.Value)
		return Value{
			I64: &i64,
		}
	case *ast.StringValue:
		return Value{
			String: &t.Value,
		}
	case *ast.FloatValue:
		return Value{
			F64: &t.Value,
		}
	case *ast.EnumValue:
		return Value{
			Reference: &Reference{
				Name: t.Value,
			},
		}
	case *ast.ListValue:
		values := make([]Value, len(t.Values))
		for i, val := range t.Values {
			values[i] = c.convertValue(val)
		}
		return Value{
			ListValue: &ListValue{
				Values: values,
			},
		}
	case *ast.ObjectValue:
		fields := make([]ObjectField, len(t.Fields))
		for i, field := range t.Fields {
			fields[i] = ObjectField{
				Name:  field.Name.Value,
				Value: c.convertValue(field.Value),
			}
		}
		return Value{
			ObjectValue: &ObjectValue{
				Fields: fields,
			},
		}
	}
	return Value{}
}

func stringValuePtr(value *ast.StringValue) *string {
	if value == nil {
		return nil
	}
	return &value.Value
}

func nameValuePtr(value *ast.Name) *string {
	if value == nil {
		return nil
	}
	return &value.Value
}
