package ast

import (
	"github.com/wapc/widl-go/language/kinds"
)

type Type interface {
	GetKind() string
	GetLoc() *Location
	String() string
}

// Ensure that all value types implements Value interface
var _ Type = (*Named)(nil)
var _ Type = (*List)(nil)
var _ Type = (*Map)(nil)
var _ Type = (*Optional)(nil)

// Named implements Node, Type
type Named struct {
	Kind string
	Loc  *Location
	Name *Name
}

func NewNamed(t *Named) *Named {
	if t == nil {
		t = &Named{}
	}
	t.Kind = kinds.Named
	return t
}

func (t *Named) GetKind() string {
	return t.Kind
}

func (t *Named) GetLoc() *Location {
	return t.Loc
}

func (t *Named) String() string {
	return t.GetKind()
}

// List implements Node, Type
type List struct {
	Kind string
	Loc  *Location
	Type Type
}

func NewList(t *List) *List {
	if t == nil {
		t = &List{}
	}
	return &List{
		Kind: kinds.List,
		Loc:  t.Loc,
		Type: t.Type,
	}
}

func (t *List) GetKind() string {
	return t.Kind
}

func (t *List) GetLoc() *Location {
	return t.Loc
}

func (t *List) String() string {
	return t.GetKind()
}

// Map implements Node, Type
type Map struct {
	Kind      string
	Loc       *Location
	KeyType   Type
	ValueType Type
}

func NewMap(t *Map) *Map {
	if t == nil {
		t = &Map{}
	}
	return &Map{
		Kind:      kinds.List,
		Loc:       t.Loc,
		KeyType:   t.KeyType,
		ValueType: t.ValueType,
	}
}

func (t *Map) GetKind() string {
	return t.Kind
}

func (t *Map) GetLoc() *Location {
	return t.Loc
}

func (t *Map) String() string {
	return t.GetKind()
}

// Optional implements Node, Type
type Optional struct {
	Kind string
	Loc  *Location
	Type Type
}

func NewOptional(t *Optional) *Optional {
	if t == nil {
		t = &Optional{}
	}
	return &Optional{
		Kind: kinds.Optional,
		Loc:  t.Loc,
		Type: t.Type,
	}
}

func (t *Optional) GetKind() string {
	return t.Kind
}

func (t *Optional) GetLoc() *Location {
	return t.Loc
}

func (t *Optional) String() string {
	return t.GetKind()
}
