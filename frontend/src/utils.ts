const clientInfo = {
    clientID: import.meta.env.VITE_APP_CLIENT_ID.toString(),
    clientSecret: import.meta.env.VITE_APP_CLIENT_SECRET.toString(),
}

const clientPub = {
    clientID: import.meta.env.VITE_APP_CLIENT_ID.toString(),
}

const spice = {
    encryptSecret: import.meta.env.VITE_APP_ENCRYPT_SECRET.toString(),
}

export { clientInfo, clientPub, spice }
