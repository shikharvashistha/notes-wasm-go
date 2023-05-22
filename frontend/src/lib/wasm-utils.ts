/*
    WASM SECTION
*/

import "$lib/wasm_exec.js";
import wasm from "$lib/main.wasm?url";

const wasm_init = async () => {
  // @ts-ignore
  const go = new Go(); // defined in wasm_exec.js
  let mod, inst;

  WebAssembly.instantiateStreaming(fetch(wasm), go.importObject).then(
    async (result) => {
      mod = result.module;
      inst = result.instance;
      await go.run(inst);
    }
  );

  console.log("WASM initialized");
};

export { wasm_init };
