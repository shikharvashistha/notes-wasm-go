<!-- <h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p> -->

<script>
    import markdownIt from "markdown-it";
    import sanitizeHtml from 'sanitize-html';
    import { afterUpdate, onMount } from 'svelte';
    import "$lib/wasm_exec.js"; // GLUE for go wasm
    import wasm from "$lib/main.wasm?url"; // WASM
    import { clientPub } from "../utils";

    // @ts-ignore
    const go = new Go();
    let mod, inst, code;
    let user

    let note = "";
    let result = "";
    let SignedIn = false;

    const md = markdownIt({
      html: true,
      linkify: true,
      typographer: true
    });

    import "../app.css";

    function saveNote() {
        // TODO
        log();

        // @ts-ignore
        encryptNotes(note);
    }

    function log() {
        console.log(note);
        console.log(result);
    }

    function SignIn() {
        // TODO: chnage scope to "repo"
        window.location.assign("https://github.com/login/oauth/authorize?client_id=" + clientPub.clientID)
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

        WebAssembly.instantiateStreaming(fetch(wasm), go.importObject).then(async (result) => {
        mod     = result.module;
        inst    = result.instance;
        await go.run(inst);
    });
    })

    afterUpdate(() => {
		result = sanitizeHtml(md.render(note), {
			allowedTags: sanitizeHtml.defaults.allowedTags.concat([ 'h1', 'h2', 'img' ])
		});
	});

</script>

<div class="my-auto p-16 rounted-lg text-center grid gap-4">

    <div class="text-center">
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
        <div class="w-full mb-4 border border-gray-200 rounded-lg bg-gray-50 dark:bg-gray-700 dark:border-gray-600">
            <div class="flex items-center justify-between px-3 py-2 border-b dark:border-gray-600">
                <div class="flex flex-wrap items-center divide-gray-200 sm:divide-x dark:divide-gray-600">
                    <div class="flex items-center space-x-1 sm:pr-4">
                        <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M8 4a3 3 0 00-3 3v4a5 5 0 0010 0V7a1 1 0 112 0v4a7 7 0 11-14 0V7a5 5 0 0110 0v4a3 3 0 11-6 0V7a1 1 0 012 0v4a1 1 0 102 0V7a3 3 0 00-3-3z" clip-rule="evenodd"></path></svg>
                            <span class="sr-only">Attach file</span>
                        </button>
                        <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"></path></svg>
                            <span class="sr-only">Embed map</span>
                        </button>
                        <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"></path></svg>
                            <span class="sr-only">Upload image</span>
                        </button>
                        <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
                             <span class="sr-only">Format code</span>
                        </button>
                        <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                            <svg aria-hidden="true" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM7 9a1 1 0 100-2 1 1 0 000 2zm7-1a1 1 0 11-2 0 1 1 0 012 0zm-.464 5.535a1 1 0 10-1.415-1.414 3 3 0 01-4.242 0 1 1 0 00-1.415 1.414 5 5 0 007.072 0z" clip-rule="evenodd"></path></svg>
                            <span class="sr-only">Add emoji</span>
                        </button>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-2">
                <div class="bg-white rounded-b-lg dark:bg-gray-800">
                    <label for="editor" class="sr-only">Publish post</label>
                    <textarea bind:value={note} id="editor" rows="18" 
                        class="outline-none block w-full px-0 text-sm text-gray-800 bg-white border-1 dark:bg-gray-800 dark:text-white dark:placeholder-gray-400" placeholder="Write an article..." required></textarea>
                </div>
    
                <div class="prose dark:prose-invert max-w-none text-left w-full py-4 overflow-scroll h-full dark:text-white border border-gray-200 rounded-lg bg-gray-50 dark:bg-gray-700 dark:border-gray-600">
                    {@html result}
                </div>
            </div>

        </div>
    
        <div class="text-center">
            <button type="submit" on:click={saveNote} class="px-5 py-2.5 text-sm font-medium text-center text-white bg-blue-700 rounded-lg focus:ring-4 focus:ring-blue-200 dark:focus:ring-blue-900 hover:bg-blue-800">
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
