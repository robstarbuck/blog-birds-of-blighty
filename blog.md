# InDesign Alternative using GoLang and Chrome

A client came to me asking if could put together some Infographics along with some page-layouts. The first thing that came to mind was InDesign which isn't software that I'm all that familiar with. Not wanting to commit the time nor pay £18 a month to use the software, I turned the business down.

A week later, it occurred to me that I could quite easily have used html and CSS and output it as a PDF through the __Save as PDF__ functionality in chrome. In short web-developer can expand their service with the tools they already know. I'd even say they'd be better equipped.

Keen to explore what could be achieved I set out to create a guide to British birds (Birds of Blighty), that couldn't easily be done in InDesign. This post is really just a demonstration of what's possible, you can download the finished result here.

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

The body of our PDF is going to be produced from a single page of HTML. __/layouts/section/birds.html__ is the only layout we need to produce our page.

A huge selling point of creating PDFs from HTML is we can use fragment identifiers to navigate our PDF. As we have some 267 birds to peruse I'm going to make them navigable by first initial as well as family.

This works exactly as it would on the web, with fragments are referenced with a hash prefix `#` linking to an element of the same id.

    <a href="#initial_{{ .Key }}">
      <span>{{ .Key }}</span>
    </a>

    {{ range .Data.Pages.GroupByParam "initial" }}
      <section class="cf initial" id="initial_{{ .Key }}">
      …
      </section>
    {{ end }}

I'm defining absolute dimensions for the page 94mm x 151mm, the same size as my Nexus 7. Importantly we define this against the @page selector, now chrome will know what size to output our document.

    @page{
        overflow:hidden;
        // Nexus 2013
        size: 94mm 151mm;
        margin:10mm;
        margin-bottom:15mm;
    }

The first initials page and every bird page will have no overflow, so we're specifying the same measurements for width and height.

    .paper,
    .bird,
    {
        page-break-before: always;
        width: 94mm;
        height: 151mm;
        box-sizing: border-box;
        position: relative;
    }

__page-break-before__ state the top of the page should run over "the fold" of the page above. Effectively you're saying "start this on a new page". __page-break-after__ and __page-break-before__ are similarly useful these properties are enormous selling points as no-one want their content arbitrarily spilling into new pages.

With width and height set we can even stretch our content over the height of the page with `height: 100%;`.

Coding for print takes a little getting used to, bringing up the print-preview can prove a little arduous but it's a cinch once you're used to it. 

# Finishing Touches

The last thing to do is export our document, in Chrome simply set your destination to "Save as PDF" and ensure that "Background graphics" are ticked and the job's done.

To finishing things off, I added a front and back cover to the PDF. This would have proven difficult in the HTML as our page margins were set globally which meant any background image or colour would still be pushed in at the margins. Instead I put together the designs in Affinity Designer* and combined the files in OSX's Preview with drag and drop.

# RESOURCES

[Repository](https://github.com/robstarbuck/blog-birds-of-blighty)
[Takana, SASS preprocessor for sublime users](https://github.com/mechio/takana)
[Designing for Print with CSS](https://www.smashingmagazine.com/2015/01/designing-for-print-with-css/)
[Selector for print (@page)](https://developer.mozilla.org/en/docs/Web/CSS/@page) 

