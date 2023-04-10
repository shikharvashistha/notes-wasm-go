const clientInfo = {
    clientID: import.meta.env.VITE_APP_CLIENT_ID.toString(),
    clientSecret: import.meta.env.VITE_APP_CLIENT_SECRET.toString(),
}

const clientPub = {
    clientID: import.meta.env.VITE_APP_CLIENT_ID.toString(),
}

export { clientInfo, clientPub }
