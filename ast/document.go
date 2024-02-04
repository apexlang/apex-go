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

import (
	"github.com/apexlang/apex-go/kinds"
)

// Document implements Node
type Document struct {
	BaseNode
	Definitions []Node `json:"definitions"`
}

func NewDocument(loc *Location, definitions []Node) *Document {
	return &Document{
		BaseNode:    BaseNode{kinds.Document, loc},
		Definitions: definitions,
	}
}

func (d *Document) Accept(context Context, visitor Visitor) {
	visitor.VisitDocumentBefore(context)

	if context.Namespace != nil {
		context.Namespace.Accept(context, visitor)
	}

	visitor.VisitImportsBefore(context)
	for _, importDef := range context.Imports {
		c := context
		c.Import = importDef
		importDef.Accept(c, visitor)
	}
	visitor.VisitImportsAfter(context)

	visitor.VisitDirectivesBefore(context)
	for _, directive := range context.Directives {
		c := context
		c.Directive = directive
		directive.Accept(c, visitor)
	}
	visitor.VisitDirectivesAfter(context)

	visitor.VisitAliasesBefore(context)
	for _, alias := range context.Aliases {
		c := context
		c.Alias = alias
		alias.Accept(c, visitor)
	}
	visitor.VisitAliasesAfter(context)

	visitor.VisitAllOperationsBefore(context)

	visitor.VisitFunctionsBefore(context)
	for _, function := range context.Functions {
		c := context
		c.Function = function
		function.Accept(c, visitor)
	}
	visitor.VisitFunctionsAfter(context)

	visitor.VisitInterfacesBefore(context)
	for _, iface := range context.Interfaces {
		c := context
		c.Interface = iface
		iface.Accept(c, visitor)
	}
	visitor.VisitInterfacesAfter(context)

	visitor.VisitAllOperationsAfter(context)

	visitor.VisitTypesBefore(context)
	for _, t := range context.Types {
		if t.Annotation("novisit") != nil {
			continue
		}
		c := context
		c.Type = t
		t.Accept(c, visitor)
	}
	visitor.VisitTypesAfter(context)

	visitor.VisitUnionsBefore(context)
	for _, union := range context.Unions {
		c := context
		c.Union = union
		union.Accept(c, visitor)
	}
	visitor.VisitUnionsAfter(context)

	visitor.VisitEnumsBefore(context)
	for _, enumDef := range context.Enums {
		if enumDef.Annotation("novisit") != nil {
			continue
		}
		c := context
		c.Enum = enumDef
		enumDef.Accept(c, visitor)
	}
	visitor.VisitEnumsAfter(context)

	visitor.VisitDocumentAfter(context)
}
