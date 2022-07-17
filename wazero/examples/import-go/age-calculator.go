package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tetratelabs/wazero"
)

// ageCalculatorWasm was generated by the following:
//	cd testdata; wat2wasm --debug-names age_calculator.wat
//go:embed testdata/age_calculator.wasm
var ageCalculatorWasm []byte

// main shows how to define, import and call a Go-defined function from a
// WebAssembly-defined function.
//
// See README.md for a full description.
func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime()
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate a Go-defined module named "env" that exports functions to
	// get the current year and log to the console.
	//
	// Note: As noted on ExportFunction documentation, function signatures are
	// constrained to a subset of numeric types.
	// Note: "env" is a module name conventionally used for arbitrary
	// host-defined functions, but any name would do.
	_, err := r.NewModuleBuilder("env").
		ExportFunction("log_i32", func(v uint32) {
			fmt.Println("log_i32 >>", v)
		}).
		ExportFunction("current_year", func() uint32 {
			if envYear, err := strconv.ParseUint(os.Getenv("CURRENT_YEAR"), 10, 64); err == nil {
				return uint32(envYear) // Allow env-override to prevent annual test maintenance!
			}
			return uint32(time.Now().Year())
		}).
		Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}

	// Instantiate a WebAssembly module named "age-calculator" that imports
	// functions defined in "env".
	//
	// Note: The import syntax in both Text and Binary format is the same
	// regardless of if the function was defined in Go or WebAssembly.
	ageCalculator, err := r.InstantiateModuleFromBinary(ctx, ageCalculatorWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Read the birthYear from the arguments to main
	birthYear, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Panicf("invalid arg %v: %v", os.Args[1], err)
	}

	// First, try calling the "get_age" function and printing to the console externally.
	results, err := ageCalculator.ExportedFunction("get_age").Call(ctx, birthYear)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("println >>", results[0])

	// First, try calling the "log_age" function and printing to the console externally.
	_, err = ageCalculator.ExportedFunction("log_age").Call(ctx, birthYear)
	if err != nil {
		log.Panicln(err)
	}
}