/*
    WASM SECTION
*/

import "$lib/wasm_exec.js"
import wasm from "$lib/main.wasm?url";

declare function git_clone(url: string): Promise<string>;
declare function touchNcat(path: string, content: string): Promise<string>;

interface push {
    url: string;
    AccessTocken: string;
    userName: string;
    userEmail: string;
    fileName: string;
    commitMessage: string;
}
declare function git_push(push: push): Promise<string>;

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
}

export { wasm_init };
