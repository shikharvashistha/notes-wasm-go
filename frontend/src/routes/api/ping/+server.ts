import { json } from "@sveltejs/kit";

export async function GET() {
    // return pong
    return json(
        {
            pong: "alive"
        }
    )
}
