package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"syscall/js"

	"github.com/go-git/go-billy/v5/memfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	storagefs "github.com/go-git/go-git/v5/storage/filesystem"
)

var Filesystem = memfs.New() // create a new in-memory filesystem
/*
	* A note on the filesystem
	* This piece of Go code creates an in-memory filesystem using the memfs package from the go-git library. In-memory filesystems are useful when working with WebAssembly (WASM) because browsers typically do not allow direct access to the device storage, which can cause problems when trying to perform read/write operations. By using an in-memory filesystem, we can create a virtual filesystem in memory and perform I/O operations on it with its abstractions.
	* The code creates a new instance of an in-memory filesystem called "Filesystem" using the memfs.New() function. This filesystem will be used throughout the program to perform I/O operations. The rest of the code is just a skeleton and needs to be written to perform specific tasks.
	*
	* Most of the code is yet to be written and is just a skeleton 
*/


type GitRepo struct {
	storage *storagefs.Storage
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
			fmt.Println("Recovered from panic: ", r)
		}
	}()
	
	PATH := "repo"
	worktreeFs, _ := Filesystem.Chroot(PATH)
	dotGitFs, _ := Filesystem.Chroot(filepath.Join(PATH, ".git"))
	storage := storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	var re GitRepo

	go func() {
		repo, repoErr := gogit.Clone(storage, worktreeFs, &gogit.CloneOptions{
			URL: url,
		})

		check(repoErr)
		re = GitRepo{storage, repo}
	}()

	return re
}

func expose_git_clone(this js.Value, i []js.Value) interface{} {
	// A playground function to test the git_clone function

	url := i[0].String()
	git_clone(url)

	return nil
}

func git_push(repo GitRepo, creds credentials) {
	Repo := repo.gitRepo

	worktree, err := Repo.Worktree()
	check(err)

	_, err = worktree.Add(".")
	check(err)

	_, err = worktree.Commit("commit", &gogit.CommitOptions{})
	check(err)

	auth := &http.BasicAuth{
		Username: creds.username,
		Password: creds.password,
	}

	err = Repo.Push(&gogit.PushOptions{
		Auth: auth,
	})
	check(err)
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

func encrypt(key string) {
	// TODO
}

func decrypt(key string) {
	// TODO
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
