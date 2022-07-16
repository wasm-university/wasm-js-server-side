
"use strict";
const fs = require("fs");
const { WASI } = require("wasi");
const wasi = new WASI();
const importObject = { wasi_snapshot_preview1: wasi.wasiImport };

const http = require('http');
const port = process.env.PLANETOID_HTTP_PORT || 8080


function getMemoryAddressFor(text, moduleInstance) {

  // Get the address of the writable memory.
  let addr = moduleInstance.exports.getBuffer()
  let buffer = moduleInstance.exports.memory.buffer

  let mem = new Int8Array(buffer)
  let view = mem.subarray(addr, addr + text.length)

  for (let i = 0; i < text.length; i++) {
     view[i] = text.charCodeAt(i)
  }

  // Return the address we started at.
  return addr
}


(async () => {
  const wasm = await WebAssembly.compile(
    fs.readFileSync("./function/hello.wasm")
  );
  const moduleInstance = await WebAssembly.instantiate(wasm, importObject);

  wasi.start(moduleInstance);

  const requestHandler = (request, response) => {
    response.writeHead(200, {'Content-Type': 'application/json; charset=utf-8'})

    let body = ''
    request.on('data', chunk => {
      body += chunk.toString() // convert Buffer to string
    })
    request.on('end', () => {
      //console.log("body", body)
      //TODO: handle exceptions

      // call the function
      let handleValue = moduleInstance.exports.Handle(getMemoryAddressFor(body, moduleInstance), body.length)
      let memory = moduleInstance.exports.memory

      const buffer = new Uint8Array(memory.buffer, handleValue, 100)
      const str = new TextDecoder("utf8").decode(buffer)


      response.end(JSON.stringify(str))
    })

  }

  const server = http.createServer(requestHandler)

  server.listen(port, (err) => {
    if (err) {
      return console.log('😡 something bad happened', err)
    }
    console.log(`🌍 serving on ${port}`)
  })



})();

// $ node --experimental-wasi-unstable-preview1 index.js

