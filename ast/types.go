package ast

import (
	"github.com/apexlang/apex-go/kinds"
)

type Type interface {
	Node
}

// Ensure that all value types implements Value interface
var _ Type = (*Named)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*MapType)(nil)
var _ Type = (*Optional)(nil)

// Named implements Node, Type
type Named struct {
	BaseNode
	Name *Name `json:"name"`
}

func NewNamed(loc *Location, name *Name) *Named {
	return &Named{
		BaseNode: BaseNode{kinds.Named, loc},
		Name:     name,
	}
}

// ListType implements Node, Type
type ListType struct {
	BaseNode
	Type Type `json:"type"`
}

func NewListType(loc *Location, t Type) *ListType {
	return &ListType{
		BaseNode: BaseNode{kinds.ListType, loc},
		Type:     t,
	}
}

// MapType implements Node, Type
type MapType struct {
	BaseNode
	KeyType   Type `json:"keyType"`
	ValueType Type `json:"valueType"`
}

func NewMapType(loc *Location, keyType, valueType Type) *MapType {
	return &MapType{
		BaseNode:  BaseNode{kinds.MapType, loc},
		KeyType:   keyType,
		ValueType: valueType,
	}
}

// Optional implements Node, Type
type Optional struct {
	BaseNode
	Type Type `json:"type"`
}

func NewOptional(loc *Location, t Type) *Optional {
	return &Optional{
		BaseNode: BaseNode{kinds.Optional, loc},
		Type:     t,
	}
}
