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

package ast

type Context struct {
	Document   *Document
	Namespace  *NamespaceDefinition
	Imports    []*ImportDefinition
	Import     *ImportDefinition
	Directives []*DirectiveDefinition
	Directive  *DirectiveDefinition
	Aliases    []*AliasDefinition
	Alias      *AliasDefinition
	Unions     []*UnionDefinition
	Union      *UnionDefinition
	Functions  []*OperationDefinition
	Function   *OperationDefinition
	Types      []*TypeDefinition
	Type       *TypeDefinition
	Fields     []*FieldDefinition
	Field      *FieldDefinition
	FieldIndex int
	Enums      []*EnumDefinition
	Enum       *EnumDefinition
	EnumValues []*EnumValueDefinition
	EnumValue  *EnumValueDefinition
	Interfaces []*InterfaceDefinition
	Interface  *InterfaceDefinition
	Operations []*OperationDefinition
	Operation  *OperationDefinition
	Parameters []*ParameterDefinition
	Parameter  *ParameterDefinition

	Annotation *Annotation

	Named map[string]Definition

	persistent *contextPersistent
}

type contextPersistent struct {
	Errors []error
}

func NewContext(doc *Document) Context {
	named := make(map[string]Definition)
	c := Context{
		Document:   doc,
		Named:      named,
		persistent: &contextPersistent{},
	}
	for _, def := range doc.Definitions {
		switch t := def.(type) {
		case *NamespaceDefinition:
			c.Namespace = t
		case *ImportDefinition:
			c.Imports = append(c.Imports, t)
		case *DirectiveDefinition:
			c.Directives = append(c.Directives, t)
		case *AliasDefinition:
			c.Aliases = append(c.Aliases, t)
			named[t.Name.Value] = t
		case *UnionDefinition:
			c.Unions = append(c.Unions, t)
			named[t.Name.Value] = t
		case *EnumDefinition:
			c.Enums = append(c.Enums, t)
			named[t.Name.Value] = t
		case *OperationDefinition:
			c.Functions = append(c.Functions, t)
		case *TypeDefinition:
			c.Types = append(c.Types, t)
			named[t.Name.Value] = t
		case *InterfaceDefinition:
			c.Interfaces = append(c.Interfaces, t)
		}
	}
	return c
}

func (c *Context) ReportError(err error) {
	c.persistent.Errors = append(c.persistent.Errors, err)
}

func (c *Context) Errors() []error {
	return c.persistent.Errors
}

type Visitor interface {
	VisitDocumentBefore(context Context)
	VisitNamespace(context Context)

	VisitImportsBefore(context Context)
	VisitImport(context Context)
	VisitImportsAfter(context Context)

	VisitDirectivesBefore(context Context)
	VisitDirectiveBefore(context Context)
	VisitDirective(context Context)
	VisitDirectiveParametersBefore(context Context)
	VisitDirectiveParameter(context Context)
	VisitDirectiveParametersAfter(context Context)
	VisitDirectiveAfter(context Context)
	VisitDirectivesAfter(context Context)

	VisitAliasesBefore(context Context)
	VisitAliasBefore(context Context)
	VisitAlias(context Context)
	VisitAliasAfter(context Context)
	VisitAliasesAfter(context Context)

	VisitAllOperationsBefore(context Context)
	VisitFunctionsBefore(context Context)
	VisitFunctionBefore(context Context)
	VisitFunction(context Context)
	VisitFunctionAfter(context Context)
	VisitFunctionsAfter(context Context)
	VisitInterfacesBefore(context Context)
	VisitInterfaceBefore(context Context)
	VisitInterface(context Context)
	VisitOperationsBefore(context Context)
	VisitOperationBefore(context Context)
	VisitOperation(context Context)
	VisitParametersBefore(context Context)
	VisitParameter(context Context)
	VisitParametersAfter(context Context)
	VisitOperationAfter(context Context)
	VisitOperationsAfter(context Context)
	VisitInterfaceAfter(context Context)
	VisitInterfacesAfter(context Context)
	VisitAllOperationsAfter(context Context)

	VisitTypesBefore(context Context)
	VisitTypeBefore(context Context)
	VisitType(context Context)
	VisitTypeFieldsBefore(context Context)
	VisitTypeField(context Context)
	VisitTypeFieldsAfter(context Context)
	VisitTypeAfter(context Context)
	VisitTypesAfter(context Context)

	VisitEnumsBefore(context Context)
	VisitEnumBefore(context Context)
	VisitEnum(context Context)
	VisitEnumValuesBefore(context Context)
	VisitEnumValue(context Context)
	VisitEnumValuesAfter(context Context)
	VisitEnumAfter(context Context)
	VisitEnumsAfter(context Context)

	VisitUnionsBefore(context Context)
	VisitUnion(context Context)
	VisitUnionMembersBefore(context Context)
	VisitUnionMember(context Context)
	VisitUnionMembersAfter(context Context)
	VisitUnionsAfter(context Context)

	VisitAnnotationsBefore(context Context)
	VisitAnnotationBefore(context Context)
	VisitAnnotation(context Context)
	VisitAnnotationArgumentsBefore(context Context)
	VisitAnnotationArgument(context Context)
	VisitAnnotationArgumentsAfter(context Context)
	VisitAnnotationAfter(context Context)
	VisitAnnotationsAfter(context Context)

	VisitDocumentAfter(context Context)
}

type BaseVisitor struct{}

var _ = Visitor((*BaseVisitor)(nil))

func (b *BaseVisitor) VisitDocumentBefore(context Context) {}
func (b *BaseVisitor) VisitNamespace(context Context)      {}

func (b *BaseVisitor) VisitImportsBefore(context Context) {}
func (b *BaseVisitor) VisitImport(context Context)        {}
func (b *BaseVisitor) VisitImportsAfter(context Context)  {}

func (b *BaseVisitor) VisitDirectivesBefore(context Context)          {}
func (b *BaseVisitor) VisitDirectiveBefore(context Context)           {}
func (b *BaseVisitor) VisitDirective(context Context)                 {}
func (b *BaseVisitor) VisitDirectiveParametersBefore(context Context) {}
func (b *BaseVisitor) VisitDirectiveParameter(context Context)        {}
func (b *BaseVisitor) VisitDirectiveParametersAfter(context Context)  {}
func (b *BaseVisitor) VisitDirectiveAfter(context Context)            {}
func (b *BaseVisitor) VisitDirectivesAfter(context Context)           {}

func (b *BaseVisitor) VisitAliasesBefore(context Context) {}
func (b *BaseVisitor) VisitAliasBefore(context Context)   {}
func (b *BaseVisitor) VisitAlias(context Context)         {}
func (b *BaseVisitor) VisitAliasAfter(context Context)    {}
func (b *BaseVisitor) VisitAliasesAfter(context Context)  {}

func (b *BaseVisitor) VisitAllOperationsBefore(context Context) {}
func (b *BaseVisitor) VisitFunctionsBefore(context Context)     {}
func (b *BaseVisitor) VisitFunctionBefore(context Context)      {}
func (b *BaseVisitor) VisitFunction(context Context)            {}
func (b *BaseVisitor) VisitFunctionAfter(context Context)       {}
func (b *BaseVisitor) VisitFunctionsAfter(context Context)      {}
func (b *BaseVisitor) VisitInterfacesBefore(context Context)    {}
func (b *BaseVisitor) VisitInterfaceBefore(context Context)     {}
func (b *BaseVisitor) VisitInterface(context Context)           {}
func (b *BaseVisitor) VisitOperationsBefore(context Context)    {}
func (b *BaseVisitor) VisitOperationBefore(context Context)     {}
func (b *BaseVisitor) VisitOperation(context Context)           {}
func (b *BaseVisitor) VisitParametersBefore(context Context)    {}
func (b *BaseVisitor) VisitParameter(context Context)           {}
func (b *BaseVisitor) VisitParametersAfter(context Context)     {}
func (b *BaseVisitor) VisitOperationAfter(context Context)      {}
func (b *BaseVisitor) VisitOperationsAfter(context Context)     {}
func (b *BaseVisitor) VisitInterfaceAfter(context Context)      {}
func (b *BaseVisitor) VisitInterfacesAfter(context Context)     {}
func (b *BaseVisitor) VisitAllOperationsAfter(context Context)  {}

func (b *BaseVisitor) VisitTypesBefore(context Context)      {}
func (b *BaseVisitor) VisitTypeBefore(context Context)       {}
func (b *BaseVisitor) VisitType(context Context)             {}
func (b *BaseVisitor) VisitTypeFieldsBefore(context Context) {}
func (b *BaseVisitor) VisitTypeField(context Context)        {}
func (b *BaseVisitor) VisitTypeFieldsAfter(context Context)  {}
func (b *BaseVisitor) VisitTypeAfter(context Context)        {}
func (b *BaseVisitor) VisitTypesAfter(context Context)       {}

func (b *BaseVisitor) VisitEnumsBefore(context Context)      {}
func (b *BaseVisitor) VisitEnumBefore(context Context)       {}
func (b *BaseVisitor) VisitEnum(context Context)             {}
func (b *BaseVisitor) VisitEnumValuesBefore(context Context) {}
func (b *BaseVisitor) VisitEnumValue(context Context)        {}
func (b *BaseVisitor) VisitEnumValuesAfter(context Context)  {}
func (b *BaseVisitor) VisitEnumAfter(context Context)        {}
func (b *BaseVisitor) VisitEnumsAfter(context Context)       {}

func (b *BaseVisitor) VisitUnionsBefore(context Context)       {}
func (b *BaseVisitor) VisitUnion(context Context)              {}
func (b *BaseVisitor) VisitUnionMembersBefore(context Context) {}
func (b *BaseVisitor) VisitUnionMember(context Context)        {}
func (b *BaseVisitor) VisitUnionMembersAfter(context Context)  {}
func (b *BaseVisitor) VisitUnionsAfter(context Context)        {}

func (b *BaseVisitor) VisitAnnotationsBefore(context Context)         {}
func (b *BaseVisitor) VisitAnnotationBefore(context Context)          {}
func (b *BaseVisitor) VisitAnnotation(context Context)                {}
func (b *BaseVisitor) VisitAnnotationArgumentsBefore(context Context) {}
func (b *BaseVisitor) VisitAnnotationArgument(context Context)        {}
func (b *BaseVisitor) VisitAnnotationArgumentsAfter(context Context)  {}
func (b *BaseVisitor) VisitAnnotationAfter(context Context)           {}
func (b *BaseVisitor) VisitAnnotationsAfter(context Context)          {}

func (b *BaseVisitor) VisitDocumentAfter(context Context) {}

type MultiVisitor struct {
	_visitors []Visitor
}

func NewMultiVisitor(visitors ...Visitor) *MultiVisitor {
	return &MultiVisitor{
		_visitors: visitors,
	}
}

func (m *MultiVisitor) visit(fn func(v Visitor)) {
	for _, v := range m._visitors {
		fn(v)
	}
}

var _ = Visitor((*MultiVisitor)(nil))

func (m *MultiVisitor) VisitDocumentBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitDocumentBefore(context) })
}
func (m *MultiVisitor) VisitNamespace(context Context) {
	m.visit(func(v Visitor) { v.VisitNamespace(context) })
}

func (m *MultiVisitor) VisitImportsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitImportsBefore(context) })
}
func (m *MultiVisitor) VisitImport(context Context) {
	m.visit(func(v Visitor) { v.VisitImport(context) })
}
func (m *MultiVisitor) VisitImportsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitImportsAfter(context) })
}

func (m *MultiVisitor) VisitDirectivesBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectivesBefore(context) })
}
func (m *MultiVisitor) VisitDirectiveBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectiveBefore(context) })
}
func (m *MultiVisitor) VisitDirective(context Context) {
	m.visit(func(v Visitor) { v.VisitDirective(context) })
}
func (m *MultiVisitor) VisitDirectiveParametersBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectiveParametersBefore(context) })
}
func (m *MultiVisitor) VisitDirectiveParameter(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectiveParameter(context) })
}
func (m *MultiVisitor) VisitDirectiveParametersAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectiveParametersAfter(context) })
}
func (m *MultiVisitor) VisitDirectiveAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectiveAfter(context) })
}
func (m *MultiVisitor) VisitDirectivesAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitDirectivesAfter(context) })
}

func (m *MultiVisitor) VisitAliasesBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitAliasesBefore(context) })
}
func (m *MultiVisitor) VisitAliasBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitAliasBefore(context) })
}
func (m *MultiVisitor) VisitAlias(context Context) {
	m.visit(func(v Visitor) { v.VisitAlias(context) })
}
func (m *MultiVisitor) VisitAliasAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitAliasAfter(context) })
}
func (m *MultiVisitor) VisitAliasesAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitAliasesAfter(context) })
}

func (m *MultiVisitor) VisitAllOperationsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitAllOperationsBefore(context) })
}
func (m *MultiVisitor) VisitFunctionsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitFunctionsBefore(context) })
}
func (m *MultiVisitor) VisitFunctionBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitFunctionBefore(context) })
}
func (m *MultiVisitor) VisitFunction(context Context) {
	m.visit(func(v Visitor) { v.VisitFunction(context) })
}
func (m *MultiVisitor) VisitFunctionAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitFunctionAfter(context) })
}
func (m *MultiVisitor) VisitFunctionsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitFunctionsAfter(context) })
}
func (m *MultiVisitor) VisitInterfacesBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitInterfacesBefore(context) })
}
func (m *MultiVisitor) VisitInterfaceBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitInterfaceBefore(context) })
}
func (m *MultiVisitor) VisitInterface(context Context) {
	m.visit(func(v Visitor) { v.VisitInterface(context) })
}
func (m *MultiVisitor) VisitOperationsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitOperationsBefore(context) })
}
func (m *MultiVisitor) VisitOperationBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitOperationBefore(context) })
}
func (m *MultiVisitor) VisitOperation(context Context) {
	m.visit(func(v Visitor) { v.VisitOperation(context) })
}
func (m *MultiVisitor) VisitParametersBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitParametersBefore(context) })
}
func (m *MultiVisitor) VisitParameter(context Context) {
	m.visit(func(v Visitor) { v.VisitParameter(context) })
}
func (m *MultiVisitor) VisitParametersAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitParametersAfter(context) })
}
func (m *MultiVisitor) VisitOperationAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitOperationAfter(context) })
}
func (m *MultiVisitor) VisitOperationsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitOperationsAfter(context) })
}
func (m *MultiVisitor) VisitInterfaceAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitInterfaceAfter(context) })
}
func (m *MultiVisitor) VisitInterfacesAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitInterfacesAfter(context) })
}
func (m *MultiVisitor) VisitAllOperationsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitAllOperationsAfter(context) })
}

func (m *MultiVisitor) VisitTypesBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitTypesBefore(context) })
}
func (m *MultiVisitor) VisitTypeBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitTypeBefore(context) })
}
func (m *MultiVisitor) VisitType(context Context) {
	m.visit(func(v Visitor) { v.VisitType(context) })
}
func (m *MultiVisitor) VisitTypeFieldsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitTypeFieldsBefore(context) })
}
func (m *MultiVisitor) VisitTypeField(context Context) {
	m.visit(func(v Visitor) { v.VisitTypeField(context) })
}
func (m *MultiVisitor) VisitTypeFieldsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitTypeFieldsAfter(context) })
}
func (m *MultiVisitor) VisitTypeAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitTypeAfter(context) })
}
func (m *MultiVisitor) VisitTypesAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitTypesAfter(context) })
}

func (m *MultiVisitor) VisitEnumsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumsBefore(context) })
}
func (m *MultiVisitor) VisitEnumBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumBefore(context) })
}
func (m *MultiVisitor) VisitEnum(context Context) {
	m.visit(func(v Visitor) { v.VisitEnum(context) })
}
func (m *MultiVisitor) VisitEnumValuesBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumValuesBefore(context) })
}
func (m *MultiVisitor) VisitEnumValue(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumValue(context) })
}
func (m *MultiVisitor) VisitEnumValuesAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumValuesAfter(context) })
}
func (m *MultiVisitor) VisitEnumAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumAfter(context) })
}
func (m *MultiVisitor) VisitEnumsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitEnumsAfter(context) })
}

func (m *MultiVisitor) VisitUnionsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitUnionsBefore(context) })
}
func (m *MultiVisitor) VisitUnion(context Context) {
	m.visit(func(v Visitor) { v.VisitUnionsBefore(context) })
}
func (m *MultiVisitor) VisitUnionMembersBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitUnionMembersBefore(context) })
}
func (m *MultiVisitor) VisitUnionMember(context Context) {
	m.visit(func(v Visitor) { v.VisitUnionMember(context) })
}
func (m *MultiVisitor) VisitUnionMembersAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitUnionMembersAfter(context) })
}
func (m *MultiVisitor) VisitUnionsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitUnionsAfter(context) })
}

func (m *MultiVisitor) VisitAnnotationsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationsBefore(context) })
}
func (m *MultiVisitor) VisitAnnotationBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationBefore(context) })
}
func (m *MultiVisitor) VisitAnnotation(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotation(context) })
}
func (m *MultiVisitor) VisitAnnotationArgumentsBefore(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationArgumentsBefore(context) })
}
func (m *MultiVisitor) VisitAnnotationArgument(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationArgument(context) })
}
func (m *MultiVisitor) VisitAnnotationArgumentsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationArgumentsAfter(context) })
}
func (m *MultiVisitor) VisitAnnotationAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationAfter(context) })
}
func (m *MultiVisitor) VisitAnnotationsAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitAnnotationsAfter(context) })
}

func (m *MultiVisitor) VisitDocumentAfter(context Context) {
	m.visit(func(v Visitor) { v.VisitDocumentAfter(context) })
}
