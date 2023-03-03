package main

import (
	// "context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	// "fmt"
	"os"
	"syscall/js"

	"github.com/go-git/go-git/v5"
	// "github.com/google/uuid"

	logger "github.com/sirupsen/logrus"
)

//GOOS=js GOARCH=wasm go build -o main.wasm

func gitClone(this js.Value, i []js.Value) interface{} {
	url := i[0].String()
	println("Cloning " + url)

	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git",
		Progress: os.Stdout,
	})
	
	if err != nil {
		logger.Warn("Error cloning repository")
	}

	return nil
}

func encryptNotes(this js.Value, i []js.Value) interface{} {
	/*
		1. Generate a random key
			(32 bytes)
		2. encode the key to hex -> keyString
		3. encrypt the notes using the key
		4. encode the encrypted notes to hex -> hexEncrypted
	*/
	
	key := make([]byte, 32)
	_, err := rand.Read(key)
	
	if err != nil {
		logger.Warn("Error generating key")
	}
	
	keyString := hex.EncodeToString(key)

	block, err := aes.NewCipher(key)

	notes := i[0].String()
	
	stream := cipher.NewCTR(block, key[:block.BlockSize()])
	encrypted := make([]byte, len(notes))
	stream.XORKeyStream(encrypted, []byte(notes))

	hexEncypted := hex.EncodeToString(encrypted)
	hexKey := keyString

	// TODO: remove this in production
	println("encypted(hex) : " + js.ValueOf(hexEncypted).String())
	println("keystring(hex): " + js.ValueOf(hexKey).String())

	return hex.EncodeToString(encrypted)
}

func decryptNotes(this js.Value, i []js.Value) interface{} {
	/*
		1. Decode the encrypted notes and key from hex
		2. Use the key to create a new AES cipher block
		3. Create a new CTR stream using the cipher block and the nonce
		4. Decrypt the encrypted notes using the CTR stream
	*/

	hexEncrypted := i[0].String()
	hexKey := i[1].String()

	encrypted, err := hex.DecodeString(hexEncrypted)
	if err != nil {
		logger.Warn("Error decoding encrypted notes")
		return nil
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		logger.Warn("Error decoding key")
		return nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Warn("Error creating cipher block")
		return nil
	}

	stream := cipher.NewCTR(block, key[:block.BlockSize()])

	decrypted := make([]byte, len(encrypted))
	stream.XORKeyStream(decrypted, encrypted)

	return string(decrypted)
}

func registerCallbacks() {
	println("Registering callbacks ...")
	
	println(":\tencryptNotes()")
	js.Global().Set("encryptNotes", js.FuncOf(encryptNotes))
	
	println(":\tdecryptNotes()")
	js.Global().Set("decryptNotes", js.FuncOf(decryptNotes))

	println(":\tgitClone()")
	js.Global().Set("gitClone", js.FuncOf(gitClone))
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	registerCallbacks()
	<-c
}
