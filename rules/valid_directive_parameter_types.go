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
	"github.com/apexlang/apex-go/kinds"
)

func ValidDirectiveParameterTypes() ast.Visitor { return &validDirectiveParameterTypes{} }

type validDirectiveParameterTypes struct{ ast.BaseVisitor }

var validTypes = map[kinds.Kind]struct{}{
	kinds.TypeDefinition: {},
	kinds.EnumDefinition: {},
}

func (r *validDirectiveParameterTypes) VisitDirectiveParameter(context ast.Context) {
	dir := context.Directive
	param := context.Parameter
	r.check(context, dir, param.Type)
}

func (r *validDirectiveParameterTypes) check(context ast.Context, dir *ast.DirectiveDefinition, t ast.Type) {
	switch v := t.(type) {
	case *ast.Optional:
		r.check(context, dir, v.Type)

	case *ast.Named:
		typeDef, ok := context.Named[v.Name.Value]
		if !ok {
			return
		}

		if _, has := validTypes[typeDef.GetKind()]; has {
			context.ReportError(
				ValidationError(
					t,
					"invalid type used in directive %q: only Types, Enums and built-in types are allowed",
					dir.Name.Value),
			)
		}

	case *ast.ListType:
		r.check(context, dir, v.Type)

	case *ast.MapType:
		r.check(context, dir, v.ValueType)
	}
}
