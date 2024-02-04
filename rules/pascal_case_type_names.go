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

package rules

import (
	"github.com/apexlang/apex-go/ast"
	"github.com/iancoleman/strcase"
)

func PascalCaseTypeNames() ast.Visitor { return &pascelCaseTypeNames{} }

type pascelCaseTypeNames struct{ ast.BaseVisitor }

func (r *pascelCaseTypeNames) VisitNamespace(context ast.Context) {
	pos := 0
	for _, def := range context.Document.Definitions {
		switch v := def.(type) {
		case *ast.ImportDefinition, *ast.DirectiveDefinition:
			// Ignore the position
		case *ast.NamespaceDefinition:
			if pos == 0 {
				return
			}
			context.ReportError(
				ValidationError(
					v,
					"namespace must be defined before any other definition",
				),
			)
		default:
			pos++
		}
	}
}

func (r *pascelCaseTypeNames) VisitAlias(context ast.Context) {
	alias := context.Alias
	name := alias.Name.Value
	if name != strcase.ToCamel(name) {
		context.ReportError(
			ValidationError(alias.Name, "alias %q should be pascal case", name),
		)
	}
}

func (r *pascelCaseTypeNames) VisitType(context ast.Context) {
	t := context.Type
	name := t.Name.Value
	if name != strcase.ToCamel(name) {
		context.ReportError(
			ValidationError(t.Name, "type %q should be pascal case", name),
		)
	}
}

func (r *pascelCaseTypeNames) VisitEnum(context ast.Context) {
	enumDef := context.Enum
	name := enumDef.Name.Value
	if name != strcase.ToCamel(name) {
		context.ReportError(
			ValidationError(enumDef.Name, "enum %q should be pascal case", name),
		)
	}
}

func (r *pascelCaseTypeNames) VisitUnion(context ast.Context) {
	union := context.Union
	name := union.Name.Value
	if name != strcase.ToCamel(name) {
		context.ReportError(
			ValidationError(union.Name, "union %q should be pascal case", name),
		)
	}
}
