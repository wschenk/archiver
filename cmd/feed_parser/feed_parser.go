package main

import (
	"fmt"

	"github.com/wschenk/archiver"
	"github.com/wschenk/archiver/repository"
	"github.com/wschenk/archiver/web"
)

func main() {
	testFeed("/tmp/feeds/willschenk.com", "http://willschenk.com/feed.xml")
	testFeed("/tmp/feeds/replyall", "http://feeds.gimletmedia.com/hearreplyall")
}

func testFeed(repoPath, feedUrl string) {
	var f archiver.Feed

	repo, err := repository.CreateFileRepository(repoPath)

	fetcher := web.CreateWebClient()

	if err != nil {
		panic(err)
	}

	f = feed.NewRssFeed(repo, feedUrl)

	refreshed, err := f.Refresh(fetcher)

	if err != nil {
		panic(err)
	}

	if refreshed {
		fmt.Println("New items in the feed")
	}

	fmt.Printf("Title   %s\n", f.Title())
	fmt.Printf("Type    %s\n", f.Type())
	fmt.Printf("Author  %s\n", f.Author())
	fmt.Printf("Updated %s\n", f.Updated())
	fmt.Printf("Link    %s\n", f.Link())
	fmt.Printf("Desc    %s\n", f.Description())
	fmt.Printf("Lang    %s\n", f.Language())

	fmt.Printf("There are %d items in the feed\n", len(f.Items()))
}
