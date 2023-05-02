import { json } from "@sveltejs/kit";

const USER_ENDPOINT = "https://api.github.com/user"

export const POST = (async ({ request }: any) => {
    const res = await request.text();
    const body = JSON.parse(res)
    const Auth = body.Auth

    await fetch(USER_ENDPOINT, {
        method: "GET",
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json",
            "Authorization": `Bearer ${Auth}`
        }
    }).then((response) => {
        return json(response.json())
    })
})
