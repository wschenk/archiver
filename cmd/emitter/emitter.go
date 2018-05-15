package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/wschenk/archiver"
	"github.com/wschenk/archiver/emitter"
	"github.com/wschenk/archiver/repository"
	"github.com/wschenk/archiver/web"
)

var debug, logger, errLogger *log.Logger

func main() {
	logger = log.New(os.Stdout, "[emitter] ", log.Ldate|log.Ltime|log.Lshortfile)
	errLogger = log.New(os.Stderr, "[emitter] ", log.Ldate|log.Ltime|log.Lshortfile)
	debug = log.New(os.Stdout, "[emitter] ", log.Ldate|log.Ltime|log.Lshortfile)

	debug.Printf("We have %d arguments\n", len(os.Args))
	if len(os.Args) != 3 {
		errLogger.Println("Usage:")
		errLogger.Printf("%s [provider] [account]\n", os.Args[0])
		os.Exit(1)
	}

	provider := os.Args[1]
	account := os.Args[2]

	debug.Printf("Looking for account %s on %s\n", account, provider)

	// repoPath := filepath.Join(os.TempDir(), account, provider)
	repoPath := filepath.Join("/tmp", provider, account)

	repo, err := repository.CreateFileRepository(repoPath)

	if err != nil {
		panic(err)
	}

	var feed archiver.Emitter

	switch provider {
	case "instagram":
		feed = emitter.CreateInstagramEmitter(repo, account)
	case "github_stars":
		feed = emitter.CreateGithubStarsEmitter(repo, account)
	case "medium_articles":
		feed = emitter.CreateMediumArticlesEmitter(repo, account)
	default:
		errLogger.Printf("Unknown prodiver %s\n", provider)
		os.Exit(1)
	}

	fetcher := web.CreateWebClient()

	if err != nil {
		panic(err)
	}

	newItems, err := feed.Refresh(fetcher)

	if err != nil {
		panic(err)
	}

	if newItems {
		debug.Println("New images found")
	} else {
		debug.Println("Post something!")
	}
}
