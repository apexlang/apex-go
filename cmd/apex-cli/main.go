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

package main

import (
	"io"
	"os"

	"github.com/apexlang/apex-go/errors"
	"github.com/apexlang/apex-go/model"
	"github.com/apexlang/apex-go/parser"
	"github.com/apexlang/apex-go/rules"
)

func main() {
	specBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		errors.Write(err)
		return
	}

	doc, err := parser.Parse(parser.ParseParams{
		Source: string(specBytes),
		Options: parser.ParseOptions{
			NoSource: true,
		},
	})
	if err != nil {
		errors.Write(err)
		return
	}

	errs := rules.Validate(doc, rules.Rules...)
	if len(errs) > 0 {
		errors.Write(errs...)
		return
	}

	ns, errs := model.Convert(doc)
	if len(errs) > 0 {
		errors.Write(errs...)
		return
	}

	jsonBytes, err := ns.MarshalJSON()
	if err != nil {
		errors.Write(err)
		return
	}

	os.Stdout.Write(jsonBytes)
}
