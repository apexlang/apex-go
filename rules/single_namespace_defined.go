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

func SingleNamespaceDefined() ast.Visitor { return &singleNamespaceDefined{} }

type singleNamespaceDefined struct {
	ast.BaseVisitor
	found bool
}

func (r *singleNamespaceDefined) VisitNamespace(context ast.Context) {
	if !r.found {
		r.found = true
		return
	}

	context.ReportError(
		ValidationError(context.Namespace, "only one namespace can be defined"),
	)
}
