package emitter

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wschenk/archiver"
)

type MediumAccount struct {
	repo          archiver.Repository
	cache         archiver.Cache
	username      string
	subfeed       string
	pendingEvents []string
}

func CreateMediumArticlesEmitter(repo archiver.Repository, username string) *MediumAccount {
	return &MediumAccount{
		repo:     repo,
		cache:    archiver.CreateRepoCache(repo),
		username: username,
		subfeed:  "medium_articles",
	}
}

func (medium *MediumAccount) Info() archiver.EmitterInfo {
	return archiver.EmitterInfo{
		Type:   medium.subfeed,
		Author: medium.username,
	}
}

func (medium *MediumAccount) Refresh(fetcher archiver.Fetcher) (newStars bool, err error) {
	url := "https://medium.com/@" + medium.username + "/latest"

	data, err := medium.cache.Get("latest.html", func() ([]byte, error) {
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

	// html.css( ".postArticle").each do |div|
	//   title = div.css( "h3").text
	//   link = div.css( ".layoutSingleColumn a").first['href']
	//   link = link.gsub( /\?.*/, "")
	//   out << [title,link]
	// end

	foundExistingStar := false
	nextPage := true

	for nextPage {
		doc.Find("h1").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			log.Println("title: " + title)
			selection := s.ParentsFiltered("a")
			// log.Println(selection.Html())
			link, exists := selection.First().Attr("href")

			if exists {
				log.Println("title : " + title)
				log.Println("href  : " + link)
				// repoName := repoName[1:len(repoName)]
				// description := s.Find("p").First().Text()
				// // log.Println(repoName)
				// // fmt.Println(description)
				//
				// if medium.repo.HasKey(repoName) {
				// 	foundExistingStar = true
				// } else {
				// 	newStars = true
				// 	medium.repo.Put(repoName, []byte(description))
				// 	medium.AddEvent(repoName)
				// }
			}
		})

		nextPage = false

		nextPageLink := doc.Find("a.next_page").First()

		if nextPageLink != nil {
			if !foundExistingStar {
				url, exists := nextPageLink.Attr("href")
				if exists {
					nextPage = true
					url = "https://medium.com" + url
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

	medium.FlushEvents()

	return newStars, nil
}

func (medium *MediumAccount) AddEvent(key string) {
	if medium.repo.HasKey("index") {
		log.Printf("Appending key %s\n", key)
		indexBytes, _ := medium.repo.Get("index")
		latestIndex, _ := strconv.Atoi(string(indexBytes))
		nextIndex := latestIndex + 1
		nextIndexString := strconv.Itoa(nextIndex)
		medium.repo.Put(nextIndexString, []byte(key))
		medium.repo.Put("index", []byte(nextIndexString))
	} else {
		log.Printf("Queueing event %s\n", key)
		medium.pendingEvents = append(medium.pendingEvents, key)
	}
}

func (medium *MediumAccount) FlushEvents() {
	if medium.repo.HasKey("index") {
		return
	}

	log.Println("Flushing events")
	count := 1

	for i := len(medium.pendingEvents) - 1; i >= 0; i-- {
		nextIndexString := strconv.Itoa(count)

		if medium.pendingEvents[i] != "" {
			medium.repo.Put(nextIndexString, []byte(medium.pendingEvents[i]))
			count += 1
		}
	}

	indexString := strconv.Itoa(count)
	medium.repo.Put("index", []byte(indexString))
}
