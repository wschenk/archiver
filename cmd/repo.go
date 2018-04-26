package main

import (
	"fmt"
	"github.com/wschenk/archiver"
	"github.com/wschenk/archiver/ipfs"
	"github.com/wschenk/archiver/repository"
)

var test_hash = "QmcskP4cXKKnRm6GEhoTw52e8KzXBRJrvb6WLnVKuUx7zv"

func main() {
	var store archiver.ArchiveStore
	store = ipfs.CreateService()

	var repo archiver.Repository
	var err error
	repo, err = repository.CreateFileRepoFromStore(store, test_hash)

	if err != nil {
		panic(err)
	}

	contents, err := repo.Get("content/posts/my-first-post.md")

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", contents)

	if repo.Dirty() {
		fmt.Println("Repo dirty")
	} else {
		fmt.Println("Repo clean")
	}

	fmt.Println("Adding a file")
	err = repo.Put("content/posts/second-post.md", []byte("this is my file"))

	if err != nil {
		panic(err)
	}

	if repo.Dirty() {
		fmt.Println("Repo dirty")
	} else {
		fmt.Println("Repo clean")
	}

	fmt.Println("Saving repo")
	key, err := store.Put(repo)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Saved at %s\n", key)
	// repo.Clean()
}
