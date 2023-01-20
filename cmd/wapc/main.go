package main

import (
	"github.com/apexlang/apex-go/model"
	"github.com/apexlang/apex-go/wapc"
)

func main() {
	// Create providers
	resolverProvider := wapc.NewResolver()

	// Create services
	parserService := model.NewParser(resolverProvider)

	// Register services
	wapc.RegisterParser(parserService)
}
