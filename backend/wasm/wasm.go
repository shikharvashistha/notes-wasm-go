package main

import (
	// "crypto/md5"
	// "encoding/hex"
	// "errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall/js"
	"time"

	// "github.com/go-git/go-git/v5"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"

	// "github.com/go-git/go-git/v5/plumbing/object"

	// "github.com/go-git/go-git/storage/memory"
	gogit "github.com/go-git/go-git/v5"
	storagefs "github.com/go-git/go-git/v5/storage/filesystem"
)

type Entry struct {
	Host      string
	Path      string
	URL       string
	GogitRepo *gogit.Repository
}

var Filesystem = memfs.New()
var AllRepositories = make(map[string]*Entry, 0)

func GetRepositoryList(this js.Value, i []js.Value) interface{} {
	retRepos := make([]interface{}, len(AllRepositories))

	repoIndex := 0
	for path, entry := range AllRepositories {
		// cfg, err := entry.GogitRepo.Config()
		// if err != nil {
		// 	return nil, nil
		// }
		repo := make(map[string]interface{}, 0)
		repo["path"] = path
		repo["host"] = entry.Host
		repo["path"] = entry.Path
		repo["url"] = entry.URL
		// repo["author"] = cfg.Author.Name
		// repo["author-email"] = cfg.Author.Email
		retRepos[repoIndex] = repo
		repoIndex += 1
	}

	return retRepos
}

func gitClone(this js.Value, i []js.Value) interface{} {
	url := i[0].String()
	path := i[1].String()

	worktreeFs, err := Filesystem.Chroot(path)
	if err != nil {
		return nil
	}

	dotGitFs, err := Filesystem.Chroot(filepath.Join(path, ".git"))
	if err != nil {
		return nil
	}

	storage := storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())

	go func() {
		repo, err := gogit.Clone(storage, worktreeFs, &gogit.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		repo.Log(&gogit.LogOptions{
			Order: gogit.LogOrderCommitterTime,
		})

		if err != nil {
			// if true {
			println("gogit.Clone() failed: ", err.Error())
			fmt.Println(err.Error())
		} else {
			fmt.Println("::: Cloned repository successfully.")
		}

		ref, _ := repo.Head()

		since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
		until := time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC)
		cIter, err := repo.Log(&gogit.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
			// ... just iterates over the commits, printing it
		err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)

		return nil
		})
		fmt.Println(cIter);
		// list files in repo
	}()


	return nil
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
func git_log(this js.Value, i []js.Value) interface{} {
	worktreeFs := Filesystem
	dir := i[0].String()

	fmt.Println("Filesystem contents:", worktreeFs.Root())
	repofiles, err := worktreeFs.ReadDir("/")
	if err != nil {
		fmt.Println("Error reading repo files:", err.Error())
	}

	// print out repofile names
	for _, file := range repofiles {
		fmt.Println("File:", file.Name())
	}

	// open dir
	dirFs, err := worktreeFs.Chroot(dir)
	if err != nil {
		fmt.Println("Error opening dir:", err.Error())
	}

	// list files in dir
	dirfiles, err := dirFs.ReadDir("/")
	if err != nil {
		fmt.Println("Error reading dir files:", err.Error())
	}

	// print out dirfile names
	for _, file := range dirfiles {
		fmt.Println("File:", file.Name())
	}
	return nil
}

func registerCallbacks() {
	println("Registering callbacks ...")

	// println(":\tencryptNotes()")
	// js.Global().Set("encryptNotes", js.FuncOf(encryptNotes))

	// println(":\tdecryptNotes()")
	// js.Global().Set("decryptNotes", js.FuncOf(decryptNotes))

	println(":\tgitClone()")
	js.Global().Set("gitClone", js.FuncOf(gitClone))
	print(":\t GetRepositoryList")
	js.Global().Set("lsrepo", js.FuncOf(GetRepositoryList))
	js.Global().Set("log", js.FuncOf(git_log))
	js.Global().Set("saveFile", js.FuncOf(saveFile))
	// js.Global().Set("gitCloneA", js.FuncOf(gitCloneA))
	// js.Global().Set("ro", js.FuncOf(ro))
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	registerCallbacks()
	fmt.Println("::: Cloned repository successfully.")
	<-c
}
