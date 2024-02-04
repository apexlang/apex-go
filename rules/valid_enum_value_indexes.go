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

func ValidEnumValueIndexes() ast.Visitor { return &validEnumValueIndexes{} }

type validEnumValueIndexes struct {
	ast.BaseVisitor
	parentName string
}

func (r *validEnumValueIndexes) VisitEnumBefore(context ast.Context) {
	r.parentName = context.Enum.Name.Value
}

func (r *validEnumValueIndexes) VisitEnumValue(context ast.Context) {
	enumValue := context.EnumValue
	value := enumValue.Index.Value
	if value < 0 {
		context.ReportError(
			ValidationError(
				enumValue.Index,
				"value index %d in enum %q must be a non-negative integer", value, r.parentName),
		)
	}
}
