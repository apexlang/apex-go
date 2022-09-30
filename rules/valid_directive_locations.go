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

func ValidDirectiveLocation() ast.Visitor { return &validDirectiveLocation{} }

type validDirectiveLocation struct{ ast.BaseVisitor }

var validLocationNames = map[string]struct{}{
	"NAMESPACE":  {},
	"INTERFACE":  {},
	"OPERATION":  {},
	"PARAMETER":  {},
	"TYPE":       {},
	"FIELD":      {},
	"ENUM":       {},
	"ENUM_VALUE": {},
	"UNION":      {},
	"ALIAS":      {},
}

func (r *validDirectiveLocation) VisitDirective(context ast.Context) {
	dir := context.Directive
	dirName := dir.Name.Value
	locationNames := make(map[string]struct{})

	for _, loc := range dir.Locations {
		if _, valid := validLocationNames[loc.Value]; !valid {
			context.ReportError(
				ValidationError(
					loc,
					"invalid directive location %q on %q", loc.Value, dirName),
			)
		}
		if _, duplicate := locationNames[loc.Value]; duplicate {
			context.ReportError(
				ValidationError(
					loc,
					"duplicate directive location %q on %q", loc.Value, dirName),
			)
		}

		locationNames[loc.Value] = struct{}{}
	}

	for _, req := range dir.Requires {
		requireLocationNames := make(map[string]struct{})

		for _, loc := range req.Locations {
			_, valid := validLocationNames[loc.Value]
			if loc.Value != "SELF" && !valid {
				context.ReportError(
					ValidationError(
						loc,
						"invalid directive location %q on %q", loc.Value, dirName),
				)
			}
			if _, duplicate := requireLocationNames[loc.Value]; duplicate {
				context.ReportError(
					ValidationError(
						loc,
						"duplicate directive location %q on %q", loc.Value, dirName),
				)
			}

			requireLocationNames[loc.Value] = struct{}{}
		}
	}
}
