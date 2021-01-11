package ast

import (
	"github.com/wapc/widl-go/language/kinds"
)

// Annotation implements Node
type Annotation struct {
	Kind      string
	Loc       *Location
	Name      *Name
	Arguments []*Argument
}

func NewAnnotation(an *Annotation) *Annotation {
	if an == nil {
		an = &Annotation{}
	}
	return &Annotation{
		Kind:      kinds.Annotation,
		Loc:       an.Loc,
		Name:      an.Name,
		Arguments: an.Arguments,
	}
}

func (an *Annotation) GetKind() string {
	return an.Kind
}

func (an *Annotation) GetLoc() *Location {
	return an.Loc
}
