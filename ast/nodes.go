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
