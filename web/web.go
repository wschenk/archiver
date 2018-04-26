package web

import (
	"github.com/wschenk/archiver"
	"io/ioutil"
	"net/http"
)

func Fetch(repo archiver.Repository, repoPath string, url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		err = repo.Put(repoPath, data)
		if err != nil {
			return nil, err
		}
	}

	return data, err
}

// // URL
// func LoadUrl(string URL) string {

// }

// // URL -> URLS -> Page
// func GetPage(string URL) string {

// }

// {Site, User} -> Feed

// rss feed
// Twitter feed
// Instagram posts
// Youtube likes
// sound cloud likes
// spotify likes ?
// link feed -- posted, browsed
// saved linkes screen shot, archive
// saved videos

// youtube history
// browsing history
