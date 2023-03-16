package main

import (
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
	* the Filesystem is a in-memory filesystem
	* usually browsers sanbox features wont allow direct access to device storage which cause WASM to fail when it comes read/write operations when loaded in page
	* memfs comes handy in this case as it allows us to create a virtual filesystem in memory and perfrom I/O operations on it ( but with its abstractions )
	* 
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

	const data = i[0].String()
	return nil
}

func registerCallbacks() {
	js.Global().Set("AddNew", js.FuncOf(E_AddNew))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}


func main() {
	c := make(chan struct{}, 0)
	fmt.Println("WASM Go Initialized")
	registerCallbacks()
	<-c
}
