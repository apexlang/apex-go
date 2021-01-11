package kinds

const (
	// Name
	Name = "Name"

	// Document
	Document            = "Document"
	OperationDefinition = "OperationDefinition"
	VariableDefinition  = "VariableDefinition"
	Variable            = "Variable"
	SelectionSet        = "SelectionSet"
	Field               = "Field"
	Argument            = "Argument"

	// // Fragments
	// FragmentSpread     = "FragmentSpread"
	// InlineFragment     = "InlineFragment"
	// FragmentDefinition = "FragmentDefinition"

	// Values
	IntValue     = "IntValue"
	FloatValue   = "FloatValue"
	StringValue  = "StringValue"
	BooleanValue = "BooleanValue"
	EnumValue    = "EnumValue"
	ListValue    = "ListValue"
	MapValue     = "MapValue"
	ObjectValue  = "ObjectValue"
	ObjectField  = "ObjectField"

	// Annotations
	Annotation = "Annotation"

	// Types
	Named    = "Named" // previously NamedType
	List     = "List"  // previously ListType
	Map      = "Map"
	Optional = "Optional"

	// Type System Definitions
	SchemaDefinition        = "SchemaDefinition"
	OperationTypeDefinition = "OperationTypeDefinition"

	// Types Definitions
	InputValueDefinition = "InputValueDefinition"

	// Types Extensions
	NamespaceDefinition = "NamespaceDefinition"
	ScalarDefinition    = "ScalarDefinition" // previously ScalarTypeDefinition
	ObjectDefinition    = "ObjectDefinition" // previously ObjectTypeDefinition
	FieldDefinition     = "FieldDefinition"
	InterfaceDefinition = "InterfaceDefinition" // previously InterfaceTypeDefinition
	UnionDefinition     = "UnionDefinition"     // previously UnionTypeDefinition
	EnumDefinition      = "EnumDefinition"      // previously EnumTypeDefinition
	EnumValueDefinition = "EnumValueDefinition"

	// Annotation Definitions
	AnnotationDefinition = "AnnotationDefinition"
)
