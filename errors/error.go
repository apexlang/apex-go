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

package errors

import (
	"fmt"
	"reflect"

	"github.com/apexlang/apex-go/ast"
	"github.com/apexlang/apex-go/location"
	"github.com/apexlang/apex-go/source"
)

type Error struct {
	Message       string                    `json:"message"`
	Stack         string                    `json:"stack,omitempty"`
	Nodes         []ast.Node                `json:"-"`
	Source        *source.Source            `json:"source,omitempty"`
	Positions     []uint                    `json:"positions,omitempty"`
	Locations     []location.SourceLocation `json:"locations,omitempty"`
	OriginalError error                     `json:"-"`
	Path          []interface{}             `json:"path,omitempty"`
}

type Errors []*Error

// implements Golang's built-in `error` interface
func (g Error) Error() string {
	return fmt.Sprintf("%v", g.Message)
}

func NewError(message string, nodes []ast.Node, stack string, source *source.Source, positions []uint, origError error) *Error {
	return newError(message, nodes, stack, source, positions, nil, origError)
}

func NewErrorWithPath(message string, nodes []ast.Node, stack string, source *source.Source, positions []uint, path []interface{}, origError error) *Error {
	return newError(message, nodes, stack, source, positions, path, origError)
}

func newError(message string, nodes []ast.Node, stack string, source *source.Source, positions []uint, path []interface{}, origError error) *Error {
	// if stack == "" && message != "" {
	// 	stack = message
	// }
	if source == nil {
		for _, node := range nodes {
			// get source from first node
			if node == nil || reflect.ValueOf(node).IsNil() {
				continue
			}
			if node.GetLoc() != nil {
				source = node.GetLoc().Source
			}
			break
		}
	}
	if len(positions) == 0 && len(nodes) > 0 {
		for _, node := range nodes {
			if node == nil || reflect.ValueOf(node).IsNil() {
				continue
			}
			if node.GetLoc() == nil {
				continue
			}
			positions = append(positions, node.GetLoc().Start)
		}
	}
	locations := []location.SourceLocation{}
	for _, pos := range positions {
		loc := location.GetLocation(source, pos)
		locations = append(locations, loc)
	}
	return &Error{
		Message:       message,
		Stack:         stack,
		Nodes:         nodes,
		Source:        source,
		Positions:     positions,
		Locations:     locations,
		OriginalError: origError,
		Path:          path,
	}
}
