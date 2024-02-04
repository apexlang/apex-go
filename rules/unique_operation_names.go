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

	"github.com/apexlang/apex-go/ast"
)

func UniqueOperationNames() ast.Visitor {
	return &uniqueOperationNames{names: map[string]struct{}{}}
}

type uniqueOperationNames struct {
	ast.BaseVisitor
	parentName string
	names      map[string]struct{}
}

func (r *uniqueOperationNames) VisitInterfaceBefore(context ast.Context) {
	r.parentName = fmt.Sprintf("interface %q", context.Interface.Name.Value)
	r.names = map[string]struct{}{}
}

func (r *uniqueOperationNames) VisitOperation(context ast.Context) {
	oper := context.Operation
	name := oper.Name.Value
	if _, duplicate := r.names[name]; duplicate {
		context.ReportError(
			ValidationError(oper.Name, "duplicate operation %q in %s", name, r.parentName),
		)
		return
	}

	r.names[name] = struct{}{}
}
