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

package rules

import (
	"github.com/apexlang/apex-go/ast"
)

func NamespaceFirst() ast.Visitor { return &namespaceFirst{} }

type namespaceFirst struct{ ast.BaseVisitor }

func (c *namespaceFirst) VisitNamespace(context ast.Context) {
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
