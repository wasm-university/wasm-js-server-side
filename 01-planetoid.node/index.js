const http = require('http')
const fs = require('fs/promises')

require("./wasm_exec")

const port = process.env.PLANETOID_HTTP_PORT || 8080
let wasmModule = process.argv[2]

function runWasm(wasmFile, args) {
  const go = new Go()

  // 🖐️ hack for tiny go
  go.importObject.env["syscall/js.finalizeRef"] = () => {}

  return new Promise((resolve, reject) => {
    WebAssembly.instantiate(wasmFile, go.importObject)
    .then(result => {
      if(args) go.argv = args
      go.run(result.instance)
      resolve(result.instance)
    })
    .catch(error => {
      reject(error)
    })
  })
}

const requestHandler = (request, response) => {
  response.writeHead(200, {'Content-Type': 'application/json; charset=utf-8'})

  let body = ''
  request.on('data', chunk => {
    body += chunk.toString() // convert Buffer to string
  })
  request.on('end', () => {
    //console.log("body", body)
    //TODO: handle exceptions
    response.end(JSON.stringify(Handle(body)))
  })

}

fs.readFile(wasmModule)
  .then(wasmFile => runWasm(wasmFile))
  .then(wasm => {

    //console.log(Handle("Jane Doe"))

    console.log("wasm module loaded")

    const server = http.createServer(requestHandler)

    server.listen(port, (err) => {
      if (err) {
        return console.log('😡 something bad happened', err)
      }
      console.log(`🌍 serving on ${port}`)
    })




  })
  .catch(error => {
    console.log("ouch", error)
  })

/*
http POST http://localhost:8080 "bob morane"
http --form POST http://localhost:8080  data="hello world"
curl -d "bob morane" -X POST  http://localhost:8080
*/
