package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	// "github.com/apexlang/apex-go/parser"
)

//go:embed apex-api.wasm
var apexWasm []byte

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: apex-host <apex file>")
		return
	}
	specFile := os.Args[1]
	ctx := context.Background()
	specBytes, err := os.ReadFile(specFile)
	if err != nil {
		panic(err)
	}

	config := wazero.NewModuleConfig().
		WithStdout(os.Stdout).WithStderr(os.Stderr)

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithCoreFeatures(api.CoreFeaturesV2))
	defer r.Close(ctx)

	homeDir, err := getHomeDirectory()
	if err != nil {
		panic(err)
	}
	definitionsDir := filepath.Join(homeDir, "definitions")

	var malloc, free api.Function

	returnString := func(m api.Module, value string) uint64 {
		size := uint64(len(value))
		results, err := malloc.Call(ctx, size)
		if err != nil {
			panic(err)
		}

		ptr := uintptr(results[0])

		m.Memory().Write(ctx, uint32(ptr), []byte(value))
		return (uint64(ptr) << uint64(32)) | uint64(size)
	}

	resolve := func(ctx context.Context, m api.Module, locationPtr, locationLen, fromPtr, fromLen uint32) uint64 {
		locationBuf, ok := m.Memory().Read(ctx, locationPtr, locationLen)
		if !ok {
			return returnString(m, fmt.Sprintf("error: %v", err))
		}
		location := string(locationBuf)

		loc := filepath.Join(definitionsDir, filepath.Join(strings.Split(location, "/")...))
		if filepath.Ext(loc) != ".apex" {
			specLoc := loc + ".apex"
			found := false
			stat, err := os.Stat(specLoc)
			if err == nil && !stat.IsDir() {
				found = true
				loc = specLoc
			}

			if !found {
				stat, err := os.Stat(loc)
				if err != nil {
					return returnString(m, fmt.Sprintf("error: %v", err))
				}
				if stat.IsDir() {
					loc = filepath.Join(loc, "index.apex")
				} else {
					loc += ".apex"
				}
			}
		}

		data, err := os.ReadFile(loc)
		if err != nil {
			return returnString(m, fmt.Sprintf("error: %v", err))
		}

		source := string(data)
		return returnString(m, source)
	}

	_, err = r.NewHostModuleBuilder("apex").
		ExportFunction("resolve", resolve,
			"resolve",
			"location_ptr", "location_len",
			"from_ptr", "from_len").
		Instantiate(ctx, r)
	if err != nil {
		panic(err)
	}

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		panic(err)
	}

	code, err := r.CompileModule(ctx, apexWasm)
	if err != nil {
		panic(err)
	}

	g, err := r.InstantiateModule(ctx, code, config)
	if err != nil {
		panic(err)
	}

	parse := g.ExportedFunction("parse")
	malloc = g.ExportedFunction("_malloc")
	free = g.ExportedFunction("_free")

	specSize := uint64(len(specBytes))

	results, err := malloc.Call(ctx, specSize)
	if err != nil {
		panic(err)
	}

	bufferPtr := results[0]
	defer free.Call(ctx, bufferPtr)

	g.Memory().Write(ctx, uint32(bufferPtr), specBytes)
	results, err = parse.Call(ctx, bufferPtr, specSize)
	if err != nil {
		panic(err)
	}

	ret := results[0]
	if ret == 0 {
		os.Exit(1)
		return
	}
	size := uint32(ret & 0xFFFFFFFF)
	ptr := uint32(ret >> uint64(32))

	docBytes, _ := g.Memory().Read(ctx, ptr, size)

	// doc, err := parser.Parse(parser.ParseParams{
	// 	Source: specBytes,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// docBytes, err := json.MarshalIndent(doc, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(string(docBytes))
}

func getHomeDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	home, err = homedir.Expand(home)
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".apex"), nil
}
