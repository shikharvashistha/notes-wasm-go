import { json } from "@sveltejs/kit";
import { clientInfo } from "../../../utils";

export const POST = async ({ request }: any) => {
  const res = await request.text();
  const body = JSON.parse(res);
  const GITHUB_OAUTH_URL = "https://github.com/login/oauth/access_token";
  const GITHUB_OAUTH_URL_PARAMS = `?client_id=${clientInfo.clientID}&client_secret=${clientInfo.clientSecret}&code=${body.Auth}`;

  const response = await fetch(GITHUB_OAUTH_URL + GITHUB_OAUTH_URL_PARAMS, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
  });
  return json(await response.json());
};
