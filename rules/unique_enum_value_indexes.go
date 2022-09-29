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

func UniqueEnumValueIndexes() ast.Visitor {
	return &uniqueEnumValueIndexes{values: map[int]struct{}{}}
}

type uniqueEnumValueIndexes struct {
	ast.BaseVisitor
	parentName string
	values     map[int]struct{}
}

func (r *uniqueEnumValueIndexes) VisitEnumBefore(context ast.Context) {
	r.parentName = context.Enum.Name.Value
	r.values = map[int]struct{}{}
}

func (r *uniqueEnumValueIndexes) VisitEnumValue(context ast.Context) {
	enumValue := context.EnumValue
	value := enumValue.Index.Value
	if _, duplicate := r.values[value]; duplicate {
		context.ReportError(
			ValidationError(enumValue.Index, "duplicate index %d in enum %q", value, r.parentName),
		)
		return
	}

	r.values[value] = struct{}{}
}
