package ast

import (
	"github.com/apexlang/apex-go/kinds"
)

// Document implements Node
type Document struct {
	BaseNode
	Definitions []Node `json:"definitions"`
}

func NewDocument(loc *Location, definitions []Node) *Document {
	return &Document{
		BaseNode:    BaseNode{kinds.Document, loc},
		Definitions: definitions,
	}
}
