<script lang="ts">
    import { Editor } from "bytemd";
    import { afterUpdate, onMount } from "svelte";
    import { Button } from 'flowbite-svelte'
    import { GH_Helper as GH } from "$lib/github";
    import { userEmail, userName, SignIn } from "../stores";
    import { wasm_init } from "$lib/wasm-utils";
    import { plugins } from "$lib/editor-plugins";
    import { Octokit } from 'octokit'
    import "../app.css" // global styles

    let note = "";
    let SignedIn = false;
    function handleChange(e: { detail: { value: string } }) {
        note = e.detail.value;
    }

    afterUpdate(() => {
        // FEATURE: save notes to local storage
        localStorage.setItem("note", note);
    })

    onMount(async () => {
        // FEATURE: load notes from local storage
        note = localStorage.getItem("note") || "";
        
        // FEATURE: check if user is signed in
        const code = new URLSearchParams(window.location.search).get("code");

        if (code) {
            // get access token
            GH.getAccessTocken(code).then((res) => {
                if (!res.error) {
                if (res.access_token) {
                    // store token in local storage
                    localStorage.setItem("GITHUB_ACCESS_TOKEN", res.access_token);
                    SignIn.set(true);
                    } else {
                        console.error("Got neither error nor access token");
                    }
                } else {
                    console.error("AUTH_ERROR ? -> " + res.error);
                }
            });

            // remove code from url
            window.history.replaceState({}, document.title, "/");
        } else {
            // check if access token is present in local storage
            if (localStorage.getItem("GITHUB_ACCESS_TOKEN")) {
                SignIn.set(true);
            }
        }

        wasm_init();
    })

    userEmail.subscribe((value) => {
        console.log("userEmail", value);
    });
    SignIn.subscribe(async (value) => {
        SignedIn = value;
        
        if (SignedIn) {
            const octokit = new Octokit({
                auth: localStorage.getItem("GITHUB_ACCESS_TOKEN"),
            });

            const user = await octokit.rest.users.getAuthenticated().then((res) => {
                return res.data.login;
            });
            const id = await octokit.rest.users.getAuthenticated().then((res) => {
                return res.data.id
            });
            
            // no reply email
            const email = `${id}+${user}@users.noreply.github.com`
            userName.update(val => val = user);
            userEmail.set(email);
        } else {
            userName.update(val => val = "Guest");
            userEmail.set("")
        }
    });

</script>

<div>
    <Editor
        value={note}
        {plugins}
        on:change={handleChange}
    />

    {#if !SignedIn}
        <Button color="light" on:click={GH.SignInGitHub} class="m-2">
            <svg class="w-5 h-5 mr-1" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8C0 11.54 2.29 14.53 5.47 15.59C5.87 15.66 6.02 15.42 6.02 15.21C6.02 15.02 6.01 14.39 6.01 13.72C4 14.09 3.48 13.23 3.32 12.78C3.23 12.55 2.84 11.84 2.5 11.65C2.22 11.5 1.82 11.13 2.49 11.12C3.12 11.11 3.57 11.7 3.72 11.94C4.44 13.15 5.59 12.81 6.05 12.6C6.12 12.08 6.33 11.73 6.56 11.53C4.78 11.33 2.92 10.64 2.92 7.58C2.92 6.71 3.23 5.99 3.74 5.43C3.66 5.23 3.38 4.41 3.82 3.31C3.82 3.31 4.49 3.1 6.02 4.13C6.66 3.95 7.34 3.86 8.02 3.86C8.7 3.86 9.38 3.95 10.02 4.13C11.55 3.09 12.22 3.31 12.22 3.31C12.66 4.41 12.38 5.23 12.3 5.43C12.81 5.99 13.12 6.7 13.12 7.58C13.12 10.65 11.25 11.33 9.47 11.53C9.76 11.78 10.01 12.26 10.01 13.01C10.01 14.08 10 14.94 10 15.21C10 15.42 10.15 15.67 10.55 15.59C13.71 14.53 16 11.53 16 8C16 3.58 12.42 0 8 0Z" fill="white" />
            </svg>
            Sign In
        </Button>
    {:else}
        <Button color="light" on:click={GH.SignOutGitHub} class="m-2">
            <svg width="16" class="mr-2" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd" d="M2 2.75C2 1.7835 2.7835 1 3.75 1H6.25C6.66421 1 7 1.33579 7 1.75C7 2.16421 6.66421 2.5 6.25 2.5H3.75C3.61193 2.5 3.5 2.61193 3.5 2.75V13.25C3.5 13.3881 3.61193 13.5 3.75 13.5H6.25C6.66421 13.5 7 13.8358 7 14.25C7 14.6642 6.66421 15 6.25 15H3.75C2.7835 15 2 14.2165 2 13.25V2.75ZM12.4393 7.25H6.75002C6.33581 7.25 6.00002 7.58579 6.00002 8C6.00002 8.41422 6.33581 8.75 6.75002 8.75H12.4393L10.4697 10.7197C10.1768 11.0126 10.1768 11.4874 10.4697 11.7803C10.7626 12.0732 11.2374 12.0732 11.5303 11.7803L14.7803 8.53033C15.0732 8.23744 15.0732 7.76256 14.7803 7.46967L11.5303 4.21967C11.2374 3.92678 10.7626 3.92678 10.4697 4.21967C10.1768 4.51256 10.1768 4.98744 10.4697 5.28033L12.4393 7.25Z" fill="white" />
            </svg>
            Sign Out
        </Button>
    {/if}
</div>

<style>
    :global(.bytemd) {
        height: calc(100vh - 200px);
    }

    /* disable Github Permalink */
    :global(.bytemd-toolbar-right [bytemd-tippy-path='5']) {
      display: none;
    }
</style>
