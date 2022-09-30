.PHONY: all wasm-cli wasm-api codegen

all: codegen wasm-cli wasm-api wasm-host

wasm-cli:
	tinygo build -o apex-cli.wasm -scheduler=none -target=wasi -wasm-abi=generic -no-debug cmd/apex-cli/main.go
	wasm-opt -O apex-cli.wasm -o apex-cli.wasm

wasm-api:
	tinygo build -o apex-api.wasm -scheduler=none -target=wasi -wasm-abi=generic -no-debug cmd/apex-api/main.go
	wasm-opt -O apex-api.wasm -o apex-api.wasm
	cp apex-api.wasm cmd/host/apex-api.wasm

wasm-host:
	go build -o apex-host cmd/host/main.go

codegen:
	apex generate
