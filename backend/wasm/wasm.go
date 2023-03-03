package main

import (
	// "context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	// "fmt"
	// "os"
	"syscall/js"

	// "github.com/google/go-github/v50/github"
	// "github.com/google/uuid"

	logger "github.com/sirupsen/logrus"
)

//GOOS=js GOARCH=wasm go build -o main.wasm

func encryptNotes(this js.Value, i []js.Value) interface{} {
	
	key := make([]byte, 32)
	_, err := rand.Read(key)
	
	if err != nil {
		logger.Warn("Error generating key")
	}
	
	keyString := hex.EncodeToString(key)

	block, err := aes.NewCipher(key)
	// if err != nil {
	// 	logger.Warn("Error creating cipher block")
	// 	return "", err
	// }

	notes := i[0].String()
	
	stream := cipher.NewCTR(block, key[:block.BlockSize()])
	encrypted := make([]byte, len(notes))
	stream.XORKeyStream(encrypted, []byte(notes))


	// file, err := os.Create("notes.txt")
	// if err != nil {
	// 	logger.Warn("Error creating file")
	// 	return "", err
	// }
	// defer file.Close()
	// _, err = file.Write(encrypted)
	// if err != nil {
	// 	logger.Warn("Error writing to file")
	// 	return "", err
	// }
	// github := github.NewClient(nil)
	// _, _, err = github.Repositories.CreateFile(context.Background(), "shikharvashistha", "notes-wasm-go", "notes/"+uuid.New().String()+"./notes.txt", nil)
	// if err != nil {
	// 	logger.Warn("Error uploading file to github")
	// 	return "", err
	// }
	hexEncypted := hex.EncodeToString(encrypted)
	hexKey := keyString

	// TODO: remove this in production
	println("encypted(hex) : " + js.ValueOf(hexEncypted).String())
	println("keystring(hex): " + js.ValueOf(hexKey).String())

	return hex.EncodeToString(encrypted)
}

func decryptNotes(this js.Value, i []js.Value) interface{} {
	notes := i[0].String()
	key := i[1].String()

	keyGen, err := hex.DecodeString(key)
	if err != nil {
		logger.Warn("Error decoding key")
	}
	ciphertext, err := hex.DecodeString(notes)
	if err != nil {
		logger.Warn("Error decoding notes")
	}
	block, err := aes.NewCipher(keyGen)
	if err != nil {
		logger.Warn("Error creating cipher block")
	}

	stream := cipher.NewCFBDecrypter(block, keyGen[:block.BlockSize()])
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func registerCallbacks() {
	println("Registering callbacks ...")
	
	println(":\tencryptNotes()")
	js.Global().Set("encryptNotes", js.FuncOf(encryptNotes))
	
	println(":\tdecryptNotes()")
	js.Global().Set("decryptNotes", js.FuncOf(decryptNotes))
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	registerCallbacks()
	<-c
}
