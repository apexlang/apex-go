package location

import (
	"regexp"

	"github.com/apexlang/apex-go/source"
)

type SourceLocation struct {
	Line   uint `json:"line"`
	Column uint `json:"column"`
}

func GetLocation(s *source.Source, position uint) SourceLocation {
	body := []byte{}
	if s != nil {
		body = s.Body
	}
	line := uint(1)
	column := position + 1
	lineRegexp := regexp.MustCompile("\r\n|[\n\r]")
	matches := lineRegexp.FindAllIndex(body, -1)
	for _, match := range matches {
		matchIndex := uint(match[0])
		if matchIndex < position {
			line++
			l := uint(len(s.Body[match[0]:match[1]]))
			column = position + 1 - (matchIndex + l)
			continue
		} else {
			break
		}
	}
	return SourceLocation{Line: line, Column: column}
}
