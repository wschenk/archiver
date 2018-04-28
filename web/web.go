package web

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

type WebClient struct {
	client *http.Client
}

func CreateWebClient() *WebClient {
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	return &WebClient{client}
}

func (client *WebClient) Fetch(url string) ([]byte, error) {
	resp, err := client.client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Printf("Status %s\n", resp.Status)

	fmt.Println("Cookies", resp.Cookies())

	return ioutil.ReadAll(resp.Body)
}
