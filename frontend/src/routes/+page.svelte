<!-- <h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p> -->

<script lang="ts">
    import { afterUpdate, onMount } from 'svelte';
    import "$lib/wasm_exec.js"; // GLUE for go wasm
    import wasm from "$lib/main.wasm?url"; // WASM
    import { clientPub, spice } from "../utils";
    import gfm from '@bytemd/plugin-gfm'
    import highlight from '@bytemd/plugin-highlight-ssr'
    import math from '@bytemd/plugin-math-ssr'
    import gemoji from '@bytemd/plugin-gemoji'
    import breaks from '@bytemd/plugin-breaks'
    import frontmatter from '@bytemd/plugin-frontmatter'
    import { Editor } from 'bytemd'
    import { Button, DarkMode, Heading, P, Span } from 'flowbite-svelte'

    import 'bytemd/dist/index.css'
    import 'github-markdown-css/github-markdown-light.css'
    import 'highlight.js/styles/atom-one-dark.css'
    import 'katex/dist/katex.css'

    import "../wasm-utils.ts"
    
    // @ts-ignore
    const go = new Go();
    // @ts-ignore
    let mod, inst, code;

    let note = "";
    let SignedIn = false;
    let localeKey = 'en'

    
    function stripPrefixes(obj: Record<string, any>) {
        return Object.entries(obj).reduce((p, [key, value]) => {
          p[key.split('/').slice(-1)[0].replace('.json', '')] = value
          // console.log(p)
          return p
        }, {} as Record<string, any>)
    }
    const gfmLocales = stripPrefixes(
        import.meta.glob('/node_modules/@bytemd/plugin-gfm/locales/*.json', {
          eager: true,
        })
    )
    const mathLocales = stripPrefixes(
        import.meta.glob('/node_modules/@bytemd/plugin-math/locales/*.json', {
          eager: true,
        })
    )

    const plugins = [
        gfm({
            locale: gfmLocales[localeKey],
        }),
        highlight(),
        math({
            locale: mathLocales[localeKey],
            katexOptions: { output: 'html' }, // https://github.com/KaTeX/KaTeX/issues/2796
        }),
        gemoji(),
        breaks(),
        frontmatter(),
    ]

    /**
     * @param {{ detail: { value: any; }; }} e
     */
    function handleChange(e: { detail: { value: string; }; }) {
        note = e.detail.value
    }

    import "../app.css";

    function SignIn() {
        window.location.assign("https://github.com/login/oauth/authorize?client_id=" + clientPub.clientID+ "&scope=repo")
    }

    function SignOut() {
        // remove token from local storage
        localStorage.removeItem("AccessToken_GH");
        SignedIn = false;

        // remove ?code= from url
        window.history.replaceState({}, document.title, "/");
    }

    function TsaveNote() {
        // strategy 1 - LOCAL
        // - create a unique id for the note
        // - save the note to the local storage
        // - save the id to the local storage with key "storate-id-array"
        //
        // state 2 - REMOTE
        // TODO ?
        //


        var id = Math.floor(Math.random() * 1000000000);
        var encrypted = ""
        encryptHandler(note).then((res) => {
            localStorage.setItem("storage-"+String(id), res);
            encrypted = res
        })
       
        if (localStorage.getItem("storage-id-array") == null) {
            localStorage.setItem("storage-id-array", String(id));
        } else {
            localStorage.setItem("storage-id-array", localStorage.getItem("storage-id-array") + "," + String(id));
        }

        if (SignedIn) {
            // upload to remove
            var url = "http://localhost:8081/?https://github.com/SaicharanKandukuri/test-re"
            // @ts-ignore
            var AccessTocken: string = localStorage.getItem("AccessToken_GH")

            gitClone(url).then((res) => {
                console.log(res)
                writeNoteToFile("wasm-repo/storage-" + String(id) , encrypted).then((res) => {
                    console.log(res)

                    gitPush(url,
                            AccessTocken,
                            "GitHub Action",
                            "action@github.com",
                            "storage-"+String(id),
                            "WasmUpload JOB").then((res) => {
                                console.log(res)
                            })
                })
            })
        }
    }

    function wipeNotes() {
        // - get the id array from local storage
        // - loop through the array and remove all notes from local storage
        // - remove the id array from local storage

        if (localStorage.getItem("storage-id-array") == null) {
            return;
        }

        // @ts-ignore
        var idArray = localStorage.getItem("storage-id-array").split(",");
        for (var i = 0; i < idArray.length; i++) {
            localStorage.removeItem("storage-"+idArray[i]);
        }
        localStorage.removeItem("storage-id-array");
    }

    function trigger() {
        // @ts-ignore
        getAccessTocken(code).then((res) => {
            if (!res.error) {
                if (res.access_token) {
                    // store token in local storage
                    localStorage.setItem("AccessToken_GH", res.access_token);
                } else {
                    console.error("Got neither error nor access token")
                }
            } else {
                console.error("AUTH_ERROR ? -> "+ res.error)
            }
        })
    }

    onMount(async () => {
        const LurlParms = new URLSearchParams(window.location.search);
        code = LurlParms.get("code");

        if (!localStorage.getItem("AccessToken_GH")) {
            if (code != null) {
                trigger();
            }
        } else {
            if (code != null) {
                window.history.replaceState({}, document.title, "/");
            }
            SignedIn = true;
        }

        // if note is empty onMount -> try to get it from local storage
        if (note == "") {
            // @ts-ignore
            if (localStorage.getItem("note")) {
                // @ts-ignore
                note = localStorage.getItem("note");
            }
        }

        WebAssembly.instantiateStreaming(fetch(wasm), go.importObject).then(async (result) => {
        mod     = result.module;
        inst    = result.instance;
        await go.run(inst);
    });
    })

    afterUpdate(() => {
        localStorage.setItem("note", note);
	});


</script>

<div class="my-auto h-screen p-12 rounted-lg flex flex-col">
    <div class="text-center mb-4">
        <Heading tag="h1" class="mb-4 antialiased" >Notes App 📝</Heading>   
        {#if !SignedIn}
            <p class="select-none antialiased text-sl dark:text-slate-50">
                A note implementation using  
                <a class="underline decoration-pink-500" href="https://svelte.dev/">SvelteKit</a>
                and
                <a class="underline decoration-sky-500" href="https://webassembly.org/">WASM</a>
            </p>
        {/if}
    
    </div>
    
    <form>
        <div>
            <Editor value={note} plugins={plugins} on:change={handleChange} />
        </div>
    
        <div class="text-center flex items-center justify-center flex-row">
            
            {#if !SignedIn}
            <Button color="light" on:click={SignIn} class="m-2" >
                <svg class="w-5 h-5 mr-1" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8C0 11.54 2.29 14.53 5.47 15.59C5.87 15.66 6.02 15.42 6.02 15.21C6.02 15.02 6.01 14.39 6.01 13.72C4 14.09 3.48 13.23 3.32 12.78C3.23 12.55 2.84 11.84 2.5 11.65C2.22 11.5 1.82 11.13 2.49 11.12C3.12 11.11 3.57 11.7 3.72 11.94C4.44 13.15 5.59 12.81 6.05 12.6C6.12 12.08 6.33 11.73 6.56 11.53C4.78 11.33 2.92 10.64 2.92 7.58C2.92 6.71 3.23 5.99 3.74 5.43C3.66 5.23 3.38 4.41 3.82 3.31C3.82 3.31 4.49 3.1 6.02 4.13C6.66 3.95 7.34 3.86 8.02 3.86C8.7 3.86 9.38 3.95 10.02 4.13C11.55 3.09 12.22 3.31 12.22 3.31C12.66 4.41 12.38 5.23 12.3 5.43C12.81 5.99 13.12 6.7 13.12 7.58C13.12 10.65 11.25 11.33 9.47 11.53C9.76 11.78 10.01 12.26 10.01 13.01C10.01 14.08 10 14.94 10 15.21C10 15.42 10.15 15.67 10.55 15.59C13.71 14.53 16 11.53 16 8C16 3.58 12.42 0 8 0Z" fill="white"/>
                </svg>
                Sign In
            </Button>
            {/if}

            {#if SignedIn}
                <Button color="light" on:click={SignOut} class="m-2" >
                    <svg width="16" class="mr-2" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path fill-rule="evenodd" clip-rule="evenodd" d="M2 2.75C2 1.7835 2.7835 1 3.75 1H6.25C6.66421 1 7 1.33579 7 1.75C7 2.16421 6.66421 2.5 6.25 2.5H3.75C3.61193 2.5 3.5 2.61193 3.5 2.75V13.25C3.5 13.3881 3.61193 13.5 3.75 13.5H6.25C6.66421 13.5 7 13.8358 7 14.25C7 14.6642 6.66421 15 6.25 15H3.75C2.7835 15 2 14.2165 2 13.25V2.75ZM12.4393 7.25H6.75002C6.33581 7.25 6.00002 7.58579 6.00002 8C6.00002 8.41422 6.33581 8.75 6.75002 8.75H12.4393L10.4697 10.7197C10.1768 11.0126 10.1768 11.4874 10.4697 11.7803C10.7626 12.0732 11.2374 12.0732 11.5303 11.7803L14.7803 8.53033C15.0732 8.23744 15.0732 7.76256 14.7803 7.46967L11.5303 4.21967C11.2374 3.92678 10.7626 3.92678 10.4697 4.21967C10.1768 4.51256 10.1768 4.98744 10.4697 5.28033L12.4393 7.25Z" fill="white"/>
                    </svg>    
                    Sign Out
                </Button>
            {/if}
            
            <Button color="purple" on:click={wipeNotes} class="m-2" >
                <svg width="16" class="mr-2" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M3.75 1.5C3.61193 1.5 3.5 1.61193 3.5 1.75V14.25C3.5 14.3881 3.61193 14.5 3.75 14.5H13.25C13.3881 14.5 13.5 14.3881 13.5 14.25V4.66421C13.5 4.59791 13.4737 4.53432 13.4268 4.48744L10.5126 1.57322C10.4657 1.52634 10.4021 1.5 10.3358 1.5H3.75ZM2 1.75C2 0.783502 2.7835 0 3.75 0H10.3358C10.7999 0 11.245 0.184374 11.5732 0.512563L14.4874 3.42678C14.8156 3.75497 15 4.20008 15 4.66421V14.25C15 15.2165 14.2165 16 13.25 16H3.75C2.7835 16 2 15.2165 2 14.25V1.75ZM8.25 7.5H10.4922C10.9064 7.5 11.2422 7.83578 11.2422 8.25C11.2422 8.66421 10.9064 9 10.4922 9H8.25254L5.99781 9.01498C5.58361 9.01773 5.2456 8.68418 5.24284 8.26998C5.24009 7.85577 5.57364 7.51776 5.98785 7.51501L8.25 7.5Z" fill="white"/>
                </svg>

                Delete
            </Button>  
            <Button color="green" on:click={TsaveNote} class="m-2" >
                <svg width="16" class="mr-2" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M3.75 1.5C3.61193 1.5 3.5 1.61193 3.5 1.75V14.25C3.5 14.3881 3.61193 14.5 3.75 14.5H13.25C13.3881 14.5 13.5 14.3881 13.5 14.25V4.66421C13.5 4.59791 13.4737 4.53432 13.4268 4.48744L10.5126 1.57322C10.4657 1.52634 10.4021 1.5 10.3358 1.5H3.75ZM2 1.75C2 0.783502 2.7835 0 3.75 0H10.3358C10.7999 0 11.245 0.184374 11.5732 0.512563L14.4874 3.42678C14.8156 3.75497 15 4.20008 15 4.66421V14.25C15 15.2165 14.2165 16 13.25 16H3.75C2.7835 16 2 15.2165 2 14.25V1.75ZM8.22999 5.25785C8.64419 5.25508 8.98222 5.58861 8.98499 6.00281L8.995 7.5H10.4922C10.9064 7.5 11.2422 7.83579 11.2422 8.25C11.2422 8.66421 10.9064 9 10.4922 9H9V10.5072C9 10.9214 8.66421 11.2572 8.25 11.2572C7.83579 11.2572 7.5 10.9214 7.5 10.5072V9.005L5.99782 9.01498C5.58361 9.01773 5.2456 8.68419 5.24285 8.26998C5.2401 7.85578 5.57365 7.51777 5.98785 7.51501L7.495 7.505L7.48502 6.01285C7.48225 5.59864 7.81578 5.26062 8.22999 5.25785Z" fill="white"/>
                </svg>
                    
                Save Notes
            </Button>
            
            <DarkMode class="text-lg">
                <svelte:fragment slot="lightIcon">
                    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" preserveAspectRatio="xMidYMid meet" viewBox="0 0 32 32"><path fill="currentColor" d="M16 12.005a4 4 0 1 1-4 4a4.005 4.005 0 0 1 4-4m0-2a6 6 0 1 0 6 6a6 6 0 0 0-6-6ZM5.394 6.813L6.81 5.399l3.505 3.506L8.9 10.319zM2 15.005h5v2H2zm3.394 10.193L8.9 21.692l1.414 1.414l-3.505 3.506zM15 25.005h2v5h-2zm6.687-1.9l1.414-1.414l3.506 3.506l-1.414 1.414zm3.313-8.1h5v2h-5zm-3.313-6.101l3.506-3.506l1.414 1.414l-3.506 3.506zM15 2.005h2v5h-2z"/></svg>
                </svelte:fragment>
                <svelte:fragment slot="darkIcon">
                    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" preserveAspectRatio="xMidYMid meet" viewBox="0 0 32 32"><path fill="currentColor" d="M13.502 5.414a15.075 15.075 0 0 0 11.594 18.194a11.113 11.113 0 0 1-7.975 3.39c-.138 0-.278.005-.418 0a11.094 11.094 0 0 1-3.2-21.584M14.98 3a1.002 1.002 0 0 0-.175.016a13.096 13.096 0 0 0 1.825 25.981c.164.006.328 0 .49 0a13.072 13.072 0 0 0 10.703-5.555a1.01 1.01 0 0 0-.783-1.565A13.08 13.08 0 0 1 15.89 4.38A1.015 1.015 0 0 0 14.98 3Z"/></svg>
                </svelte:fragment>
            </DarkMode>
        </div>
    </form>
</div>

<style>
:global(.bytemd) {
  height: calc(100vh - 200px);
}
:global(.medium-zoom-overlay) {
  z-index: 100;
}
:global(.medium-zoom-image--opened) {
  z-index: 101;
}
</style>
