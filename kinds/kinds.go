package kinds

type Kind string

const (
	// Nodes
	Document         Kind = "Document"
	Name             Kind = "Name"
	Annotation       Kind = "Annotation"
	Argument         Kind = "Argument"
	DirectiveRequire Kind = "DirectiveRequire"
	ImportName       Kind = "ImportName"

	// Values
	IntValue     Kind = "IntValue"
	FloatValue   Kind = "FloatValue"
	StringValue  Kind = "StringValue"
	BooleanValue Kind = "BooleanValue"
	EnumValue    Kind = "EnumValue"
	ListValue    Kind = "ListValue"
	MapValue     Kind = "MapValue"
	ObjectValue  Kind = "ObjectValue"
	ObjectField  Kind = "ObjectField"

	// Types
	Named    Kind = "Named"
	ListType Kind = "ListType"
	MapType  Kind = "MapType"
	Optional Kind = "Optional"

	// Definitions
	NamespaceDefinition Kind = "NamespaceDefinition"
	ImportDefinition    Kind = "ImportDefinition"
	AliasDefinition     Kind = "AliasDefinition"
	InterfaceDefinition Kind = "InterfaceDefinition"
	RoleDefinition      Kind = "RoleDefinition"
	OperationDefinition Kind = "OperationDefinition"
	ParameterDefinition Kind = "ParameterDefinition"
	TypeDefinition      Kind = "TypeDefinition"
	FieldDefinition     Kind = "FieldDefinition"
	UnionDefinition     Kind = "UnionDefinition"
	EnumDefinition      Kind = "EnumDefinition"
	EnumValueDefinition Kind = "EnumValueDefinition"
	DirectiveDefinition Kind = "DirectiveDefinition"
)
