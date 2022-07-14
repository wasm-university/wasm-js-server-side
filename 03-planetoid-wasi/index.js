"use strict";
const fs = require("fs");
const { WASI } = require("wasi");
const wasi = new WASI();
const importObject = { wasi_snapshot_preview1: wasi.wasiImport };

// Read a string from the instance's memory.
function readString(ptr, len, instance) {
  var m = new Uint8Array(instance.exports.memory.buffer, ptr, len);
  var decoder = new TextDecoder("utf-8");
  // return a slice of size `len` from the module's
  // memory, starting at offset `ptr`
  return decoder.decode(m.slice(0, len));
}

(async () => {
  const wasm = await WebAssembly.compile(
    fs.readFileSync("./function/hello.wasm")
  );
  const instance = await WebAssembly.instantiate(wasm, importObject);

  wasi.start(instance);
  //console.log(instance)
  console.log(instance.exports.add(1,23))
  //console.log(instance.exports.hello("Sam")) // BigInt

  let helloValue = instance.exports.hello("Sam")
  let heyValue = instance.exports.hey()
  let memory = instance.exports.memory

  console.log(memory.buffer)


  //const values = new Uint32Array(memory.buffer);
  //console.log(values[0]);

  const buffer = new Uint8Array(memory.buffer, heyValue, 11)
  const str = new TextDecoder("utf8").decode(buffer)
  console.log(`📝: ${str}`)

  const buffer2 = new Uint8Array(memory.buffer, helloValue, 9)
  const str2 = new TextDecoder("utf8").decode(buffer2)
  console.log(`📝: ${str2}`)



})();

// $ node --experimental-wasi-unstable-preview1 index.js

/*
  - read the memory
  - get buffer at the position of the value
  - transform the buffer to string
*/

/*
const greetValue = instance.exports.greet();
const memory = instance.exports.memory;
const buffer = new Uint8Array(memory.buffer, greetValue, 12);
const str = new TextDecoder("utf8").decode(buffer);
console.log(`📝: ${str}`)
*/
