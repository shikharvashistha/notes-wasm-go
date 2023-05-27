#!/bin/bash

WASM_NAME="main.wasm"
FRONTEND_PATH="../../frontend/src/lib"

if ! command -v go &> /dev/null
then
    echo "go could not be found"
    exit
fi

cd wasm || exit
go mod tidy
GOOS=js GOARCH=wasm go build -o ../$WASM_NAME
while [ $# -gt 0 ]; do
    case $1 in
        --frontend-path)
            FRONTEND_PATH=$2
            shift 2
            ;;
        *)
            echo "Argument: $1 ignored"
            ;;
    esac
done


cp -v ../$WASM_NAME $FRONTEND_PATH
