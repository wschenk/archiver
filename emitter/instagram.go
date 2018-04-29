package emitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/wschenk/archiver"
	"log"
	"strings"
	"time"
)

type InstagramAccount struct {
	repo     archiver.Repository
	cache    archiver.Cache
	username string
}

type user struct {
	Id       string `json:"id"`
	Timeline struct {
		Edges []struct {
			Node node `json:"node"`
		} `json:"edges"`
		PageInfo pageInfo `json:"page_info"`
	} `json:"edge_owner_to_timeline_media"`
}

type node struct {
	ImageURL   string `json:"display_url"`
	Shortcode  string `json:"shortcode"`
	IsVideo    bool   `json:"is_video"`
	Date       int    `json:"taken_at_timestamp"`
	Dimensions struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
}

type pageInfo struct {
	EndCursor string `json:"end_cursor"`
	NextPage  bool   `json:"has_next_page"`
}

// const nextPageURLTemplate string = `https://www.instagram.com/graphql/query/?query_id=17888483320059182&variables={"id":"%s","first":12,"after":"%s"}`

const nextPageURLTemplate string = `https://www.instagram.com/graphql/query/?query_hash=42323d64886122307be10013ad2dcc44&variables={"id":"%s","first":12,"after":"%s"}`

func CreateInstagramEmitter(repo archiver.Repository, username string) *InstagramAccount {
	return &InstagramAccount{repo: repo,
		cache:    archiver.CreateRepoCache(repo),
		username: username}
}

func (insta *InstagramAccount) Info() archiver.EmitterInfo {
	return archiver.EmitterInfo{
		Type:   "instagram",
		Author: insta.username,
	}
}

func (insta *InstagramAccount) Refresh(fetcher archiver.Fetcher) (newImages bool, err error) {
	url := "https://instagram.com/" + insta.username

	data, err := insta.cache.Get("index.html", func() ([]byte, error) {
		log.Printf("Loading url %s\n", url)
		return fetcher.Fetch(url)
	}, time.Second*60)

	if err != nil {
		panic(err)
	}

	log.Printf("Data is %d big\n", len(data))

	var actualUserId string

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	foundPreviouslyLoadedImage := false

	// Find the first script block
	doc.Find("body > script").Each(func(i int, s *goquery.Selection) {
		scriptString := s.Text()

		if i == 0 {
			jsonData := scriptString[strings.Index(scriptString, "{") : len(scriptString)-1]

			// fmt.Println(jsonData)

			// jsonData := scriptString.slice(scriptString.Index(scriptString, "{"), len(scriptString)-1)
			data := struct {
				EntryData struct {
					ProfilePage []struct {
						Graphql struct {
							User user `json:"user"`
						} `json:"graphql"`
					} `json:"ProfilePage"`
				} `json:"entry_data"`
			}{}
			err := json.Unmarshal([]byte(jsonData), &data)

			if err != nil {
				panic(err)
			}

			log.Println(data)

			page := data.EntryData.ProfilePage[0]
			actualUserId = page.Graphql.User.Id

			loadEntries := func(u user) {
				for _, obj := range u.Timeline.Edges {
					// skip videos
					if obj.Node.IsVideo {
						continue
					}

					if !foundPreviouslyLoadedImage {
						loaded := insta.EnsureObject(fetcher, obj.Node)
						if loaded {
							newImages = true
						} else {
							foundPreviouslyLoadedImage = true
						}
					}
				}
			}

			loadEntries(page.Graphql.User)

			pageInfo := page.Graphql.User.Timeline.PageInfo
			// if !foundPreviouslyLoadedImage &&
			if pageInfo.NextPage {
				log.Println("Scraping for next page of images")
				nextPageURL := fmt.Sprintf(nextPageURLTemplate, actualUserId, pageInfo.EndCursor)

				log.Println("Loading", nextPageURL)

				data, err := fetcher.Fetch(nextPageURL)

				if err != nil {
					panic(err)
					return
				}
				log.Println(string(data))
			}
		}
	})
	return newImages, nil
}

func (i *InstagramAccount) EnsureObject(fetcher archiver.Fetcher, node node) (loaded bool) {
	log.Printf("Ensuring %s\n", node.Shortcode)
	log.Printf("Timeline %d\n", node.Date)

	path := fmt.Sprintf("%d/image.jpg", node.Date)
	loaded = false
	_, err := i.cache.Get(path, func() ([]byte, error) {
		log.Printf("Loading %s\n", node.ImageURL)
		loaded = true
		return fetcher.Fetch(node.ImageURL)
	}, time.Second*60)

	if err != nil {
		panic(err)
	}

	path = fmt.Sprintf("%d/shortcode", node.Date)
	i.repo.Put(path, []byte(node.Shortcode))

	return loaded
}
