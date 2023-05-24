<script lang="ts">
  import { Editor } from "bytemd";
  import { afterUpdate, onMount } from "svelte";
  import { Button } from "flowbite-svelte";
  import { GH_Helper as GH } from "$lib/github";
  import { userEmail, userName, SignIn } from "../stores";
  import { wasm_init } from "$lib/wasm-utils";
  import { plugins } from "$lib/editor-plugins";
  import { Octokit } from "octokit";
  import { spice } from "../utils";
  import { Input } from 'flowbite-svelte';
  import { Repo, Branch } from "../repo.json";
  import "../app.css"; // global styles

  let note = "";
  let noteName = "";
  let SignedIn = false;
  let user = "";
  let user_email = "";

  function handleChange(e: { detail: { value: string } }) {
    note = e.detail.value;
  }

  afterUpdate(() => {
    // FEATURE: save notes to local storage
    localStorage.setItem("note", note);
  });

  let clonedOnce = false
  async function saveNote() {
    // const url = String(clientPub.repoURL) <= this wont work (e: not implemented on js at syscall/js.valueNew (was)
    // using json as a workaround
    
      if (SignedIn) {
      const url = "https://cors.isomorphic-git.org/github.com/"+Repo
      if (!clonedOnce) {
        //@ts-ignore
        const clone = await git_clone(url)
        clonedOnce = true
      } else {
        const clone = "Already cloned"
      }
      const filteredNoteName = noteName.replace(/[^a-zA-Z0-9]/g, "")
      const fileName = filteredNoteName+".jaef"
      
      //@ts-ignore
      const encryptNote = await encrypt_text(note, spice.encryptSecret)
      //@ts-ignore
      const write = await touchNcat("wasm-repo/"+fileName, encryptNote)
      // get history
      const history = await getHistory()
      const newHistory = JSON.parse(history)
      // remove all special characters from noteName
      newHistory.push({
        name: noteName,
        fileName: filteredNoteName+".jaef",
        user: user,
        date: new Date().toISOString(),
      })

      //@ts-ignore
      const writeNewHistoryFile = await touchNcat("wasm-repo/history.json", JSON.stringify(newHistory))
      
      //@ts-ignore
      const push = await git_push(
        url,
        localStorage.getItem("GITHUB_ACCESS_TOKEN"),
        user, user_email,
        fileName,
        "WASM Commit: Added "+fileName // jaef -> just an encrypted file
      )

      //@ts-ignore
      const pushHistory = await git_push(
        url,
        localStorage.getItem("GITHUB_ACCESS_TOKEN"),
        user, user_email,
        "history.json",
        "WASM Commit: history.json" // jaef -> just an encrypted file
      )
      console.log("WASM "  + push)
      console.log("WASM "  + pushHistory)
    } else {
      console.warn("You are not signed in");
    }

    // finally fetch history
    await fetchHistory()
  }

  async function getHistory() {
    const url = "https://raw.githubusercontent.com/"+Repo+"/"+Branch+"/history.json"
    const repourl = "http://localhost:8081/?https://github.com/"+Repo
      
    // fetch file
    /// if error, return
    const file = await fetch(url).then(async (res) => {
      if (res.ok) {
        return res.text()
      } else {
        if (!clonedOnce) {
          //@ts-ignore
          const clone = await git_clone(url)
        }
        //@ts-ignore
        const writeNewHistoryFile = await touchNcat("wasm-repo/history.json", "[]")
        //@ts-ignore
        const push = await git_push(
          repourl,
          localStorage.getItem("GITHUB_ACCESS_TOKEN"),
          user, user_email,
          "history.json",
          "WASM Commit: Added history.json" // jaef -> just an encrypted file
        )
        return "[]"
      }
    })
    console.log(file)
    return file
  }
  async function fetchHistory() {
    if (SignedIn) {
      const rawUrl = "https://raw.githubusercontent.com/" + Repo + "/"+ Branch +"/"
      const res = await fetch(rawUrl + "history.json", {cache: "no-cache"})
      const data = await res.json()
      if (data.length == 0) {
        emptyHistory = true
      }
      console.log(data)
      historyData = data
    }
  }

  // // HISTORY IMP
  import { Heading, Badge } from 'flowbite-svelte'
  let historyData: any = [{ name: "Loading...", fileName: "Loading...", user: "Loading...", date: "Loading..."}]
  let emptyHistory=false
  async function handleHistoryClick(filename: string) {
    console.log(filename)
    const encryptedNote = await GH.getFileContents(filename)
    //@ts-ignore
    const decryptedNote = await decrypt_text(encryptedNote, spice.encryptSecret)
    note = decryptedNote
    noteName = filename.split(".")[0]
  }

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

    // History
    fetchHistory()
    wasm_init();
  });

  userEmail.subscribe((value) => {
    user_email = value;
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
        return res.data.id;
      });

      // no reply email
      const email = `${id}+${user}@users.noreply.github.com`;
      userName.update((val) => (val = user));
      userEmail.set(email);
    } else {
      userName.update((val) => (val = "Guest"));
      userEmail.set("");
    }

    // trigger history fetch
    fetchHistory()
  });
  userName.subscribe((value) => {
    user = value;
  });
</script>

<div>
  <Editor value={note}
  uploadImages={(files) => {
    return Promise.all(
      files.map(async (file) => {

        if (!SignedIn) {
          return {
            url: "https://i.imgur.com/IYn3ARk.jpg",
          }
        }

        function fileToBase64(file) {
          return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;
            reader.readAsDataURL(file);
          });
        }

        const res = await fileToBase64(file)
        //@ts-ignore
        const base64 = res.split(",")[1]
        const uploadRes = await GH.uploadFile(base64, file.name, user, user_email)

        return {
          url: uploadRes,
        }
      })
    )
  }}
  {plugins} on:change={handleChange} />

  {#if !SignedIn}
    <Button color="light" on:click={GH.SignInGitHub} class="m-2">
      <svg
        class="w-5 h-5 mr-1"
        viewBox="0 0 16 16"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M8 0C3.58 0 0 3.58 0 8C0 11.54 2.29 14.53 5.47 15.59C5.87 15.66 6.02 15.42 6.02 15.21C6.02 15.02 6.01 14.39 6.01 13.72C4 14.09 3.48 13.23 3.32 12.78C3.23 12.55 2.84 11.84 2.5 11.65C2.22 11.5 1.82 11.13 2.49 11.12C3.12 11.11 3.57 11.7 3.72 11.94C4.44 13.15 5.59 12.81 6.05 12.6C6.12 12.08 6.33 11.73 6.56 11.53C4.78 11.33 2.92 10.64 2.92 7.58C2.92 6.71 3.23 5.99 3.74 5.43C3.66 5.23 3.38 4.41 3.82 3.31C3.82 3.31 4.49 3.1 6.02 4.13C6.66 3.95 7.34 3.86 8.02 3.86C8.7 3.86 9.38 3.95 10.02 4.13C11.55 3.09 12.22 3.31 12.22 3.31C12.66 4.41 12.38 5.23 12.3 5.43C12.81 5.99 13.12 6.7 13.12 7.58C13.12 10.65 11.25 11.33 9.47 11.53C9.76 11.78 10.01 12.26 10.01 13.01C10.01 14.08 10 14.94 10 15.21C10 15.42 10.15 15.67 10.55 15.59C13.71 14.53 16 11.53 16 8C16 3.58 12.42 0 8 0Z"
          fill="white"
        />
      </svg>
      Sign In
    </Button>
  {:else}
    <div class="flex flex-col text-center justify-center items-center md:flex-row m-3">
      <div>
        <Button color="light" on:click={GH.SignOutGitHub} class="ml-2 mt-2">
          <svg
            width="16"
            class="mr-2"
            height="16"
            viewBox="0 0 16 16"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M2 2.75C2 1.7835 2.7835 1 3.75 1H6.25C6.66421 1 7 1.33579 7 1.75C7 2.16421 6.66421 2.5 6.25 2.5H3.75C3.61193 2.5 3.5 2.61193 3.5 2.75V13.25C3.5 13.3881 3.61193 13.5 3.75 13.5H6.25C6.66421 13.5 7 13.8358 7 14.25C7 14.6642 6.66421 15 6.25 15H3.75C2.7835 15 2 14.2165 2 13.25V2.75ZM12.4393 7.25H6.75002C6.33581 7.25 6.00002 7.58579 6.00002 8C6.00002 8.41422 6.33581 8.75 6.75002 8.75H12.4393L10.4697 10.7197C10.1768 11.0126 10.1768 11.4874 10.4697 11.7803C10.7626 12.0732 11.2374 12.0732 11.5303 11.7803L14.7803 8.53033C15.0732 8.23744 15.0732 7.76256 14.7803 7.46967L11.5303 4.21967C11.2374 3.92678 10.7626 3.92678 10.4697 4.21967C10.1768 4.51256 10.1768 4.98744 10.4697 5.28033L12.4393 7.25Z"
              fill="white"
            />
          </svg>
          Sign Out
        </Button>
      </div>
      <div class="w-60 mt-2 ml-2">
        <Input type="text" size="md" placeholder="🐣 Enter a Name for Note" bind:value={noteName}></Input>
      </div>
      <div>
        <Button type="submit" color="green" on:click={saveNote}  class="ml-2 mt-2">
          save
        </Button>
      </div>
    </div>
  {/if}

  {#if SignedIn}
    {#if !emptyHistory}
      <hr class="h-px my-8 bg-gray-200 border-0 dark:bg-gray-700 mb-2">
      <div class="text-center">
        <Heading tag="h3" class="font-semibold mb-2 underline">History <Badge class="text-2xl font-semibold ml-2">Broken</Badge></Heading>
      </div>
      <div class="flex flex-row justify-center">
        <div class="flex flex-col md:flex-row">
          {#each historyData as item}
              <Button class="m-2" color="light" on:click="{() => handleHistoryClick(item.fileName)}">{item.name}</Button>
          {/each}
        </div>
      </div>
    {/if}
  {/if}
</div>

<style>
  :global(.bytemd) {
    height: calc(100vh - 200px);
  }

  /* disable Github Permalink */
  :global(.bytemd-toolbar-right [bytemd-tippy-path="5"]) {
    display: none;
  }
</style>
