package main

import (
	"github.com/suborbital/reactr/api/tinygo/runnable" 
)

type Hello struct{}

func (h Hello) Run(input []byte) ([]byte, error) {
	return []byte("Hello, " + string(input)), nil
}

// initialize runnable, do not edit //
func main() {
	runnable.Use(Hello{})
}
