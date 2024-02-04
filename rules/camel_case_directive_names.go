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
	"github.com/iancoleman/strcase"

	"github.com/apexlang/apex-go/ast"
)

func CamelCaseDirectiveNames() ast.Visitor { return &camelCaseDirectiveNames{} }

type camelCaseDirectiveNames struct{ ast.BaseVisitor }

func (c *camelCaseDirectiveNames) VisitDirective(context ast.Context) {
	directive := context.Directive
	name := directive.Name.Value
	if name != strcase.ToLowerCamel(name) {
		context.ReportError(
			ValidationError(
				directive.Name,
				"directive %s should be camel case",
				directive.Name.Value,
			),
		)
	}
}
