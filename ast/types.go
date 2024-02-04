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
var _ Type = (*Stream)(nil)

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

// Stream implements Node, Type
type Stream struct {
	BaseNode
	Type Type `json:"type"`
}

func NewStream(loc *Location, t Type) *Stream {
	return &Stream{
		BaseNode: BaseNode{kinds.Stream, loc},
		Type:     t,
	}
}
