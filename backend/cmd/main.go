package main

import (
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"
	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Running"))
	r.Write(w)
}

func main() {
	//wasmBytes, _ := ioutil.ReadFile("/workspace/notes-wasm-go/backend/wasm/main.wasm")
	wasmBytes, err := os.ReadFile("/workspace/notes-wasm-go/backend/wasm/main.wasm")
	if err != nil {
		logger.Warn("Error reading wasm file")
	}
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		logger.Warn("Error creating module")
	}
	importObject := wasmer.NewImportObject()
	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		logger.Warn("Error creating instance")
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		encrypt, _ := instance.Exports.GetFunction("EncryptNotes")
		notesBytes := make([]byte, 1000)
		notesBody, err := r.Body.Read(notesBytes)
		if err != nil {
			logger.Warn("Error reading notes")
		}
		notes := string(notesBytes[:notesBody])
		key, err := encrypt(notes)
		if err != nil {
			logger.Warn("Error encrypting notes")
		}
		bytesKey := key.([]byte)

		w.Write(bytesKey)
	})

	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		decrypt, _ := instance.Exports.GetFunction("DecryptNotes")
		notesBytes := make([]byte, 1000)
		notesBody, err := r.Body.Read(notesBytes)
		if err != nil {
			logger.Warn("Error reading notes")
		}
		notesString := string(notesBytes[:notesBody])
		key := string(r.Header.Get("key"))
		notes, err := decrypt(notesString, key)
		if err != nil {
			logger.Warn("Error decrypting notes")
		}
		byteNotes := notes.([]byte)
		w.Write(byteNotes)
	})

	http.ListenAndServe(":8080", nil)
	logger.Info("Hello, world!")
}
