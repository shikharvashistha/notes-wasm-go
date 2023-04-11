package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"syscall/js"
	"time"

	// "crypto/aes"
	// "crypto/cipher"
	// "crypto/rand"
	// "encoding/hex"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	storagefs "github.com/go-git/go-git/v5/storage/filesystem"
)

/*
	* A note on the filesystem
	* This piece of Go code creates an in-memory filesystem using the memfs package from the go-git library. In-memory filesystems are useful when working with WebAssembly (WASM) because browsers typically do not allow direct access to the device storage, which can cause problems when trying to perform read/write operations. By using an in-memory filesystem, we can create a virtual filesystem in memory and perform I/O operations on it with its abstractions.
	* The code creates a new instance of an in-memory filesystem called "Filesystem" using the memfs.New() function. This filesystem will be used throughout the program to perform I/O operations. The rest of the code is just a skeleton and needs to be written to perform specific tasks.
	*
	* ----- Work in progress -----
	*
	*
	* ISSUES and WORKAROUND for golang WASM
	*
	* 1) WASM does not support multithreading and most of gorooutines are panics
	*      - trigger methods which uses routines internally with goroutines ( i.e "go func() { //code }()" )
	*      - use time.Sleep() to wait for the routine to finish [ which is not a good solution but works for now ]
	*
	* 2) WASM does not support file system access
	*      - use in-memory filesystems ( i.e "github.com/go-git/go-billy/v5/memfs" )
	*      - use "github.com/go-git/go-git/v5/plumbing/cache" to cache objects
	*      - use "github.com/go-git/go-git/v5/storage/filesystem" to store objects
	*
	* 3) WASM memory is not system memory ( https://radu-matei.com/blog/practical-guide-to-wasm-memory/ )
	* 
*/
	

var Filesystem = memfs.New() // create a new in-memory filesystem
var PATH = "repo"

type GitRepo struct {
	storage *billy.Filesystem
	gitRepo *gogit.Repository
}

type credentials struct {
	username string
	password string
}

func git_clone(url string) GitRepo {

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("git_clone: Recovered from panic: ", r)
		}
	}()
	
	PATH := "repo"
	worktreeFs, _ := Filesystem.Chroot(PATH)
	dotGitFs, _ := Filesystem.Chroot(filepath.Join(PATH, ".git"))
	storage := storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	var re GitRepo

	go func() {
		_, repoErr := gogit.Clone(storage, worktreeFs, &gogit.CloneOptions{
			URL: url,
		})

		check(repoErr)
	}()
	
	time.Sleep(10 * time.Second)
	repo, err := gogit.Open(storage, worktreeFs)
	checkErr(err)
	re = GitRepo{ &worktreeFs, repo }
	
	return re
}

func expose_git_clone(this js.Value, i []js.Value) interface{} {
	// A playground function to test the git_clone function

	url := i[0].String()
	git_clone(url)

	return nil
}

func git_push(repo GitRepo, creds credentials) {
	Repo := *repo.gitRepo
	storage := *repo.storage
	
	_, e := storage.Create("test.txt")

	check(e)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("git_push: Recovered from panic: ", r)
		}
	}()
	
	fmt.Println("entered git_push")

	worktree, err := Repo.Worktree()
	check(err)
	fmt.Println("got worktree")

	_, err = worktree.Add(".")
	check(err)
	fmt.Println("added files")

	_, err = worktree.Commit("commit", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "",
			Email: "",
			When:  time.Now(),
		},
	})
	check(err)
	fmt.Println("committed")

	auth := &http.BasicAuth{
		Username: creds.username,
		Password: creds.password,
	}
	fmt.Println("got auth")

	go func() {
		err := Repo.Push(&gogit.PushOptions{
			Auth: auth,
		})
		check(err)
	}()
}

func saveFile(this js.Value, i []js.Value) interface{} {
	//save the file in the repository
	path := i[0].String()
	content := i[1].String()

	fmt.Println("Saving file:", path)
	fmt.Println("File contents:", content)

	dirFs, err := Filesystem.Chroot(path)
	if err != nil {
		fmt.Println("Error opening dir:", err.Error())
	}
	fileName := filepath.Base(path)
	file, err := dirFs.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
	}
	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Println("Error writing file:", err.Error())
	}
	return nil
}

func encrypt(note string) {
	// TODO get saved key from server
}

func decrypt(note string, key string) {
	
}

func init_notes_fs(filesystem *storagefs.Storage) {
	// TODO
}

func E_AddNew(this js.Value, i []js.Value) interface{} {
	/*
		* [TODO] trigger setup action that sets up the filesystem
		* [TODO] get the notes -> encrypt -> store in filesystem
		* [TODO] push to remote
	*/

	var data = i[0].String()
	todoMsg("AddNew called with data: " + data)
	return nil
}


// test driver function to simulate the AddNew function
func testPush(this js.Value, i []js.Value) interface{} {
	
	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TEST: Recovered from panic: ", r)
		}
	}()

	var url = ""
	var auth = ""
	
	if len(i) > 0 {
		url = i[0].String()
		auth = i[1].String()
	}

	fmt.Println(url + "|" + auth)
	touch("/.preserve")
	fmt.Println("touch done")
	var repo = git_clone(url)
	fmt.Println("clone done")

	fs, err := Filesystem.Chroot(PATH)
	checkErr(err)

	/*
	* 	ISSUE: test.txt never got created and no error ?
	*/
	file, err := fs.Create("test.txt")
	checkErr(err)
	file.Close()

	fmt.Println("create done")

	var creds = credentials{"saicharankandukuri1x1", auth}
	git_push(repo, creds)
	
	return true
}

func lsDir(url string) {
	filesystem := Filesystem

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic: ", r)
		}
	}()

	_, err := filesystem.Stat(url)
	check(err)

	files, err := filesystem.ReadDir(url)
	checkErr(err)

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func touch(path string) {
	filesystem := Filesystem

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic: ", r)
		}
	}()

	file, err := filesystem.Create(path)
	checkErr(err)
	file.Close()
}

func writetofile(buff string, path string) {
	filesystem := Filesystem

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic: ", r)
		}
	}()

	file, err := filesystem.Create(path)
	checkErr(err)
	file.Write([]byte(buff))
	file.Close()
}

func cat(path string) {
	// cat the file at path
	filesystem := Filesystem

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic: ", r)
		}
	}()

	file, err := filesystem.Open(path)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err.Error())
	}

}

func wasm_lsDir(this js.Value, i []js.Value) interface{} {
	// A playground function to test the lsDir function
	var url ="/"

	if len(i) > 0 {
		url = i[0].String()
	}

	lsDir(url)

	return true
}

func wasm_touch(this js.Value, i []js.Value) interface{} {
	// A playground function to test the touch function
	var path = ""

	if len(i) > 0 {
		path = i[0].String()
	} else {
		fmt.Println("usage func(\"<filename>\"): requires one argument")
		return false
	}

	touch(path)
	return true
}

func wasm_writetofile(this js.Value, i []js.Value) interface{} {
	// A playground function to test the writetofile function

	var buff, path = "", ""

	if len(i) > 0 {
		buff = i[0].String()
		path = i[1].String()
	} else {
		fmt.Println("usage func(\"<msg>\",\"<path>\"): requires two argumets")
		return false
	}

	writetofile(buff, path)

	return true
}

func wasm_cat(this js.Value, i []js.Value) interface{} {
	// A playground function to test the cat function
	var path = ""

	if len(i) > 0 {
		path = i[0].String()
	} else {
		fmt.Println("usage func(\"<filename>\"): requires one valid argument")
		return false
	}

	cat(path)

	return true
}

func registerCallbacks() {
	js.Global().Set("AddNew", js.FuncOf(E_AddNew))
	js.Global().Set("git_clone", js.FuncOf(expose_git_clone))
	js.Global().Set("test_push", js.FuncOf(testPush))

	// playground functions
	js.Global().Set("wasm_ls", js.FuncOf(wasm_lsDir))
	js.Global().Set("wasm_touch", js.FuncOf(wasm_touch))
	js.Global().Set("wasm_write", js.FuncOf(wasm_writetofile))
	js.Global().Set("wasm_cat", js.FuncOf(wasm_cat))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkErr(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
	}
}

func todoMsg(msg string) {
	fmt.Println("TODO:", msg)
}

func main() {
	touch("/.preserve") // create a file to preserve the filesystem
	c := make(chan struct{}, 0)
	fmt.Println("WASM Go Initialized")
	registerCallbacks()
	<-c
}
