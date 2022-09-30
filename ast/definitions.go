/*
Copyright 2022 The Apex Authors.

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

package ast

import "github.com/apexlang/apex-go/kinds"

type (
	Definition interface {
		Node
		Accept(context Context, visitor Visitor)
	}

	Annotated interface {
		Annotation(name string, callback func(annotation *Annotation)) *Annotation
	}

	AnnotatedNode struct {
		Annotations []*Annotation `json:"annotations"`
	}
)

func (a *AnnotatedNode) Annotation(name string) *Annotation {
	for _, annotation := range a.Annotations {
		if annotation.Name.Value == name {
			return annotation
		}
	}
	return nil
}

// NamespaceDefinition implements Node, Definition
var _ Definition = (*NamespaceDefinition)(nil)

type NamespaceDefinition struct {
	BaseNode
	Name        *Name        `json:" name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	AnnotatedNode
}

func NewNamespaceDefinition(loc *Location, name *Name, description *StringValue, annotations []*Annotation) *NamespaceDefinition {
	return &NamespaceDefinition{
		BaseNode:      BaseNode{kinds.NamespaceDefinition, loc},
		Name:          name,
		Description:   description,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *NamespaceDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitNamespace(context)
	VisitAnnotations(context, visitor, d.Annotations)
}

// AliasDefinition implements Node, Definition
var _ Definition = (*AliasDefinition)(nil)

type AliasDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	Type        Type         `json:"type"`
	AnnotatedNode
}

func NewAliasDefinition(loc *Location, name *Name, description *StringValue, t Type, annotations []*Annotation) *AliasDefinition {
	return &AliasDefinition{
		BaseNode:      BaseNode{kinds.AliasDefinition, loc},
		Name:          name,
		Description:   description,
		Type:          t,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *AliasDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitAlias(context)
	VisitAnnotations(context, visitor, d.Annotations)
}

// ImportDefinition implements Node, Definition
var _ Definition = (*ImportDefinition)(nil)

type ImportDefinition struct {
	BaseNode
	Description *StringValue  `json:"description,omitempty"` // Optional
	All         bool          `json:"all"`
	Names       []*ImportName `json:"names"`
	From        *StringValue  `json:"from"`
	AnnotatedNode
}

func NewImportDefinition(loc *Location, description *StringValue, all bool, names []*ImportName, from *StringValue, annotations []*Annotation) *ImportDefinition {
	return &ImportDefinition{
		BaseNode:      BaseNode{kinds.ImportDefinition, loc},
		Description:   description,
		All:           all,
		Names:         names,
		From:          from,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *ImportDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitImport(context)
	VisitAnnotations(context, visitor, d.Annotations)
}

// TypeDefinition implements Node, Definition
var _ Definition = (*TypeDefinition)(nil)

type TypeDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	Interfaces  []*Named     `json:"interfaces,omitempty"`
	AnnotatedNode
	Fields []*FieldDefinition `json:"fields"`
}

func NewTypeDefinition(loc *Location, name *Name, description *StringValue, interfaces []*Named, annotations []*Annotation, fields []*FieldDefinition) *TypeDefinition {
	return &TypeDefinition{
		BaseNode:      BaseNode{kinds.TypeDefinition, loc},
		Name:          name,
		Description:   description,
		Interfaces:    interfaces,
		AnnotatedNode: AnnotatedNode{annotations},
		Fields:        fields,
	}
}

func (d *TypeDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitTypeBefore(context)
	visitor.VisitType(context)

	c := context
	c.Fields = context.Type.Fields
	visitor.VisitTypeFieldsBefore(c)
	for i, field := range c.Fields {
		c.FieldIndex = i
		c.Field = field
		field.Accept(c, visitor)
	}
	visitor.VisitTypeFieldsAfter(c)

	VisitAnnotations(context, visitor, d.Annotations)
	visitor.VisitTypeAfter(context)
}

type ValuedDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	Type        Type         `json:"type"`
	Default     Value        `json:"default,omitempty"` // Optional
	AnnotatedNode
}

// FieldDefinition implements Node, Definition
var _ Definition = (*FieldDefinition)(nil)

type FieldDefinition ValuedDefinition

func NewFieldDefinition(loc *Location, name *Name, description *StringValue, t Type, defaultValue Value, annotations []*Annotation) *FieldDefinition {
	return &FieldDefinition{
		BaseNode:      BaseNode{kinds.FieldDefinition, loc},
		Name:          name,
		Description:   description,
		Type:          t,
		Default:       defaultValue,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *FieldDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitTypeField(context)
	VisitAnnotations(context, visitor, d.Annotations)
}

// RoleDefinition implements Node, Definition
var _ Definition = (*InterfaceDefinition)(nil)

type InterfaceDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	AnnotatedNode
	Operations []*OperationDefinition `json:"operations"`
}

func NewInterfaceDefinition(loc *Location, name *Name, description *StringValue, annotations []*Annotation, operations []*OperationDefinition) *InterfaceDefinition {
	return &InterfaceDefinition{
		BaseNode:      BaseNode{kinds.InterfaceDefinition, loc},
		Name:          name,
		Description:   description,
		Operations:    operations,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *InterfaceDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitInterfaceBefore(context)
	visitor.VisitInterface(context)

	c := context
	c.Operations = c.Interface.Operations
	visitor.VisitOperationsBefore(c)
	for _, oper := range c.Operations {
		c.Operation = oper
		oper.Accept(c, visitor)
	}
	visitor.VisitOperationsAfter(c)

	VisitAnnotations(context, visitor, d.Annotations)
	visitor.VisitInterfaceAfter(context)
}

// OperationDefinition implements Node, Definition
var _ Definition = (*OperationDefinition)(nil)

type OperationDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	Type        Type         `json:"type"`
	AnnotatedNode
	Unary      bool                   `json:"unary"`
	Parameters []*ParameterDefinition `json:"parameters"`
}

func NewOperationDefinition(loc *Location, name *Name, description *StringValue, ttype Type, annotations []*Annotation, unary bool, parameters []*ParameterDefinition) *OperationDefinition {
	return &OperationDefinition{
		BaseNode:      BaseNode{kinds.OperationDefinition, loc},
		Name:          name,
		Description:   description,
		Type:          ttype,
		AnnotatedNode: AnnotatedNode{annotations},
		Unary:         unary,
		Parameters:    parameters,
	}
}

func (o *OperationDefinition) IsUnary() bool {
	return o.Unary && len(o.Parameters) == 1
}

func (d *OperationDefinition) Accept(context Context, visitor Visitor) {
	if context.Interface != nil {
		visitor.VisitOperationBefore(context)
		visitor.VisitOperation(context)
	} else {
		visitor.VisitFunctionBefore(context)
		visitor.VisitFunction(context)
	}

	c := context
	c.Parameters = c.Operation.Parameters
	visitor.VisitParametersBefore(context)
	for _, param := range c.Parameters {
		c.Parameter = param
		param.Accept(c, visitor)
	}
	visitor.VisitParametersAfter(c)

	VisitAnnotations(context, visitor, d.Annotations)
	if context.Interface != nil {
		visitor.VisitOperationAfter(context)
	} else {
		visitor.VisitFunctionAfter(context)
	}
}

// ParameterDefinition implements Node, Definition
var _ Definition = (*ParameterDefinition)(nil)

type ParameterDefinition ValuedDefinition

func NewParameterDefinition(loc *Location, name *Name, description *StringValue, t Type, defaultValue Value, annotations []*Annotation) *ParameterDefinition {
	return &ParameterDefinition{
		BaseNode:      BaseNode{kinds.ParameterDefinition, loc},
		Name:          name,
		Description:   description,
		Type:          t,
		Default:       defaultValue,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *ParameterDefinition) Accept(context Context, visitor Visitor) {
	if context.Operation != nil {
		visitor.VisitParameter(context)
	} else if context.Directive != nil {
		visitor.VisitDirectiveParameter(context)
	}
	VisitAnnotations(context, visitor, d.Annotations)
}

// UnionDefinition implements Node, Definition
var _ Definition = (*UnionDefinition)(nil)

type UnionDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	AnnotatedNode
	Types []Type `json:"types"`
}

func NewUnionDefinition(loc *Location, name *Name, description *StringValue, annotations []*Annotation, types []Type) *UnionDefinition {
	return &UnionDefinition{
		BaseNode:      BaseNode{kinds.UnionDefinition, loc},
		Name:          name,
		Description:   description,
		AnnotatedNode: AnnotatedNode{annotations},
		Types:         types,
	}
}

func (d *UnionDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitUnion(context)
	VisitAnnotations(context, visitor, d.Annotations)
}

// EnumDefinition implements Node, Definition
var _ Definition = (*EnumDefinition)(nil)

type EnumDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	AnnotatedNode
	Values []*EnumValueDefinition `json:"values"`
}

func NewEnumDefinition(loc *Location, name *Name, description *StringValue, annotations []*Annotation, values []*EnumValueDefinition) *EnumDefinition {
	return &EnumDefinition{
		BaseNode:      BaseNode{kinds.EnumDefinition, loc},
		Name:          name,
		Description:   description,
		AnnotatedNode: AnnotatedNode{annotations},
		Values:        values,
	}
}

func (d *EnumDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitEnumBefore(context)
	visitor.VisitEnum(context)

	c := context
	c.EnumValues = c.Enum.Values
	visitor.VisitEnumValuesBefore(c)
	for _, enumValue := range c.EnumValues {
		c.EnumValue = enumValue
		enumValue.Accept(c, visitor)
	}
	visitor.VisitEnumValuesAfter(c)

	VisitAnnotations(context, visitor, d.Annotations)
	visitor.VisitEnumAfter(context)
}

// EnumValueDefinition implements Node, Definition
var _ Definition = (*EnumValueDefinition)(nil)

type EnumValueDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	Index       *IntValue    `json:"index"`
	Display     *StringValue `json:"display,omitempty"` // Optional
	AnnotatedNode
}

func NewEnumValueDefinition(loc *Location, name *Name, description *StringValue, index *IntValue, display *StringValue, annotations []*Annotation) *EnumValueDefinition {
	return &EnumValueDefinition{
		BaseNode:      BaseNode{kinds.EnumValueDefinition, loc},
		Name:          name,
		Description:   description,
		Index:         index,
		Display:       display,
		AnnotatedNode: AnnotatedNode{annotations},
	}
}

func (d *EnumValueDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitEnumValue(context)
	VisitAnnotations(context, visitor, d.Annotations)
}

// DirectiveDefinition implements Node, Definition
var _ Definition = (*DirectiveDefinition)(nil)

type DirectiveDefinition struct {
	BaseNode
	Name        *Name                  `json:"name"`
	Description *StringValue           `json:"description,omitempty"` // Optional
	Parameters  []*ParameterDefinition `json:"parameters"`
	Locations   []*Name                `json:"locations"`
	Requires    []*DirectiveRequire    `json:"requires,omitempty"` // Optional
}

func NewDirectiveDefinition(loc *Location, name *Name, description *StringValue, parameters []*ParameterDefinition, locations []*Name, requires []*DirectiveRequire) *DirectiveDefinition {
	return &DirectiveDefinition{
		BaseNode:    BaseNode{kinds.DirectiveDefinition, loc},
		Name:        name,
		Description: description,
		Parameters:  parameters,
		Locations:   locations,
		Requires:    requires,
	}
}

func (d *DirectiveDefinition) Accept(context Context, visitor Visitor) {
	visitor.VisitDirectiveBefore(context)
	visitor.VisitDirective(context)

	c := context
	c.Parameters = c.Directive.Parameters
	visitor.VisitDirectiveParametersBefore(c)
	for _, param := range c.Parameters {
		c.Parameter = param
		param.Accept(c, visitor)
	}
	visitor.VisitDirectiveParametersAfter(c)

	visitor.VisitDirectiveAfter(context)
}

func VisitAnnotations(
	context Context,
	visitor Visitor,
	annotations []*Annotation,
) {
	if len(annotations) == 0 {
		return
	}

	visitor.VisitAnnotationsBefore(context)
	for _, a := range annotations {
		c := context
		c.Annotation = a
		visitor.VisitAnnotationBefore(c)
		a.Accept(c, visitor)
		visitor.VisitAnnotationAfter(c)
	}
	visitor.VisitAnnotationsAfter(context)
}
