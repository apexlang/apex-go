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
	"strings"

	"github.com/apexlang/apex-go/ast"
	"github.com/apexlang/apex-go/kinds"
)

func ValidAnnotationArguments() ast.Visitor { return &validAnnotationArguments{} }

type validAnnotationArguments struct{ ast.BaseVisitor }

var integerBuiltInTypeNames = map[string]struct{}{
	"i8":  {},
	"u8":  {},
	"i16": {},
	"u16": {},
	"i32": {},
	"u32": {},
	"i64": {},
	"u64": {},
}

var floatBuiltInTypeNames = map[string]struct{}{
	"f32": {},
	"f64": {},
}

func (r *validAnnotationArguments) VisitAnnotation(context ast.Context) {
	a := context.Annotation

	foundArgNames := make(map[string]struct{})
	for _, v := range a.Arguments {
		if _, duplicate := foundArgNames[v.Name.Value]; duplicate {
			context.ReportError(
				ValidationError(
					v,
					"duplicate argument %q in annotation %q", v.Name.Value, a.Name.Value),
			)
		}

		foundArgNames[v.Name.Value] = struct{}{}
	}

	var dir *ast.DirectiveDefinition
	for _, d := range context.Directives {
		if d.Name.Value == a.Name.Value {
			dir = d
			break
		}
	}
	if dir == nil {
		return
	}

	args := make(map[string]*ast.Argument)
	for _, arg := range a.Arguments {
		args[arg.Name.Value] = arg
	}

	for _, param := range dir.Parameters {
		arg := args[param.Name.Value]
		if arg == nil {
			if !param.Type.IsKind(kinds.Optional) {
				context.ReportError(
					ValidationError(
						a,
						"missing required argument %q in annotation %q", param.Name.Value, a.Name.Value),
				)
			}
			continue
		}
		delete(args, param.Name.Value)

		// Validate types
		r.check(context, param.Type, arg.Value, a)
	}

	for _, arg := range args {
		context.ReportError(
			ValidationError(
				arg,
				"unknown parameter %q in directive %q", arg.Name.Value, dir.Name.Value),
		)
	}
}

func (r *validAnnotationArguments) check(
	context ast.Context,
	t ast.Type,
	value ast.Value,
	annotation *ast.Annotation) {
	switch v := t.(type) {
	case *ast.Optional:
		r.check(context, v.Type, value, annotation)

	case *ast.Named:
		_, isInteger := integerBuiltInTypeNames[v.Name.Value]
		_, isFloat := floatBuiltInTypeNames[v.Name.Value]
		if v.Name.Value == "string" {
			if !value.IsKind(kinds.StringValue) {
				context.ReportError(
					ValidationError(
						value,
						"invalid value %q in annotation %q: expected a string", value.GetValue(), annotation.Name.Value),
				)
			}
		} else if isInteger {
			intValue, ok := value.(*ast.IntValue)
			if !ok {
				context.ReportError(
					ValidationError(
						value,
						`invalid value %q in annotation %q: expected an integer`, value.GetValue(), annotation.Name.Value),
				)
				return
			}
			if strings.HasPrefix(v.Name.Value, "u") && intValue.Value < 0 {
				context.ReportError(
					ValidationError(
						value,
						`invalid value %q in annotation %q: expected a non-negative integer`, value.GetValue(), annotation.Name.Value),
				)
			}
		} else if isFloat {
			_, ok := value.(*ast.FloatValue)
			if !ok {
				context.ReportError(
					ValidationError(
						value,
						`invalid value %q in annotation %q: expected a float`, value.GetValue(), annotation.Name.Value),
				)
			}
		} else if v.Name.Value == "bool" {
			_, ok := value.(*ast.BooleanValue)
			if !ok {
				context.ReportError(
					ValidationError(
						value,
						`invalid value %q in annotation %q: expected a boolean`, value.GetValue(), annotation.Name.Value),
				)
			}
		} else {
			definition := context.Named[v.Name.Value]
			if definition == nil {
				// error reported by KnownTypes
				return
			}
			switch defv := definition.(type) {
			case *ast.EnumDefinition:
				expectedEnumValue, ok := value.(*ast.EnumValue)
				if !ok {
					context.ReportError(
						ValidationError(
							value,
							`invalid value %q in annotation %q: expected an enum value`, value.GetValue(), annotation.Name.Value),
					)
					return
				}

				enumDef := defv
				found := false
				for _, v := range enumDef.Values {
					if expectedEnumValue.Value == v.Name.Value {
						found = true
						break
					}
				}
				if !found {
					context.ReportError(
						ValidationError(
							value,
							"unknown enum value %q in annotation %q: expected value from %q", expectedEnumValue.Value, annotation.Name.Value, enumDef.Name.Value),
					)
				}
			case *ast.TypeDefinition:
				{
					obj, ok := value.(*ast.ObjectValue)
					if !ok {
						context.ReportError(
							ValidationError(
								value,
								"invalid value %q in annotation %q: expected an object", value.GetValue(), annotation.Name.Value),
						)
						return
					}

					fields := make(map[string]*ast.FieldDefinition)
					for _, f := range defv.Fields {
						fields[f.Name.Value] = f
					}

					for _, field := range obj.Fields {
						f, ok := fields[field.Name.Value]
						if !ok {
							context.ReportError(
								ValidationError(
									field.Name,
									"unknown field %q for type %q in annotation %q", field.Name.Value, defv.Name.Value, annotation.Name.Value),
							)
							continue
						}
						delete(fields, field.Name.Value)

						// Validate types
						r.check(context, f.Type, field.Value, annotation)
					}

					for _, field := range fields {
						if !field.Type.IsKind(kinds.Optional) {
							context.ReportError(
								ValidationError(
									obj,
									"missing required field %q for type %q in annotation %q", field.Name.Value, defv.Name.Value, annotation.Name.Value),
							)
						}
					}
				}
			default:
				context.ReportError(
					ValidationError(
						value,
						"invalid value %q in annotation %q: expected an object", value.GetValue(), annotation.Name.Value),
				)
			}
		}
	case *ast.ListType:
		list := v
		listValue, ok := value.(*ast.ListValue)
		if !ok {
			context.ReportError(
				ValidationError(
					value,
					"invalid value %q in annotation %q: expected a list", value.GetValue(), annotation.Name.Value),
			)
			return
		}
		for _, value := range listValue.Values {
			r.check(context, list.Type, value, annotation)
		}

	case *ast.MapType:
		objectValue, ok := value.(*ast.ObjectValue)
		if !ok {
			context.ReportError(
				ValidationError(
					value,
					"invalid value %q in annotation %q: expected a map", value.GetValue(), annotation.Name.Value),
			)
			return
		}
		for _, field := range objectValue.Fields {
			r.check(context, v.ValueType, field.Value, annotation)
		}
	}
}
