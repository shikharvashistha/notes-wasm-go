<!-- <h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p> -->

<script lang="ts">
    import { afterUpdate, onMount } from 'svelte';
    import "$lib/wasm_exec.js"; // GLUE for go wasm
    import wasm from "$lib/main.wasm?url"; // WASM
    import { clientPub } from "../utils";
    import gfm from '@bytemd/plugin-gfm'
    import highlight from '@bytemd/plugin-highlight-ssr'
    import math from '@bytemd/plugin-math-ssr'
    import gemoji from '@bytemd/plugin-gemoji'
    import breaks from '@bytemd/plugin-breaks'
    import frontmatter from '@bytemd/plugin-frontmatter'
    import { Editor } from 'bytemd'

    import 'bytemd/dist/index.css'
    import 'github-markdown-css/github-markdown-light.css'
    import 'highlight.js/styles/atom-one-dark.css'
    import 'katex/dist/katex.css'
    
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

    function saveNote() {
        // TODO
        log();

        // @ts-ignore
        encryptNotes(note);

        // -> trigger AddNew [wasm function] function 
        //    with {note} and AccessToken_GH from local storage
        //    it takes care of encrypting notes and saving to Hub
        
        // @ts-ignore
        // Wasm Function
        AddNew(note, localStorage.getItem("AccessToken_GH"));
    }

    function log() {
        console.log(note);
    }

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

    /**
     * @param {String} code
     */
    async function getAccessTocken(code) {
        const ENDPOINT = "https://cint-proj-notes-frontend.vercel.app/api/getAuthCode"
        const res = await Promise.resolve(
            fetch(ENDPOINT, {
                method: 'POST',
                body: JSON.stringify({
                    "Auth": code
                })
            })
        )
        return res.json()
    }

    async function getUserData() {
        
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
        <p class="select-none antialiased my-auto text-3xl font-bold dark:text-white">
            Notes App 📝
        </p>
    
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
        <div class="">
                <Editor value={note} plugins={plugins} on:change={handleChange} />
        </div>
    
        <div class="text-center flex items-center justify-center flex-row">
            <button type="submit" on:click={saveNote} class="m-2 px-5 py-2.5 text-sm font-medium text-center text-white bg-blue-700 rounded-lg focus:ring-4 focus:ring-blue-200 dark:focus:ring-blue-900 hover:bg-blue-800">
                Add Note
            </button>
            {#if !SignedIn}
                <button type="button" on:click={SignIn} class="px-5 py-2.5 text-sm font-medium text-center text-white bg-blue-700 rounded-lg focus:ring-4 focus:ring-blue-200 dark:focus:ring-blue-900 hover:bg-blue-800">
                    Sign In Test
                </button>
            {/if}
            {#if SignedIn}
                <button type="button" on:click={SignOut} class="px-5 py-2.5 text-sm font-medium text-center text-white bg-blue-700 rounded-lg focus:ring-4 focus:ring-blue-200 dark:focus:ring-blue-900 hover:bg-blue-800">
                    Sign Out Test
                </button>
            {/if}
        </div>
    </form>
</div>
