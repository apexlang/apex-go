package ast

import (
	"github.com/apexlang/apex-go/source"
)

type Location struct {
	Start  uint
	End    uint
	Source *source.Source `json:"source,omitempty"`
}

func NewLocation(start, end uint, source *source.Source) *Location {
	return &Location{
		Start:  start,
		End:    end,
		Source: source,
	}
}
