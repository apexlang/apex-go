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

import "github.com/apexlang/apex-go/kinds"

type Node interface {
	GetKind() kinds.Kind
	IsKind(kind kinds.Kind) bool
	GetLoc() *Location
	String() string
}

type BaseNode struct {
	Kind kinds.Kind `json:"kind"`
	Loc  *Location  `json:"-"`
}

func (node *BaseNode) GetKind() kinds.Kind {
	return node.Kind
}

func (node *BaseNode) IsKind(kind kinds.Kind) bool {
	return node.Kind == kind
}

func (node *BaseNode) GetLoc() *Location {
	return node.Loc
}

func (node *BaseNode) String() string {
	return string(node.Kind)
}

// Name implements Node
var _ Node = (*Name)(nil)

type Name struct {
	BaseNode
	Value string `json:"value"`
}

func NewName(loc *Location, value string) *Name {
	return &Name{
		BaseNode: BaseNode{kinds.Name, loc},
		Value:    value,
	}
}

// Annotation implements Node
var _ Node = (*Annotation)(nil)

type Annotation struct {
	BaseNode
	Name      *Name       `json:"name"`
	Arguments []*Argument `json:"arguments"`
}

func NewAnnotation(loc *Location, name *Name, arguments []*Argument) *Annotation {
	return &Annotation{
		BaseNode:  BaseNode{kinds.Annotation, loc},
		Name:      name,
		Arguments: arguments,
	}
}

func (a *Annotation) Accept(context Context, visitor Visitor) {
	visitor.VisitAnnotation(context)
}

// Argument implements Node
var _ Node = (*Argument)(nil)

type Argument struct {
	BaseNode
	Name  *Name `json:"name"`
	Value Value `json:"value"`
}

func NewArgument(loc *Location, name *Name, value Value) *Argument {
	return &Argument{
		BaseNode: BaseNode{kinds.Argument, loc},
		Name:     name,
		Value:    value,
	}
}

// DirectiveRequire implements Node
var _ Node = (*DirectiveRequire)(nil)

type DirectiveRequire struct {
	BaseNode
	Directive *Name   `json:"directive"`
	Locations []*Name `json:"locations"`
}

func NewDirectiveRequire(loc *Location, directive *Name, locations []*Name) *DirectiveRequire {
	return &DirectiveRequire{
		BaseNode:  BaseNode{kinds.Argument, loc},
		Directive: directive,
		Locations: locations,
	}
}

func (n *DirectiveRequire) HasLocation(location string) bool {
	for _, l := range n.Locations {
		if l.Value == location {
			return true
		}
	}
	return false
}

// ImportName implements Node
var _ Node = (*ImportName)(nil)

type ImportName struct {
	BaseNode
	Name  *Name `json:"name"`
	Alias *Name `json:"alias,omitempty"` // Optional
}

func NewImportName(loc *Location, name *Name, alias *Name) *ImportName {
	return &ImportName{
		BaseNode: BaseNode{kinds.ImportName, loc},
		Name:     name,
		Alias:    alias,
	}
}
