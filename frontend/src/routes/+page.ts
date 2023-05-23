import { Repo, Branch } from '../repo.json'
export const ssr = false;

export const load = async ({ fetch }) => {
    const rawUrl = "https://raw.githubusercontent.com/" + Repo + "/"+ Branch +"/"
    const res = await fetch(rawUrl + "history.json")
    const data = await res.json()
    return {
        history: data
    }
}
