module github.com/apexlang/apex-go

go 1.22

require (
	github.com/CosmWasm/tinyjson v0.9.0
	github.com/iancoleman/strcase v0.3.0
	github.com/tetratelabs/tinymem v0.1.0
	github.com/wapc/tinygo-msgpack v0.1.8
	github.com/wapc/wapc-guest-tinygo v0.3.3
)

require github.com/josharian/intern v1.0.0 // indirect

replace github.com/CosmWasm/tinyjson v0.9.0 => github.com/apexlang/tinyjson v0.9.1-0.20220929010544-92ef7a6da107
