/*
    GITHUB SECTION

    -- As WASM runs in Client side, tokens are stored in local storage.

    SignIn:
        - goto https://github.com/login/oauth/authorize... with client_id and scope
        - github redirects to CALLBACK_URL with code
        - onMount, if the code parameter is present, call /api/getAccessTocken with code
        - /api/getAccessTocken returns access_token
        - save access_token to local storage as GITHUB_ACCESS_TOKEN
    
    SignOut:
        - remove GITHUB_ACCESS_TOKEN from local storage
        - reload page
    
    Detect:
        - if GITHUB_ACCESS_TOKEN is present in local storage, show SignOut button
        - else show SignIn button
*/

import { clientPub } from "../utils";
import { SignIn } from "../stores";
import { Octokit } from "octokit";
import { Repo, owner } from "../repo.json";

function SignInGitHub() {
  window.location.assign(
    "https://github.com/login/oauth/authorize?client_id=" +
      clientPub.clientID +
      "&scope=repo"
  );
}

function SignOutGitHub() {
  // remove access token from local storage
  localStorage.removeItem("GITHUB_ACCESS_TOKEN");

  // remove code from url
  window.history.replaceState({}, document.title, "/");

  SignIn.set(false);
}

// getAccessTocken call
async function getAccessTocken(code: string) {
  const res = await fetch("/api/getAccessTocken", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      Auth: code,
    }),
  });
  const data = await res.json();
  return data;
}

async function getFileContents(fileName: String) {
  const oktokit = new Octokit({
    auth: localStorage.getItem("GITHUB_ACCESS_TOKEN"),
  });
  const res = await oktokit.request(
    "GET /repos/" + Repo + "/contents/" + fileName,
    {
      owner: Repo.split("/")[0],
      repo: Repo.split("/")[1],
      path: fileName,
      headers: {
        "X-GitHub-Api-Version": "2022-11-28",
      },
    }
  );
  const fileURL = res.data.download_url;
  const file = await fetch(fileURL);
  const fileContents = await file.text();
  return fileContents;
}

async function uploadFile(
  content: string,
  fileName: String,
  user: string,
  email: string
) {
  const oktokit = new Octokit({
    auth: localStorage.getItem("GITHUB_ACCESS_TOKEN"),
  });
  const res = await oktokit.request(
    "PUT /repos/" + Repo + "/contents/" + fileName,
    {
      owner: owner,
      repo: Repo.split("/")[1],
      path: fileName,
      message: "Add: " + fileName,
      committer: {
        name: user,
        email: email,
      },
      content: content,
      headers: {
        "X-GitHub-Api-Version": "2022-11-28",
      },
    }
  );

  return res.data.content.download_url;
}

const GH_Helper = {
  SignIn,
  SignInGitHub,
  SignOutGitHub,
  getAccessTocken,
  getFileContents,
  uploadFile,
};

export { GH_Helper };
