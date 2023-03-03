#!/bin/bash

WASM_NAME="main.wasm"

if ! command -v go &> /dev/null
then
    echo "go could not be found"
    exit
fi

cd wasm || exit
go mod tidy
GOOS=js GOARCH=wasm go build -o ../$WASM_NAME
