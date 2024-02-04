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
)

func UniqueObjectNames() ast.Visitor { return &uniqueObjectNames{names: map[string]struct{}{}} }

type uniqueObjectNames struct {
	ast.BaseVisitor
	names map[string]struct{}
}

func (r *uniqueObjectNames) VisitInterface(context ast.Context) {
	r.check(context, context.Interface.Name, "interface")
}

func (r *uniqueObjectNames) VisitType(context ast.Context) {
	r.check(context, context.Type.Name, "type")
}

func (r *uniqueObjectNames) VisitUnion(context ast.Context) {
	r.check(context, context.Union.Name, "union")
}

func (r *uniqueObjectNames) VisitEnum(context ast.Context) {
	r.check(context, context.Enum.Name, "enum")
}

func (r *uniqueObjectNames) VisitAlias(context ast.Context) {
	r.check(context, context.Alias.Name, "alias")
}

func (r *uniqueObjectNames) check(context ast.Context, name *ast.Name, typeName string) {
	if _, duplicate := r.names[name.Value]; duplicate {
		context.ReportError(
			ValidationError(name, "duplicate %s %q", typeName, name),
		)
		return
	}

	r.names[name.Value] = struct{}{}
}
