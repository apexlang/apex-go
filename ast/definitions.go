package ast

import "github.com/wapc/widl-go/kinds"

type (
	Definition interface {
		Node
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

// InterfaceDefinition implements Node, Definition
var _ Definition = (*InterfaceDefinition)(nil)

type InterfaceDefinition struct {
	BaseNode
	Description *StringValue `json:"description,omitempty"` // Optional
	AnnotatedNode
	Operations []*OperationDefinition `json:"operations"`
}

func NewInterfaceDefinition(loc *Location, description *StringValue, annotations []*Annotation, operations []*OperationDefinition) *InterfaceDefinition {
	return &InterfaceDefinition{
		BaseNode:      BaseNode{kinds.InterfaceDefinition, loc},
		Description:   description,
		AnnotatedNode: AnnotatedNode{annotations},
		Operations:    operations,
	}
}

// RoleDefinition implements Node, Definition
var _ Definition = (*RoleDefinition)(nil)

type RoleDefinition struct {
	BaseNode
	Name        *Name        `json:"name"`
	Description *StringValue `json:"description,omitempty"` // Optional
	AnnotatedNode
	Operations []*OperationDefinition `json:"operations"`
}

func NewRoleDefinition(loc *Location, name *Name, description *StringValue, annotations []*Annotation, operations []*OperationDefinition) *RoleDefinition {
	return &RoleDefinition{
		BaseNode:      BaseNode{kinds.RoleDefinition, loc},
		Name:          name,
		Description:   description,
		Operations:    operations,
		AnnotatedNode: AnnotatedNode{annotations},
	}
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

// UnionDefinition implements Node, Definition
var _ Definition = (*UnionDefinition)(nil)

type UnionDefinition struct {
	BaseNode
	Name        *Name                  `json:"name"`
	Description *StringValue           `json:"description,omitempty"` // Optional
	Parameters  []*ParameterDefinition `json:"parameters"`
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
