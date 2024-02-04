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
	"github.com/apexlang/apex-go/errors"
	"github.com/apexlang/apex-go/source"
)

type ValidationRule func() ast.Visitor

var Rules = []ValidationRule{
	CamelCaseDirectiveNames,
	KnownTypes,
	NamespaceFirst,
	PascalCaseTypeNames,
	SingleNamespaceDefined,
	UniqueEnumValueIndexes,
	UniqueEnumValueNames,
	UniqueFunctionNames,
	UniqueObjectNames,
	UniqueOperationNames,
	UniqueParameterNames,
	UniqueTypeFieldNames,
	ValidAnnotationArguments,
	ValidAnnotationLocations,
	ValidDirectiveLocation,
	ValidDirectiveParameterTypes,
	ValidDirectiveRequires,
	ValidEnumValueIndexes,
}

func Validate(
	doc *ast.Document,
	rules ...ValidationRule,
) []error {
	context := ast.NewContext(doc)

	ruleVisitors := make([]ast.Visitor, len(rules))
	for i, rule := range rules {
		ruleVisitors[i] = rule()
	}

	visitor := ast.NewMultiVisitor(ruleVisitors...)

	doc.Accept(context, visitor)
	return context.Errors()
}

func ValidationError(node ast.Node, format string, a ...interface{}) *errors.Error {
	loc := node.GetLoc()
	var source *source.Source
	if loc != nil {
		source = loc.Source
	}

	return errors.NewError(
		fmt.Sprintf(`Validation Error: `+format, a...),
		[]ast.Node{node},
		"",
		source,
		nil,
		nil,
	)
}
