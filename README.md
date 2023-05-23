# Notes WASM GO
An experiment project ( part of an InternShip ) to see how far we can push WASM with GO to use it as a note taking app integrated with `git` and encryption logic in wasm to store notes in a secure way to a GitHub repo.

## Runing locally
Running this locally without crazy scripts requires 3 terminals
- Terminal 1: for building wasm file
- Terminal 2: for running the frontend server
- Terminal 3: for running the backend proxy server ( to avoid CORS issues )

### Terminal 1
> requires [go >1.16](https://go.dev/doc/install)

Backend has a script to build the wasm file and copy it to the frontend folder

```bash
cd backend
bash build.sh
```

### Terminal 2
> requires [node >14](https://nodejs.org/en/download/)

#### Add configuration
Create a `.env` file in the frontend folder and add the following
- `VITE_APP_CLIENT_ID` : GitHub OAuth App Client ID
- `VITE_APP_CLIENT_SECRET` : GitHub OAuth App Client Secret
- `VITE_APP_ENCRYPT_SECRET`: Encryption secret ( used to encrypt notes before sending them to GitHub )

in `frontend/src/repo.json` add json object with the following
| key | value |
|-----|-------|
| `Repo`| GitHub repo name in format: `USER/REPO` |
| `Branch`| GitHub repo branch name |
| `owner` | GitHub repo owner name |

now the final file should look like this
```json
{
  "Repo": "USER/REPO",
  "Branch": "main",
  "owner": "USER"
}
```

#### Run the server
```bash
cd frontend
npm i -g pnpm
pnpm i
pnpm run dev
```

This will start up the frontend server on a specific port (usually 5173). navigate to your address [localhost:5173](http://localhost:5173)

### Terminal 3
> requires [go >1.16](https://go.dev/doc/install)

```bash
cd proxy
go run main.go
```

