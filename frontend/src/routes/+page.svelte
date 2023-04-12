<!-- <h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p> -->

<script>
    import markdownIt from "markdown-it";
    import sanitizeHtml from 'sanitize-html';
    import { afterUpdate, onMount } from 'svelte';
    import "$lib/wasm_exec.js"; // GLUE for go wasm
    import wasm from "$lib/main.wasm?url"; // WASM
    import { clientPub } from "../utils";
    import Fa from 'svelte-fa'
    import { faBold, faItalic, faLink, faTable, faQuoteLeft, faHeading } from '@fortawesome/free-solid-svg-icons'
    // @ts-ignore
    import { tooltip } from "@svelte-plugins/tooltips";

    // @ts-ignore
    const go = new Go();
    // @ts-ignore
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

        // -> trigger AddNew [wasm function] function 
        //    with {note} and AccessToken_GH from local storage
        //    it takes care of encrypting notes and saving to Hub
        
        // @ts-ignore
        // Wasm Function
        AddNew(note, localStorage.getItem("AccessToken_GH"));
    }

    function log() {
        console.log(note);
        console.log(result);
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
        // save note to local storage after every update
        localStorage.setItem("note", note);

        // convert markdown to html and store it in result
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
                        <u use:tooltip={{ content: 'Add a heading text.', position: 'bottom', autoPosition: true, align: 'center', animation: 'slide' }}>
                            <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                                <Fa icon={faHeading} />
                            </button>
                        </u>

                        <u use:tooltip={{ content: 'Add a bold text.', position: 'bottom', autoPosition: true, align: 'center', animation: 'slide' }}>
                            <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                                <Fa icon={faBold} />
                            </button>
                        </u>
                        
                        <u use:tooltip={{ content: 'Add a Italic text.', position: 'bottom', autoPosition: true, align: 'center', animation: 'slide' }}>
                            <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                                <Fa icon={faItalic} />
                            </button>
                        </u>

                        <u use:tooltip={{ content: 'Add a 2x2 table', position: 'bottom', autoPosition: true, align: 'center', animation: 'slide' }}>
                            <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                                <Fa icon={faTable} />
                            </button>
                        </u>

                        <u use:tooltip={{ content: 'Add a Blockquotes', position: 'bottom', autoPosition: true, align: 'center', animation: 'slide' }}>
                            <button type="button" class="p-2 text-gray-500 rounded cursor-pointer hover:text-gray-900 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-white dark:hover:bg-gray-600">
                                <Fa icon={faQuoteLeft} />
                            </button>
                        </u>
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
