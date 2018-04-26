package feed

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/wschenk/archiver"
	"github.com/wschenk/archiver/web"
)

type RssFeed struct {
	repo archiver.Repository
	url  string
}

func NewRssFeed(repo archiver.Repository, url string) *RssFeed {
	return &RssFeed{repo: repo, url: url}
}

func (rss *RssFeed) Id() string {
	return rss.url
}

func (rss *RssFeed) Refresh() (bool, error) {
	data, err := web.Fetch(rss.repo, "feed.xml", rss.url)

	if err != nil {
		return false, err
	}

	fp := gofeed.NewParser()
	feed, _ := fp.ParseString(string(data))

	rss.putString("title", feed.Title)
	rss.putString("type", "rss")
	rss.putString("author", feed.Author.Name)
	rss.putString("email", feed.Author.Email)
	rss.putString("updated", feed.Updated)
	rss.putString("link", feed.Link)
	rss.putString("description", feed.Description)
	rss.putString("language", feed.Language)

	if feed.Image != nil {
		fmt.Println(feed.Image.URL)
		// web.Fetch(rss.repo, "image.png", feed.Image.URL)
		// web.Fetch
	}

	return false, nil
}

func (rss *RssFeed) Items() []archiver.FeedItem {
	s := []archiver.FeedItem{}
	return s
}

func (rss *RssFeed) Title() string {
	return rss.getString("title")
}

func (rss *RssFeed) Type() string {
	return rss.getString("type")
}

func (rss *RssFeed) Author() string {
	return rss.getString("author")
}

func (rss *RssFeed) Updated() string {
	return rss.getString("updated")
}

func (rss *RssFeed) Link() string {
	return rss.getString("link")
}

func (rss *RssFeed) Description() string {
	return rss.getString("description")
}

func (rss *RssFeed) Language() string {
	return rss.getString("language")
}

func (rss *RssFeed) getString(key string) string {
	data, _ := rss.repo.Get(key)
	return string(data)
}

func (rss *RssFeed) putString(key string, value string) {
	rss.repo.Put(key, []byte(value))
}
