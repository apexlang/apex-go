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

package main

import (
	"fmt"

	"github.com/tetratelabs/tinymem"

	"github.com/apexlang/apex-go/errors"
	"github.com/apexlang/apex-go/model"
	"github.com/apexlang/apex-go/parser"
	"github.com/apexlang/apex-go/rules"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//go:wasm-module apex
//go:export resolve
func resolve(
	locationPtr uintptr, locationSize uint32,
	fromPtr uintptr, fromSize uint32) uint64

//export parse
func Parse(ptr uintptr, size uint32) (ptrSize uint64) {
	source := tinymem.PtrToString(ptr, size)

	doc, err := parser.Parse(parser.ParseParams{
		Source: source,
		Options: parser.ParseOptions{
			NoSource: true,
			Resolver: func(location, from string) (string, error) {
				locationPtr, locationSize := tinymem.StringToPtr(location)
				fromPtr, fromSize := tinymem.StringToPtr(from)
				ptrsize := resolve(
					locationPtr, locationSize,
					fromPtr, fromSize)
				if ptrsize == 0 {
					return "", fmt.Errorf("could not find %q", location)
				}

				ptr := uintptr(ptrsize >> 32)
				size := uint32(ptrsize & 0xFFFFFFFF)
				return tinymem.PtrToString(ptr, size), nil
			},
		},
	})
	if err != nil {
		return errors.Return(err)
	}

	errs := rules.Validate(doc, rules.Rules...)
	if len(errs) > 0 {
		return errors.Return(errs...)
	}

	ns, errs := model.Convert(doc)
	if len(errs) > 0 {
		return errors.Return(errs...)
	}

	jsonBytes, err := ns.MarshalJSON()
	if err != nil {
		return errors.Return(err)
	}

	jsonString := string(jsonBytes)
	ptr, size = tinymem.StringToPtr(jsonString)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}
