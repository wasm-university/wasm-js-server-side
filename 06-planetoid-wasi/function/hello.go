package main

var buf [1024]byte

//Alloc(?)
//export getBuffer
func getBuffer() *byte {
	return &buf[0]
}

func main() {

}

//export Handle
func Handle(parameter string) *byte {

	var returnedValue [30]byte //arbitrary length

	copy(returnedValue[:], "👋 Hello :"+parameter)
	return &(returnedValue[0])
}

/*
curl -d "bob morane" -X POST  http://localhost:8080
*/
