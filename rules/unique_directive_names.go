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

func UniqueDirectiveNames() ast.Visitor { return &uniqueDirectiveNames{names: map[string]struct{}{}} }

type uniqueDirectiveNames struct {
	ast.BaseVisitor
	names map[string]struct{}
}

func (r *uniqueDirectiveNames) VisitDirective(context ast.Context) {
	directive := context.Directive
	name := directive.Name.Value
	if _, duplicate := r.names[name]; duplicate {
		context.ReportError(
			ValidationError(directive.Name, "duplicate directive %q", name),
		)
		return
	}

	r.names[name] = struct{}{}
}
