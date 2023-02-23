FROM golang:latest

WORKDIR /backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 8080

#RUN GOOS=js GOARCH=wasm go build /app/cmd/main.go
RUN go build -o main /backend/cmd/main.go

CMD ["./main"]