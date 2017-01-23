# InDesign Alternative using GoLang and Chrome

A client came to me asking if could put together some Infographics along with some page-layouts. The first thing that came to mind was InDesign which isn't software that I'm all that familiar with. Not wanting to commit the time nor pay £18 a month to use the software, I turned the business down.

A week later, it occurred to me that I could quite easily have used html and CSS and output it as a PDF through the __Save as PDF__ functionality in chrome. In short web-developer can expand their service with the tools they already know. I'd even say they'd be better equipped.

Keen to explore what could be achieved I set out to create a guide to British birds (Birds of Blighty), that couldn't easily be done in InDesign.

# Scraping Data with Golang

First of all, we need information on birds and I know of a society that can help (no names mentioned). I use the Golang for scripting wherever I can, I really encourage developers to take the tour*, the Gopher put me off at first but I've learned to love it, it's the one nonsensical thing about GoLang.

I split the scrape into stages. Yes I could have nested loops, I've seen Inception, I understood it but let's keep things simple, in the words Rob Pike, "Complexity is multiplicative".

The scripts I used are here:

__index.go__ - Gets a list of all the bird urls from the site
__images.go__ - Downloads images for each bird
__content.go__ - Creates a json file for each bird and registers any images that have been downloaded for it

The code is [here](https://github.com/robstarbuck/birds-of-blighty-scrape), please get in touch if you've any questions. This post is not about Golang though I'll summarise a few points:

Using [GoQuery](https://github.com/PuerkitoBio/goquery) allows the use of jQuery selectors to scrape info from a page. A list of available selectors can be found [here](https://github.com/andybalholm/cascadia/blob/master/selector_test.go). Testing selectors in Chrome DevTools beforehand makes life easy.

<!-- TODO Image of chrome dev tools -->

__json/encoding__ is used to _marshal_ our data. You'll see a lot of examples of field 'tags' in the code. Tags allow us to change the property names we store, but fields must be "exported" in order to be included which means they must begin with capital letters. `omitempty` acts as it suggests since we don't have data for every field.

`json:"uk_wintering,omitempty"`

The reason our data is pulled into separate __.md__ files containing markdown will be understood in the next section.

I'm using panic rather than returning errors, for scripts like returning errors is a little laborious though I recognise it's good practice. 

    if err != nil {
        panic(err)
    }


# Producing HTML with Hugo and Takana

Our bird data is now contained within our __index.md__ files.

    ├── content
    │   ├── aquaticwarbler
    │   │   └── index.md
    │   ├── arcticskua
    │   │   └── index.md
    │   └── arctictern
    │       └── index.md

    {
        "name": "Aquatic warbler",
        "initial": "a",
        "id": "aquaticwarbler",
        "sourceurl": "/birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/a/aquaticwarbler/index.aspx",
        "status": "red",
        "intro": "The aquatic warbler is a regular but scarce autumn migrant to certain areas in southern Britain, visiting on its way between breeding grounds in eastern Europe and its winter home in West Africa.  Its dependence on a specialised and vulnerable breeding habitat means it has become a globally threatened and declining species. It is more yellow-brown and streaked than the simliar sedge warbler.",
        "latin": "Acrocephalus paludicola",

This data structure was chosen to be used with GoHugo*. A static site generator written in Golang, I now favour it's use over Jekyll (used for github pages), it's faster and more configurable. It's good to use SASS rather than CSS. Jekyll supports SASS pre-processing whereas Golang doesn't. I use Takana* to get around this problem.

# RESOURCES

[SASS preprocessor for sublime users](https://github.com/mechio/takana)