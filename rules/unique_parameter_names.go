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
	"fmt"

	"github.com/apexlang/apex-go/ast"
)

func UniqueParameterNames() ast.Visitor {
	return &uniqueParameterNames{names: map[string]struct{}{}}
}

type uniqueParameterNames struct {
	ast.BaseVisitor
	parentName string
	names      map[string]struct{}
}

func (r *uniqueParameterNames) VisitFunctionBefore(context ast.Context) {
	function := context.Function
	r.parentName = fmt.Sprintf("func %q", function.Name.Value)
	r.names = map[string]struct{}{}
}

func (r *uniqueParameterNames) VisitOperationBefore(context ast.Context) {
	iface := context.Interface
	oper := context.Operation
	r.parentName = fmt.Sprintf("operation \"%s::%s\"", iface.Name.Value, oper.Name.Value)
	r.names = map[string]struct{}{}
}

func (r *uniqueParameterNames) VisitParameter(context ast.Context) {
	oper := context.Parameter
	name := oper.Name.Value
	if _, duplicate := r.names[name]; duplicate {
		context.ReportError(
			ValidationError(oper.Name, "duplicate parameter %q in %s", name, r.parentName),
		)
		return
	}

	r.names[name] = struct{}{}
}
