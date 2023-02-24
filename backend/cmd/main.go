package main

import (
	"io"
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"
	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
	r.Write(w)
}

func main() {
	wasmBytes, _ := os.ReadFile("/workspace/notes-wasm-go/backend/wasm/main.wasm")
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	module, _ := wasmer.NewModule(store, wasmBytes)
	importObject := wasmer.NewImportObject()
	instance, _ := wasmer.NewInstance(module, importObject)

	http.HandleFunc("/", handler)
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		encrypt, _ := instance.Exports.GetFunction("EncryptNotes")
		notesBody, _ := io.ReadAll(r.Body)
		key, err := encrypt(notesBody)
		if err != nil {
			logger.Warn("Error encrypting notes")
		}
		bytesKey := key.([]byte)

		w.Write(bytesKey)
	})

	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		decrypt, _ := instance.Exports.GetFunction("DecryptNotes")
		notesBody, _ := io.ReadAll(r.Body)
		notes, err := decrypt(notesBody)
		if err != nil {
			logger.Warn("Error decrypting notes")
		}
		byteNotes := notes.([]byte)
		w.Write(byteNotes)
	})

	http.ListenAndServe(":8080", nil)
	logger.Info("Hello, world!")
}
