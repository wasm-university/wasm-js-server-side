package main

import (
	"syscall/js"
)

func Handle(this js.Value, args []js.Value) interface{} {

	firstParam := args[0].String()

	return map[string]interface{}{
		"message": "👋 Hello :" + firstParam,
	}

}

func main() {

	js.Global().Set("Handle", js.FuncOf(Handle))

	<-make(chan bool)
}
