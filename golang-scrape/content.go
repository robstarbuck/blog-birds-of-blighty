package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Birdscontainer struct {
	Birds []Birdcontent `json:"birds"`
}

type Birdcontent struct {
	Name      string `json:"name"`
	Initial   string `json:"initial"`
	Id        string `json:"id"`
	SourceUrl string `json:"sourceurl"`

	Status string `json:"status"`
	Intro  string `json:"intro"`
	Latin  string `json:"latin"`
	Family string `json:"family"`

	Where string `json:"where"`
	When  string `json:"when"`
	Diet  string `json:"diet"`

	Population   BirdPopulation   `json:"population"`
	Distribution BirdDistribution `json:"distribution"`

	Images []string `json:"images"`
}

type BirdPopulation struct {
	Europe       string `json:"europe,omitempty"`
	UK_Breeding  string `json:"uk_breeding,omitempty"`
	UK_Wintering string `json:"uk_wintering,omitempty"`
	UK_Passage   string `json:"uk_passage,omitempty"`
}

type BirdDistribution struct {
	Europe    string `json:"europe"`
	UK        string `json:"uk"`
	Worldwide string `json:"worldwide"`
}

type Birdurl struct {
	Id string `json:"id"`

	Url  string `json:"url"`
	Name string `json:"name"`
}

func ReadFile() (birds []Birdurl) {

	file, err := ioutil.ReadFile("./data/index.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(file, &birds)

	return
}

func PullData(bird Birdurl) (b Birdcontent) {

	doc, err := goquery.NewDocument(bird.Url)

	if err != nil {
		panic(err)
	}

	b.Name = bird.Name
	// Index to the first character of the string
	b.Initial = string(bird.Id[0])
	b.Id = bird.Id
	b.SourceUrl = bird.Url

	{ // status
		status_obj := doc.Find("#guide-bocc-status a")
		status, _ := status_obj.Attr("id")
		b.Status = strings.TrimPrefix(status, "bocc-")
	}

	{ // page-content

		ct := doc.Find("#page-content")

		{ //intro
			b.Intro = ct.Find(".intro").Text()
		}
		{ //latin
			b.Latin = ct.Find("h3:containsOwn('Latin name') + p").Text()
		}
		{ //where
			b.Where = ct.Find("h3:containsOwn('Where to see them') + p").Text()
		}
		{ //when
			b.When = ct.Find("h3:containsOwn('When to see them') + p").Text()
		}
		{ //diet
			b.Diet = ct.Find("h3:containsOwn('What they eat') + p").Text()
		}

		{ //family

			rx := regexp.MustCompile(`/(?P<family>[^/\.]*)(?:\.aspx|)?$`)

			family, _ := ct.Find("h3:containsOwn('Family') + p a").Attr("href")
			match := rx.FindStringSubmatch(family)
			b.Family = match[1]

		}

		{ //population

			p := BirdPopulation{}
			ob := ct.Find("h3:containsOwn('Population') + table")

			p.Europe = ob.Find("td:nth-child(1)").Text()
			p.UK_Breeding = ob.Find("td:nth-child(2)").Text()
			p.UK_Passage = ob.Find("td:nth-child(3)").Text()
			p.UK_Wintering = ob.Find("td:nth-child(4)").Text()

			if p.Europe == "-" {
				p.Europe = ""
			}
			if p.UK_Breeding == "-" {
				p.UK_Breeding = ""
			}
			if p.UK_Passage == "-" {
				p.UK_Passage = ""
			}
			if p.UK_Wintering == "-" {
				p.UK_Wintering = ""
			}

			b.Population = p

		}

		{ // Images

			files, err := ioutil.ReadDir("./images/" + bird.Id)
			if err != nil {
				panic(err)
			}
			for _, f := range files {
				b.Images = append(b.Images, f.Name())
			}
		}

	}

	{ // extras

		ob := doc.Find("#extras")

		{ // distribution
			d := BirdDistribution{}

			d.Europe = ob.Find("dt:containsOwn('In the UK') + dd").Text()
			d.UK = ob.Find("dt:containsOwn('In Europe') + dd").Text()
			d.Worldwide = ob.Find("dt:containsOwn('Worldwide') + dd").Text()

			b.Distribution = d
		}

	}

	return
}

func main() {

	birds := ReadFile()
	container := Birdscontainer{}

	// Read the json
	for i, bird := range birds {

		b := PullData(bird)

		container.Birds = append(container.Birds, b)

		data, err := json.MarshalIndent(b, "", "\t")
		if err != nil {
			panic(err)
		}
		location := "./content/" + b.Id
		err = os.MkdirAll(location, os.ModePerm)
		ioutil.WriteFile(location+"/index.md", data, 0644)

		if i > 1 {
			// TODO Include when testing
			// break
		}
	}

	// Write to one file
	data, err := json.MarshalIndent(container, "", "\t")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("data/birds.json", data, 0644)
}
