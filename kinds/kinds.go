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
	Stream   Kind = "Stream"

	// Definitions
	NamespaceDefinition Kind = "NamespaceDefinition"
	ImportDefinition    Kind = "ImportDefinition"
	AliasDefinition     Kind = "AliasDefinition"
	InterfaceDefinition Kind = "InterfaceDefinition"
	OperationDefinition Kind = "OperationDefinition"
	ParameterDefinition Kind = "ParameterDefinition"
	TypeDefinition      Kind = "TypeDefinition"
	FieldDefinition     Kind = "FieldDefinition"
	UnionDefinition     Kind = "UnionDefinition"
	EnumDefinition      Kind = "EnumDefinition"
	EnumValueDefinition Kind = "EnumValueDefinition"
	DirectiveDefinition Kind = "DirectiveDefinition"
)
