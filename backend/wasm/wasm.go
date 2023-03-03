package webassembly

import (
	// "context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"os"

	// "github.com/google/go-github/v50/github"
	// "github.com/google/uuid"

	logger "github.com/sirupsen/logrus"
)

//GOOS=js GOARCH=wasm go build -o main.wasm

func EncryptNotes(notes string) (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		logger.Warn("Error generating key")
		return "", err
	}
	keyString := hex.EncodeToString(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Warn("Error creating cipher block")
		return "", err
	}
	stream := cipher.NewCTR(block, key[:block.BlockSize()])
	encrypted := make([]byte, len(notes))
	stream.XORKeyStream(encrypted, []byte(notes))
	file, err := os.Create("notes.txt")
	if err != nil {
		logger.Warn("Error creating file")
		return "", err
	}
	defer file.Close()
	_, err = file.Write(encrypted)
	if err != nil {
		logger.Warn("Error writing to file")
		return "", err
	}
	// github := github.NewClient(nil)
	// _, _, err = github.Repositories.CreateFile(context.Background(), "shikharvashistha", "notes-wasm-go", "notes/"+uuid.New().String()+"./notes.txt", nil)
	// if err != nil {
	// 	logger.Warn("Error uploading file to github")
	// 	return "", err
	// }

	return keyString, nil
}

func DecryptNotes(notes string, key string) (string, error) {
	keyGen, err := hex.DecodeString(key)
	if err != nil {
		logger.Warn("Error decoding key")
		return "", err
	}
	ciphertext, err := base64.URLEncoding.DecodeString(notes)
	if err != nil {
		logger.Warn("Error decoding notes")
		return "", err
	}
	block, err := aes.NewCipher(keyGen)
	if err != nil {
		logger.Warn("Error creating cipher block")
		return "", err
	}

	stream := cipher.NewCFBDecrypter(block, keyGen[:block.BlockSize()])

	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
