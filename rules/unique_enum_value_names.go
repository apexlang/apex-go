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

func UniqueEnumValueNames() ast.Visitor {
	return &uniqueEnumValueNames{names: map[string]struct{}{}}
}

type uniqueEnumValueNames struct {
	ast.BaseVisitor
	parentName string
	names      map[string]struct{}
}

func (r *uniqueEnumValueNames) VisitEnumBefore(context ast.Context) {
	r.parentName = context.Enum.Name.Value
	r.names = map[string]struct{}{}
}

func (r *uniqueEnumValueNames) VisitEnumValue(context ast.Context) {
	enumValue := context.EnumValue
	name := enumValue.Name.Value
	if _, duplicate := r.names[name]; duplicate {
		context.ReportError(
			ValidationError(enumValue.Index, "duplicate value %q in enum %q", name, r.parentName),
		)
		return
	}

	r.names[name] = struct{}{}
}
