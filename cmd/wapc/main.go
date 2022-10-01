package main

import (
	"github.com/apexlang/apex-go/model"
)

func main() {
	// Create providers
	resolverProvider := model.NewResolver()

	// Create services
	parserService := model.NewParser(resolverProvider)

	// Register services
	model.RegisterParser(parserService)
}
