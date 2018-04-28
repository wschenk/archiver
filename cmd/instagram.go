package main

import (
	"fmt"
	// "github.com/wschenk/archiver"
	"github.com/wschenk/archiver/emitter/instagram"
	"github.com/wschenk/archiver/repository"
	"os"
)

func main() {
	account := "wschenk"

	if len(os.Args) > 1 {
		account = os.Args[1]
	}

	fmt.Printf("Looking for account %s\n", account)

	repo, err := repository.CreateFileRepository("/tmp/instagram/" + account)

	if err != nil {
		panic(err)
	}

	feed := feed.CreateInstagramFeed(repo, account)

	newImages, err := feed.Refresh()

	if err != nil {
		panic(err)
	}

	if newImages {
		fmt.Println("New images found")
	} else {
		fmt.Println("Post something!")
	}
}
