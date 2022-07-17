package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
  "unsafe"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().
		// Enable WebAssembly 2.0 support, which is required for TinyGo 0.24+.
		WithWasmCore2())
	defer r.Close(ctx) // This closes everything this Runtime created.

  // === Host Function ===
	// Instantiate a Go-defined module named "env" that exports a function to
	// log to the console.
	_, err := r.NewModuleBuilder("env").
		ExportFunction("host_log", logString).
    ExportFunction("host_fourtyTwo", func() uint64 {
      return uint64(42)
    }).
    ExportFunction("host_tada", func() uint32 {
      ptr, _ := stringToPtr("tada")
      return ptr
    }).
		Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}

	// Note: testdata/greet.go doesn't use WASI, but TinyGo needs it to
	// implement functions such as panic.
	if _, err = wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Panicln(err)
	}

	// Instantiate a WebAssembly module that imports the "log" function defined
	// in "env" and exports "memory" and functions we'll use in this example.
  helloWasm, err := os.ReadFile("./function/hello.wasm")
  if err != nil {
      //fmt.Println(err)
      log.Panicln(err)
      return
  }


	mod, err := r.InstantiateModuleFromBinary(ctx, helloWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get references to WebAssembly functions we'll use in this example.
	helloWorld := mod.ExportedFunction("helloWorld")
	//🖐️ greeting := mod.ExportedFunction("greeting")
	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	// Let's use the argument to this main function in Wasm.
	name := "Bob Morane"
	nameSize := uint64(len(name))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	namePtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, namePtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(ctx, uint32(namePtr), []byte(name)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size(ctx))
	}

	// Now, we can call "greet", which reads the string we wrote to memory!
	_, err = helloWorld.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Panicln(err)
	}

	// Finally, we get the greeting message "greet" printed. This shows how to
	// read-back something allocated by TinyGo.
	//🖐️ ptrSize, err := greeting.Call(ctx, namePtr, nameSize)
	//🖐️ if err != nil {
	//🖐️ 	log.Panicln(err)
	//🖐️ }
	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	//🖐️ greetingPtr := uint32(ptrSize[0] >> 32)
	//🖐️ greetingSize := uint32(ptrSize[0])
	// The pointer is a linear memory offset, which is where we write the name.
	//🖐️ if bytes, ok := mod.Memory().Read(ctx, greetingPtr, greetingSize); !ok {
	//🖐️ 	log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
	//🖐️ 		greetingPtr, greetingSize, mod.Memory().Size(ctx))
	//🖐️ } else {
	//🖐️ 	fmt.Println("go >>", string(bytes))
	//🖐️ }
}

func logString(ctx context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(ctx, offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}
