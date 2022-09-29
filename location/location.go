/*
Copyright 2022 The Apex Authors.

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
