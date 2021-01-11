package ast

import (
	"github.com/wapc/widl-go/language/kinds"
)

// DescribableNode are nodes that have descriptions associated with them.
type DescribableNode interface {
	GetDescription() *StringValue
}

type TypeDefinition interface {
	DescribableNode
	GetKind() string
	GetLoc() *Location
}

var _ TypeDefinition = (*NamespaceDefinition)(nil)
var _ TypeDefinition = (*ScalarDefinition)(nil)
var _ TypeDefinition = (*ObjectDefinition)(nil)
var _ TypeDefinition = (*InterfaceDefinition)(nil)
var _ TypeDefinition = (*UnionDefinition)(nil)
var _ TypeDefinition = (*EnumDefinition)(nil)

type TypeSystemDefinition interface {
	GetKind() string
	GetLoc() *Location
}

var _ TypeSystemDefinition = (*SchemaDefinition)(nil)
var _ TypeSystemDefinition = (TypeDefinition)(nil)
var _ TypeSystemDefinition = (*AnnotationDefinition)(nil)

// SchemaDefinition implements Node, Definition
type SchemaDefinition struct {
	Kind           string
	Loc            *Location
	Annotations    []*Annotation
	OperationTypes []*OperationTypeDefinition
}

func NewSchemaDefinition(def *SchemaDefinition) *SchemaDefinition {
	if def == nil {
		def = &SchemaDefinition{}
	}
	return &SchemaDefinition{
		Kind:           kinds.SchemaDefinition,
		Loc:            def.Loc,
		Annotations:    def.Annotations,
		OperationTypes: def.OperationTypes,
	}
}

func (def *SchemaDefinition) GetKind() string {
	return def.Kind
}

func (def *SchemaDefinition) GetLoc() *Location {
	return def.Loc
}

// OperationTypeDefinition implements Node, Definition
type OperationTypeDefinition struct {
	Kind      string
	Loc       *Location
	Operation string
	Type      *Named
}

func NewOperationTypeDefinition(def *OperationTypeDefinition) *OperationTypeDefinition {
	if def == nil {
		def = &OperationTypeDefinition{}
	}
	return &OperationTypeDefinition{
		Kind:      kinds.OperationTypeDefinition,
		Loc:       def.Loc,
		Operation: def.Operation,
		Type:      def.Type,
	}
}

func (def *OperationTypeDefinition) GetKind() string {
	return def.Kind
}

func (def *OperationTypeDefinition) GetLoc() *Location {
	return def.Loc
}

// NamespaceDefinition implements Node, Definition
type NamespaceDefinition struct {
	Kind        string
	Loc         *Location
	Description *StringValue
	Name        *Name
	Annotations []*Annotation
}

func NewNamespaceDefinition(def *NamespaceDefinition) *NamespaceDefinition {
	if def == nil {
		def = &NamespaceDefinition{}
	}
	return &NamespaceDefinition{
		Kind:        kinds.NamespaceDefinition,
		Loc:         def.Loc,
		Description: def.Description,
		Name:        def.Name,
		Annotations: def.Annotations,
	}
}

func (def *NamespaceDefinition) GetKind() string {
	return def.Kind
}

func (def *NamespaceDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *NamespaceDefinition) GetName() *Name {
	return def.Name
}

func (def *NamespaceDefinition) GetDescription() *StringValue {
	return def.Description
}

// ScalarDefinition implements Node, Definition
type ScalarDefinition struct {
	Kind        string
	Loc         *Location
	Description *StringValue
	Name        *Name
	Type        Type
	Annotations []*Annotation
}

func NewScalarDefinition(def *ScalarDefinition) *ScalarDefinition {
	if def == nil {
		def = &ScalarDefinition{}
	}
	return &ScalarDefinition{
		Kind:        kinds.ScalarDefinition,
		Loc:         def.Loc,
		Description: def.Description,
		Name:        def.Name,
		Type:        def.Type,
		Annotations: def.Annotations,
	}
}

func (def *ScalarDefinition) GetKind() string {
	return def.Kind
}

func (def *ScalarDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *ScalarDefinition) GetName() *Name {
	return def.Name
}

func (def *ScalarDefinition) GetOperation() string {
	return ""
}

func (def *ScalarDefinition) GetDescription() *StringValue {
	return def.Description
}

// ObjectDefinition implements Node, Definition
type ObjectDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Interfaces  []*Named
	Annotations []*Annotation
	Fields      []*FieldDefinition
}

func NewObjectDefinition(def *ObjectDefinition) *ObjectDefinition {
	if def == nil {
		def = &ObjectDefinition{}
	}
	return &ObjectDefinition{
		Kind:        kinds.ObjectDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Interfaces:  def.Interfaces,
		Annotations: def.Annotations,
		Fields:      def.Fields,
	}
}

func (def *ObjectDefinition) GetKind() string {
	return def.Kind
}

func (def *ObjectDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *ObjectDefinition) GetName() *Name {
	return def.Name
}

func (def *ObjectDefinition) GetDescription() *StringValue {
	return def.Description
}

// OperationDefinition implements Node
type OperationDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Arguments   []*InputValueDefinition
	Input       *InputValueDefinition
	Type        Type
	Annotations []*Annotation
}

func NewOperationDefinition(def *OperationDefinition) *OperationDefinition {
	if def == nil {
		def = &OperationDefinition{}
	}
	return &OperationDefinition{
		Kind:        kinds.OperationDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Arguments:   def.Arguments,
		Input:       def.Input,
		Type:        def.Type,
		Annotations: def.Annotations,
	}
}

func (def *OperationDefinition) GetKind() string {
	return def.Kind
}

func (def *OperationDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *OperationDefinition) GetDescription() *StringValue {
	return def.Description
}

// FieldDefinition implements Node
type FieldDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Type        Type
	Default     Value
	Annotations []*Annotation
}

func NewFieldDefinition(def *FieldDefinition) *FieldDefinition {
	if def == nil {
		def = &FieldDefinition{}
	}
	return &FieldDefinition{
		Kind:        kinds.FieldDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Type:        def.Type,
		Default:     def.Default,
		Annotations: def.Annotations,
	}
}

func (def *FieldDefinition) GetKind() string {
	return def.Kind
}

func (def *FieldDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *FieldDefinition) GetDescription() *StringValue {
	return def.Description
}

// InputValueDefinition implements Node
type InputValueDefinition struct {
	Kind         string
	Loc          *Location
	Name         *Name
	Description  *StringValue
	Type         Type
	DefaultValue Value
	Annotations  []*Annotation
}

func NewInputValueDefinition(def *InputValueDefinition) *InputValueDefinition {
	if def == nil {
		def = &InputValueDefinition{}
	}
	return &InputValueDefinition{
		Kind:         kinds.InputValueDefinition,
		Loc:          def.Loc,
		Name:         def.Name,
		Description:  def.Description,
		Type:         def.Type,
		DefaultValue: def.DefaultValue,
		Annotations:  def.Annotations,
	}
}

func (def *InputValueDefinition) GetKind() string {
	return def.Kind
}

func (def *InputValueDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InputValueDefinition) GetDescription() *StringValue {
	return def.Description
}

// InterfaceDefinition implements Node, Definition
type InterfaceDefinition struct {
	Kind string
	Loc  *Location
	//Name        *Name
	Description *StringValue
	Annotations []*Annotation
	Operations  []*OperationDefinition
}

func NewInterfaceDefinition(def *InterfaceDefinition) *InterfaceDefinition {
	if def == nil {
		def = &InterfaceDefinition{}
	}
	return &InterfaceDefinition{
		Kind: kinds.InterfaceDefinition,
		Loc:  def.Loc,
		//Name:        def.Name,
		Description: def.Description,
		Annotations: def.Annotations,
		Operations:  def.Operations,
	}
}

func (def *InterfaceDefinition) GetKind() string {
	return def.Kind
}

func (def *InterfaceDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *InterfaceDefinition) GetDescription() *StringValue {
	return def.Description
}

// UnionDefinition implements Node, Definition
type UnionDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Annotations []*Annotation
	Types       []*Named
}

func NewUnionDefinition(def *UnionDefinition) *UnionDefinition {
	if def == nil {
		def = &UnionDefinition{}
	}
	return &UnionDefinition{
		Kind:        kinds.UnionDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Annotations: def.Annotations,
		Types:       def.Types,
	}
}

func (def *UnionDefinition) GetKind() string {
	return def.Kind
}

func (def *UnionDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *UnionDefinition) GetName() *Name {
	return def.Name
}

func (def *UnionDefinition) GetDescription() *StringValue {
	return def.Description
}

// EnumDefinition implements Node, Definition
type EnumDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Annotations []*Annotation
	Values      []*EnumValueDefinition
}

func NewEnumDefinition(def *EnumDefinition) *EnumDefinition {
	if def == nil {
		def = &EnumDefinition{}
	}
	return &EnumDefinition{
		Kind:        kinds.EnumDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Annotations: def.Annotations,
		Values:      def.Values,
	}
}

func (def *EnumDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *EnumDefinition) GetName() *Name {
	return def.Name
}

func (def *EnumDefinition) GetDescription() *StringValue {
	return def.Description
}

// EnumValueDefinition implements Node, Definition
type EnumValueDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Annotations []*Annotation
}

func NewEnumValueDefinition(def *EnumValueDefinition) *EnumValueDefinition {
	if def == nil {
		def = &EnumValueDefinition{}
	}
	return &EnumValueDefinition{
		Kind:        kinds.EnumValueDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Annotations: def.Annotations,
	}
}

func (def *EnumValueDefinition) GetKind() string {
	return def.Kind
}

func (def *EnumValueDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *EnumValueDefinition) GetDescription() *StringValue {
	return def.Description
}
