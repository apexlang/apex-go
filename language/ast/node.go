package ast

type Node interface {
	GetKind() string
	GetLoc() *Location
}

// The list of all possible AST node graphql.
// Ensure that all node types implements Node interface
var _ Node = (*Name)(nil)
var _ Node = (*Document)(nil)

var _ Node = (*IntValue)(nil)
var _ Node = (*FloatValue)(nil)
var _ Node = (*StringValue)(nil)
var _ Node = (*BooleanValue)(nil)
var _ Node = (*EnumValue)(nil)
var _ Node = (*ListValue)(nil)
var _ Node = (*ObjectValue)(nil)
var _ Node = (*ObjectField)(nil)
var _ Node = (*Annotation)(nil)
var _ Node = (*Named)(nil)
var _ Node = (*List)(nil)
var _ Node = (*Map)(nil)
var _ Node = (*Optional)(nil)
var _ Node = (*OperationTypeDefinition)(nil)
var _ Node = (*ObjectDefinition)(nil)
var _ Node = (*FieldDefinition)(nil)
var _ Node = (*InputValueDefinition)(nil)
var _ Node = (*InterfaceDefinition)(nil)
var _ Node = (*UnionDefinition)(nil)
var _ Node = (*EnumDefinition)(nil)
var _ Node = (*EnumValueDefinition)(nil)
var _ Node = (*AnnotationDefinition)(nil)
