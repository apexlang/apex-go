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
	"fmt"
	"strings"

	"github.com/apexlang/apex-go/ast"
)

func KnownTypes() ast.Visitor { return &knownTypes{} }

type knownTypes struct{ ast.BaseVisitor }

var builtInTypeNames = map[string]struct{}{
	"i8":       {},
	"u8":       {},
	"i16":      {},
	"u16":      {},
	"i32":      {},
	"u32":      {},
	"i64":      {},
	"u64":      {},
	"f32":      {},
	"f64":      {},
	"bool":     {},
	"string":   {},
	"datetime": {},
	"bytes":    {},
	"any":      {},
	"raw":      {},
	"value":    {},
}

func (c *knownTypes) VisitAlias(context ast.Context) {
	alias := context.Alias
	c.checkType(context, `alias`, alias.Name.Value, alias.Type)
}

func (c *knownTypes) VisitOperationAfter(context ast.Context) {
	oper := context.Operation
	// "void" is a special case for operations without a return.
	if named, ok := oper.Type.(*ast.Named); ok && named.Name.Value == "void" {
		return
	}

	c.checkType(
		context,
		`return`,
		oper.Name.Value,
		oper.Type,
	)
}

func (c *knownTypes) VisitParameter(context ast.Context) {
	oper := context.Operation
	param := context.Parameter
	c.checkType(
		context,
		fmt.Sprintf(`parameter %q`, param.Name.Value),
		oper.Name.Value,
		param.Type,
	)
}

func (c *knownTypes) VisitField(context ast.Context) {
	t := context.Type
	field := context.Field
	c.checkType(
		context,
		fmt.Sprintf(`field %q`, field.Name.Value),
		t.Name.Value,
		field.Type,
	)
}

func (c *knownTypes) VisitUnion(context ast.Context) {
	union := context.Union
	for _, ut := range union.Members {
		c.checkType(
			context,
			fmt.Sprintf(`union %q`, union.Name.Value),
			union.Name.Value,
			ut.Type,
		)
	}
}

func (c *knownTypes) VisitDirectiveParameter(context ast.Context) {
	directive := context.Directive
	param := context.Parameter
	c.checkType(
		context,
		fmt.Sprintf(`parameter %q`, param.Name.Value),
		directive.Name.Value,
		param.Type,
	)
}

func (c *knownTypes) checkType(
	context ast.Context,
	forName string,
	parentName string,
	t ast.Type,
) {
	switch v := t.(type) {
	case *ast.Named:
		name := v.Name.Value
		first := name[0:1]

		if first == strings.ToLower(first) {
			// Check for built-in types
			if _, ok := builtInTypeNames[name]; !ok {
				context.ReportError(
					ValidationError(
						v,
						"invalid built-in type %q for %s in %q",
						name,
						forName,
						parentName,
					),
				)
			}
		} else {
			// Check against defined types
			if _, ok := context.Named[name]; !ok {
				context.ReportError(
					ValidationError(
						v,
						"unknown type %q for %s in %q",
						name,
						forName,
						parentName,
					),
				)
			}
		}

	case *ast.Optional:
		c.checkType(context, forName, parentName, v.Type)

	case *ast.MapType:
		c.checkType(context, forName, parentName, v.KeyType)
		c.checkType(context, forName, parentName, v.ValueType)

	case *ast.ListType:
		c.checkType(context, forName, parentName, v.Type)
	}
}
