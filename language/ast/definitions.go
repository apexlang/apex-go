package ast

import (
	"github.com/wapc/widl-go/language/kinds"
)

type Definition interface {
	GetKind() string
	GetLoc() *Location
}

// Ensure that all definition types implements Definition interface
var _ Definition = (TypeSystemDefinition)(nil) // experimental non-spec addition.

// AnnotationDefinition implements Node, Definition
type AnnotationDefinition struct {
	Kind        string
	Loc         *Location
	Name        *Name
	Description *StringValue
	Arguments   []*InputValueDefinition
	Locations   []*Name
}

func NewAnnotationDefinition(def *AnnotationDefinition) *AnnotationDefinition {
	if def == nil {
		def = &AnnotationDefinition{}
	}
	return &AnnotationDefinition{
		Kind:        kinds.AnnotationDefinition,
		Loc:         def.Loc,
		Name:        def.Name,
		Description: def.Description,
		Arguments:   def.Arguments,
		Locations:   def.Locations,
	}
}

func (def *AnnotationDefinition) GetKind() string {
	return def.Kind
}

func (def *AnnotationDefinition) GetLoc() *Location {
	return def.Loc
}

func (def *AnnotationDefinition) GetDescription() *StringValue {
	return def.Description
}
