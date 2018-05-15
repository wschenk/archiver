package emitter

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wschenk/archiver"
)

type GithubAccount struct {
	repo          archiver.Repository
	cache         archiver.Cache
	username      string
	subfeed       string
	pendingEvents []string
}

func CreateGithubStarsEmitter(repo archiver.Repository, username string) *GithubAccount {
	return &GithubAccount{
		repo:          repo,
		cache:         archiver.CreateRepoCache(repo),
		username:      username,
		subfeed:       "github_stars",
		pendingEvents: make([]string, 10),
	}
}

func (github *GithubAccount) Info() archiver.EmitterInfo {
	return archiver.EmitterInfo{
		Type:   github.subfeed,
		Author: github.username,
	}
}

func (github *GithubAccount) Refresh(fetcher archiver.Fetcher) (newStars bool, err error) {
	url := "https://github.com/" + github.username + "?tab=stars&page=1"

	data, err := github.cache.Get("stars.html", func() ([]byte, error) {
		log.Printf("Loading url %s\n", url)
		return fetcher.Fetch(url)
	}, time.Second*60)

	if err != nil {
		panic(err)
	}

	log.Printf("Data is %d big\n", len(data))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	foundExistingStar := false
	nextPage := true

	for nextPage {
		doc.Find(".col-9 .d-block").Each(func(i int, s *goquery.Selection) {
			repoName, exists := s.Find("h3 a").First().Attr("href")
			if exists {
				repoName := repoName[1:len(repoName)]
				description := s.Find("p").First().Text()
				// log.Println(repoName)
				// fmt.Println(description)

				if github.repo.HasKey(repoName) {
					foundExistingStar = true
				} else {
					newStars = true
					github.repo.Put(repoName, []byte(description))
					github.AddEvent(repoName)
				}
			}
		})

		nextPage = false

		nextPageLink := doc.Find("a.next_page").First()

		if nextPageLink != nil {
			if !foundExistingStar {
				url, exists := nextPageLink.Attr("href")
				if exists {
					nextPage = true
					url = "https://github.com" + url
					log.Printf("Loading %s\n", url)

					data, err = fetcher.Fetch(url)

					if err != nil {
						return newStars, err
					}

					doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data))
					if err != nil {
						return newStars, err
					}
				}
			}
		}
	}

	github.FlushEvents()

	return newStars, nil
}

func (github *GithubAccount) AddEvent(key string) {
	if github.repo.HasKey("index") {
		log.Printf("Appending key %s\n", key)
		indexBytes, _ := github.repo.Get("index")
		latestIndex, _ := strconv.Atoi(string(indexBytes))
		nextIndex := latestIndex + 1
		nextIndexString := strconv.Itoa(nextIndex)
		github.repo.Put(nextIndexString, []byte(key))
		github.repo.Put("index", []byte(nextIndexString))
	} else {
		log.Printf("Queueing event %s\n", key)
		github.pendingEvents = append(github.pendingEvents, key)
	}
}

func (github *GithubAccount) FlushEvents() {
	if github.repo.HasKey("index") {
		return
	}

	log.Println("Flushing events")
	count := 1

	for i := len(github.pendingEvents) - 1; i >= 0; i-- {
		nextIndexString := strconv.Itoa(count)

		if github.pendingEvents[i] != "" {
			github.repo.Put(nextIndexString, []byte(github.pendingEvents[i]))
			count += 1
		}
	}

	indexString := strconv.Itoa(count)
	github.repo.Put("index", []byte(indexString))
}
