package main


var buf [1024]byte

//export getBuffer
func getBuffer() *byte {
   return &buf[0]
}

func main() {

}

//export Handle
func Handle(parameter string) *byte {

  var returnedValue [100]byte //arbitrary length

  copy(returnedValue[:], "👋 Hello :" + parameter)
  return &(returnedValue[0])
}
