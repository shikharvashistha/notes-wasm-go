FROM node:latest

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

RUN mkdir -p ./{backend, frontend}
COPY backend    ./backend
COPY frontend   ./frontend

# Setup wasm build
# install go
RUN apt-get update && apt-get install -y wget git
RUN wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# [WASM] compile and copy wasm
RUN cd ./backend && bash build.sh


# [WEB] install dependencies
RUN npm install -g pnpm
RUN cd ./frontend && pnpm install

EXPOSE 5173
COPY ./frontend/package.json ./frontend/package.json
WORKDIR /app/frontend

CMD ["pnpm", "run", "dev", "--host"]
