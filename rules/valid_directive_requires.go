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

func ValidDirectiveRequires() ast.Visitor { return &validDirectiveRequires{} }

type validDirectiveRequires struct{ ast.BaseVisitor }

func (r *validDirectiveRequires) VisitDirective(context ast.Context) {
	dir := context.Directive
	dirName := dir.Name.Value

	for _, req := range dir.Requires {
		found := false
		for _, d := range context.Directives {
			if d.Name.Value == req.Directive.Value {
				found = true
				break
			}
		}
		if !found {
			context.ReportError(
				ValidationError(
					req.Directive,
					"unknown required directive %q on %q", req.Directive.Value, dirName),
			)
		}
	}
}
