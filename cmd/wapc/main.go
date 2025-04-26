package main

import (
	"github.com/apexlang/apex-go/model"
	"github.com/apexlang/apex-go/wapc"
)

//go:wasmexport wapc_init
func Initialize() {
	// Create providers
	resolverProvider := wapc.NewResolver()

	// Create services
	parserService := model.NewParser(resolverProvider)

	// Register services
	wapc.RegisterParser(parserService)
}
