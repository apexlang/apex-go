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
	"strings"

	"github.com/apexlang/apex-go/ast"
)

func ValidAnnotationLocations() ast.Visitor { return &validAnnotationLocations{} }

type validAnnotationLocations struct{ ast.BaseVisitor }

func (r *validAnnotationLocations) VisitNamespace(context ast.Context) {
	r.check(context, context.Namespace.Annotations, "NAMESPACE")
}

func (r *validAnnotationLocations) VisitInterface(context ast.Context) {
	r.check(context, context.Interface.Annotations, "INTERFACE")
}

func (r *validAnnotationLocations) VisitOperation(context ast.Context) {
	r.check(context, context.Operation.Annotations, "OPERATION")
}

func (r *validAnnotationLocations) VisitParameter(context ast.Context) {
	r.check(context, context.Parameter.Annotations, "PARAMETER")
}

func (r *validAnnotationLocations) VisitType(context ast.Context) {
	r.check(context, context.Type.Annotations, "TYPE")
}

func (r *validAnnotationLocations) VisitField(context ast.Context) {
	r.check(context, context.Field.Annotations, "FIELD")
}

func (r *validAnnotationLocations) VisitEnum(context ast.Context) {
	r.check(context, context.Enum.Annotations, "ENUM")
}

func (r *validAnnotationLocations) VisitEnumValue(context ast.Context) {
	r.check(context, context.EnumValue.Annotations, "ENUM_VALUE")
}

func (r *validAnnotationLocations) VisitUnion(context ast.Context) {
	r.check(context, context.Union.Annotations, "UNION")
}

func (r *validAnnotationLocations) VisitAlias(context ast.Context) {
	r.check(context, context.Alias.Annotations, "ALIAS")
}

func (r *validAnnotationLocations) check(
	context ast.Context,
	annotations []*ast.Annotation,
	location string,
) {
	for _, annotation := range annotations {
		r.checkAnnotation(context, annotations, annotation, location)
	}
}

func (r *validAnnotationLocations) checkAnnotation(
	context ast.Context,
	annotations []*ast.Annotation,
	annotation *ast.Annotation,
	location string,
) {
	var dir *ast.DirectiveDefinition
	for _, d := range context.Directives {
		if d.Name.Value == annotation.Name.Value {
			dir = d
			break
		}
	}
	if dir == nil {
		return
	}
	found := false
	for _, loc := range dir.Locations {
		if loc.Value == location {
			found = true
			break
		}
	}
	if !found {
		context.ReportError(
			ValidationError(
				annotation,
				"annotation %q is not valid on a %q", annotation.Name.Value, strings.ToLower(strings.ReplaceAll(location, "_", " "))),
		)
		return
	}

dirRequiresLoop:
	for _, req := range dir.Requires {
		found := false
		for _, loc := range req.Locations {
			switch loc.Value {
			case "SELF":
				if findAnnotation(req.Directive.Value, annotations) {
					found = true
					break dirRequiresLoop
				}
			case "NAMESPACE":
				if findAnnotation(req.Directive.Value, context.Namespace.Annotations) {
					found = true
					break dirRequiresLoop
				}
			case "INTERFACE":
				if context.Interface != nil &&
					findAnnotation(
						req.Directive.Value,
						context.Interface.Annotations,
					) {
					found = true
					break dirRequiresLoop
				}
			case "PARAMETER":
				if context.Parameter != nil &&
					findAnnotation(
						req.Directive.Value,
						context.Parameter.Annotations,
					) {
					found = true
					break dirRequiresLoop
				}
			case "TYPE":
				if context.Type != nil &&
					findAnnotation(req.Directive.Value, context.Type.Annotations) {
					found = true
					break dirRequiresLoop
				}
			case "FIELD":
				if context.Field != nil &&
					findAnnotation(req.Directive.Value, context.Field.Annotations) {
					found = true
					break dirRequiresLoop
				}
			case "ENUM":
				if context.Enum != nil &&
					findAnnotation(req.Directive.Value, context.Enum.Annotations) {
					found = true
					break dirRequiresLoop
				}
			case "ENUM_VALUE":
				if context.EnumValue != nil &&
					findAnnotation(
						req.Directive.Value,
						context.EnumValue.Annotations) {
					found = true
					break dirRequiresLoop
				}
			case "UNION":
				if context.Union != nil &&
					findAnnotation(req.Directive.Value, context.Union.Annotations) {
					found = true
					break dirRequiresLoop
				}
			case "ALIAS":
				if context.Alias != nil &&
					findAnnotation(req.Directive.Value, context.Alias.Annotations) {
					found = true
					break dirRequiresLoop
				}
			}
		}
		if !found {
			locations := make([]string, len(req.Locations))
			for i, l := range req.Locations {
				locations[i] = strings.ToLower(l.Value)
			}
			context.ReportError(
				ValidationError(
					annotation,
					"annotation %q requires %q to exist on relative ${locations}", annotation.Name.Value, req.Directive.Value, strings.Join(locations, ", ")),
			)
		}
	}
}

func findAnnotation(name string, annotations []*ast.Annotation) bool {
	for _, annotation := range annotations {
		if annotation.Name.Value == name {
			return true
		}
	}

	return false
}
