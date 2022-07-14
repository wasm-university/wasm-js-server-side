package main

import (
	"fmt"
  "github.com/wasm-university/wasm-js-server-side/03-planetoid-wasi/function/helpers"

)


// Declare a main function, this is the entrypoint into our go module
// That will be run. In our example, we won't need this
func main() {
  fmt.Println("yo")
}


// This exports an add function.
// It takes in two 32-bit integer values
// And returns a 32-bit integer value.
// To make this function callable from JavaScript,
// we need to add the: "export add" comment above the function
//export add
func add(x int, y int) int {
  return x + y;
}

//export hello
func hello(parameters *int32) *byte {
	name := helpers.FromInt32PtrToString(parameters)
  fmt.Println(name)
	returnValue := "Hello " + name
	return helpers.FromStringToBytePtr(returnValue)
}


//export hey
func hey() *byte {
	returnValue := "Hello World"
	return helpers.FromStringToBytePtr(returnValue)
}
