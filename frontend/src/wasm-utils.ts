async function encryptHandler(note: string) {
    const res = await Promise.resolve(
        // @ts-ignore
        encrypt_text(note, spice.encryptSecret)
    )
    return res
}

async function decryptHandler(note: string) {
    const res = await Promise.resolve(
        // @ts-ignore
        decrypt_text(note, spice.encryptSecret)
    )
    return res
}

async function writeNoteToFile(file: string, note: string) {
    const res = await Promise.resolve(
        // @ts-ignore
        touchNcat(file, note)
    )

    return res
}

async function gitClone(url: string) {
    const res = await Promise.resolve(
        // @ts-ignore
        git_clone(url)
    )

    return res
}

async function gitPush(
    url: string,
    AccessTocken: string,
    userName: string,
    userEmail: string,
    file: string,
    commitMessage: string
) {

    const res = await Promise.resolve(
        // @ts-ignore
        git_push(
            url,
            AccessTocken,
            userName,
            userEmail,
            file,
            commitMessage
        )
    )

    return res

}

async function getAccessTocken(code: any) {
    const ENDPOINT = "http://localhost:5173/api/getAuthCode"
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
