package main

import (
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

// TODO Add a URL
const urlroot = "URLGOESHERE"

type Birdurl struct {
	Id   string `json:"id"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

var Bird_addresses = []string{
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/a/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/b/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/c/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/d/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/e/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/f/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/g/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/h/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/i/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/j/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/k/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/l/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/m/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/n/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/o/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/p/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/q/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/r/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/s/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/t/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/u/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/v/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/w/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/x/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/y/",
	"birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/z/",
}

func GetIndex() {
	// Regex for id extraction tested at https://regex-golang.appspot.com/assets/html/index.html
	rx_id := regexp.MustCompile("/(?P<id>[a-z]*)/index.aspx$")
	birds := []Birdurl{}
	for _, address := range Bird_addresses {

		doc, err := goquery.NewDocument(urlroot + address)
		if err != nil {
			return
		}
		doc.Find(".teaser-title a").Each(func(_ int, s *goquery.Selection) {
			// For each item found, get the band and title
			bird_name := s.Text()
			url, _ := s.Attr("href")

			id := rx_id.FindStringSubmatch(url)

			bird := Birdurl{id[1], urlroot + url, bird_name}
			birds = append(birds, bird)
		})
	}
	data, err := json.MarshalIndent(birds, "", "\t")
	if err != nil {
		return
	}
	ioutil.WriteFile("data/index.json", data, 0644)
}

func main() {
	GetIndex()
}
