package main

import (
	"reflect"
	"unsafe"
  "strconv"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

// log a message to the console using _log.
func log(message string) {
	ptr, size := stringToPtr(message)
	host_log(ptr, size)
}

// _log is a WebAssembly import which prints a string (linear memory offset,
// byteCount) to the console.
//
// Note: In TinyGo "//export" on a func is actually an import!
//go:wasm-module env


//export allocate_buffer
func allocateBuffer(size uint32) *byte {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of get_random_string.
	buf := make([]byte, size)
	return &buf[0]
}

//export host_get_string
func host_get_string(retBufPtr **byte, retBufSize *int)

// Get random string from the hosts.
func getString() string {
	var bufPtr *byte
	var bufSize int
	host_get_string(&bufPtr, &bufSize)
	//nolint
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(bufPtr)),
		Len:  uintptr(bufSize),
		Cap:  uintptr(bufSize),
	}))
}

//export host_log
func host_log(ptr uint32, size uint32)

//export host_fourtyTwo
func host_fourtyTwo() uint64


// helloWorld is a WebAssembly export that accepts a string pointer (linear memory
// offset) and calls greet.
//export helloWorld
func helloWorld(ptr, size uint32) {
	name := ptrToString(ptr, size)


	log("wasm >> 🖐️ hello world 🌍 " + name + " " +  strconv.FormatUint(host_fourtyTwo(), 10) + " " + getString())
}


// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	// Get a slice view of the underlying bytes in the stream. We use SliceHeader, not StringHeader
	// as it allows us to fix the capacity to what was allocated.
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size), // Tinygo requires these as uintptrs even if they are int fields.
		Cap:  uintptr(size), // ^^ See https://github.com/tinygo-org/tinygo/issues/1284
	}))
}

// stringToPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}