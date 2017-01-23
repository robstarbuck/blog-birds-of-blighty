package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"
)

type Birdurl struct {
	Id string `json:"id"`

	Url  string `json:"url"`
	Name string `json:"name"`
}

func DownloadSingle(location, filename, imageurl string) (err error) {

	err = os.MkdirAll(location, os.ModePerm)
	if err != nil {
		panic(err)
	}

	response, err := http.Get(imageurl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Open a file for writing
	file, err := os.Create(location + "/" + filename)
	if err != nil {
		panic(err)
	}
	// Use io.Copy to dump the response body to the file.
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	file.Close()

	return
}

func ReadFile() (birds []Birdurl) {

	file, err := ioutil.ReadFile("./data/index.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(file, &birds)

	return
}

// Get the src / urls of the images
func ImageSrcs(address string, width string) (srcs []string, captions []string, err error) {

	doc, err := goquery.NewDocument(address)

	if err != nil {
		return
	}

	doc.Find(".carousel-inner img").Each(func(i int, img *goquery.Selection) {

		src, exist := img.Attr("src")
		if exist {
			// Width modifies url which support image resize by parameter
			if width != "" {
				rx_removewidth := regexp.MustCompile("width=\\d*")
				src = rx_removewidth.ReplaceAllString(src, "width="+width)
			} else {
				u, err := url.Parse(src)
				if err != nil {
					panic(err)
				}
				src = u.Scheme + "://" + u.Host + u.Path
			}

			srcs = append(srcs, src)

			captions = append(captions, img.Parent().Find(".carousel-caption h4").Text())

		}
	})
	return
}

func DownloadImages(bird Birdurl) (err error) {

	srcs, captions, err := ImageSrcs(bird.Url, "2000")
	if err != nil {
		panic(err)
	}

	for i, src := range srcs {

		imagename := sanitize.Name(strconv.Itoa(i) + "_" + captions[i] + ".jpg")

		DownloadSingle("./images/"+bird.Id, imagename, src)
	}
	return
}

func main() {
	birds := ReadFile()

	for i, bird := range birds {
		DownloadImages(bird)
		if i > 1 {
			// TODO Include when testing
			// break
		}
	}
}
