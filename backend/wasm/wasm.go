package main

import (
	// "context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"os/exec"
	"time"

	"crypto/tls"
	"fmt"
	"net/http"

	// "os"
	"syscall/js"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	// "github.com/google/uuid"

	logger "github.com/sirupsen/logrus"
)

//GOOS=js GOARCH=wasm go build -o main.wasm

func gitClone(this js.Value, i []js.Value) interface{} {
	url := i[0].String()
	// println("Cloning " + url)
	// var err error
	// goChan := make(chan struct{}, 0)

	// go func() {
	// 	_, err = git.PlainClone("./", false, &git.CloneOptions{
	// 		URL: url,
	// 	})
	// 	goChan <- struct{}{}
	// 	}()

	// <-goChan

	// if err != nil {
	// 	logger.Warn("Error cloning repository")
	// }

	// return nil


	// Create a custom http(s) client with your config
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},

		// 15 second timeout
		Timeout: 15 * time.Second,

		// don't follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("https", githttp.NewClient(customClient))

	// Clone repository using the new client if the protocol is https://
	Info("git clone %s", url)

	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	if err != nil {
		logger.Warn("Error cloning repository")
	}

	return nil
}

func gitCloneA(this js.Value, i []js.Value) interface{} {
	
	url := i[0].String()
	executable, err := exec.LookPath("git")
	if err != nil {
		logger.Warn(err)
		logger.Warn("Error finding git executable")
		return nil
	}

	cmdExec := exec.Command("/usr/bin/git")
	fmt.Println(cmdExec.Run())

	cmd := exec.Command(executable, "clone", url)
	cmd.Path = executable
	println(executable)
	err = cmd.Run()
	if err != nil {
		logger.Warn(err)
		logger.Warn("Error cloning repo")
		return nil
	}

	return nil
}

func ro(this js.Value, i []js.Value) interface{} {
	c := make(chan struct{}, 0)

    // Register a function called "cloneRepo" in the global JavaScript namespace
    js.Global().Set("cloneRepo", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        // Get the URL of the Git repository to clone from the first argument
        url := args[0].String()

        // Execute the "git clone" command using the URL
        cmd := exec.Command("git", "clone", url)
        err := cmd.Run()
		if err != nil {
            fmt.Println(err.Error())
            return nil
        }

        fmt.Println("Cloned repository successfully.")
		return nil
    }))
	
    <-c
	fmt.Println("::: Cloned repository successfully.")
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
	js.Global().Set("gitCloneA", js.FuncOf(gitCloneA))
	js.Global().Set("ro", js.FuncOf(ro))
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	registerCallbacks()
	    // Register a function called "cloneRepo" in the global JavaScript namespace
		js.Global().Set("cloneRepo", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Get the URL of the Git repository to clone from the first argument
			url := args[0].String()
	
			// Execute the "git clone" command using the URL
			cmd := exec.Command("git", "clone", url)
			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
	
			fmt.Println("Cloned repository successfully.")
			return nil
		}))
	fmt.Println("::: Cloned repository successfully.")	
	<-c
}
