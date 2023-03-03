$WASM_NAME = "main.wasm"
 #
 #
if (-NOT (Get-Command go -ErrorAction SilentlyContinue)) {
    echo "go could not be found"
    exit
}
 #
 #
cd wasm
go mod tidy
$Env:GOOS = "js"
$Env:GOARCH = "wasm"
go build -o ../$WASM_NAME
 #
