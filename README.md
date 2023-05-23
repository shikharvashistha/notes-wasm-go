# Notes WASM GO
An experiment to see how far I can push WASM with GO to use it as a note taking app integrated with `git` and encryption logic in wasm to store notes in a secure way to a GitHub repo.

## Runing locally
Running this locally without crazy scripts requires 3 terminals
- Terminal 1: for building wasm file
- Terminal 2: for running the frontend server
- Terminal 4: for running the backend proxy server ( to avoid CORS issues )

### Terminal 1
> requires [go >1.16](https://go.dev/doc/install)

Backend has a script to build the wasm file and copy it to the frontend folder

```bash
cd backend
bash build.sh
```

### Terminal 2
> requires [node >14](https://nodejs.org/en/download/)

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

