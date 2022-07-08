const fs = require('fs')
require("./wasm_exec")

let wasmModule = process.argv[2]
const go = new Go()

// 🖐️ hack for tiny go
go.importObject.env["syscall/js.finalizeRef"] = () => {}

let prm = WebAssembly.instantiate(fs.readFileSync(wasmModule), go.importObject)
.then(result => {
  go.run(result.instance)
  console.log(Handle("Jane Doe"))
  return result
})
.catch(error => {
  console.log("😡", error)
})


console.log(prm)



/*
http POST http://localhost:8080 "bob morane"
http --form POST http://localhost:8080  data="hello world"
curl -d "bob morane" -X POST  http://localhost:8080
*/
