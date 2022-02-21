# Apex Language support for Golang

TODO

```golang
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/apexlang/apex-go/parser"
)

func main() {
	schema, err := os.ReadFile("schema.apex")
	if err != nil {
		panic(err)
	}
	doc, err := parser.Parse(parser.ParseParams{
		Source: string(schema),
		Options: parser.ParseOptions{
			NoLocation: true,
			NoSource:   true,
		},
	})
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
```