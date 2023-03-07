package git
/*
	Made by @SaicharanKandukuri

	A sample wasm module to clone a git repository to an in-memory filesystem.
	Uses go-git and go-billy.
*/
import (
	"fmt"
	"path/filepath"
	"syscall/js"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	storagefs "github.com/go-git/go-git/v5/storage/filesystem"
)

// create a new in-memory filesystem
var Filesystem = memfs.New()

func Sample(this js.Value, i []js.Value) interface{} {
	url 	:= i[0].String()
	path 	:= i[1].String()

	worktreeFs, _ 	:= Filesystem.Chroot(path)
	dotGitFs, _ 	:= Filesystem.Chroot(filepath.Join(path, ".git"))
	storage 		:= storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	
	go func() {
		repo, repoErr := gogit.Clone(storage, worktreeFs, &gogit.CloneOptions{
			URL: url,
		})

		if repoErr != nil {
			fmt.Println("Error: ", repoErr)
		} else {
			fmt.Println("Repository cloned to ", "?fs?/"+path)
		}

		// git log
		ref, _ := repo.Head()
		since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
		until := time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC)
		cIter, _ := repo.Log(&gogit.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
		
		// ... just iterates over the commits, printing it
		_ = cIter.ForEach(func(c *object.Commit) error {
			fmt.Println(c)
			return nil
		})
		fmt.Println(cIter);
	}()

	return nil
}

func lsdir(this js.Value, i []js.Value) interface{} {
	path := i[0].String()
	fs   := Filesystem

	dir, err := fs.Chroot(path)
	check(err);

	files, err := dir.ReadDir("/")
	check(err);

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WasmLog(msg string) {
	fmt.Println("wasm: ", msg)
}

func registerCallbacks() {
	js.Global().Set("gitSample", js.FuncOf(Sample))
	js.Global().Set("lsdir", js.FuncOf(lsdir))
	WasmLog("Callbacks registered.")
}

func main() {
	c := make(chan struct{}, 0)
	// Main code here
	registerCallbacks()
	WasmLog("Done & started..")
	<-c
}
