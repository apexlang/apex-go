module github.com/apexlang/apex-go

go 1.18

require (
	github.com/CosmWasm/tinyjson v0.9.0
	github.com/iancoleman/strcase v0.2.0
	github.com/tetratelabs/tinymem v0.1.0
)

require github.com/josharian/intern v1.0.0 // indirect

replace github.com/CosmWasm/tinyjson v0.9.0 => github.com/apexlang/tinyjson v0.9.1-0.20220929010544-92ef7a6da107
