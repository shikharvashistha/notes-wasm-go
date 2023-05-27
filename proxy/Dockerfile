FROM golang:latest

WORKDIR /proxy

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8081

RUN go mod tidy

CMD ["go", "run", "/proxy/main.go"]
