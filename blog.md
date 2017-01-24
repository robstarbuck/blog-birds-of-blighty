A client came to me asking if could put together some page-layouts for a PDF. The tool for the job seemed to be InDesign which I'm not that familiar with. Not wanting to commit the time to learn it, nor pay £18 a month to use it, I turned the business down.

A week later, it occurred to me that I could have put it together with HTML and output it as a PDF with the "Save as PDF" option in chrome. In short, web-developers can become publishers using tools they already know. I'd even say they're better equipped.

Keen to explore what could be achieved I set out to create a guide to British Birds (Birds of Blighty), that couldn't easily be done in InDesign. This post is brief walk-through the process, you can see the finished result here (NB, hyperlinks won't work in preview).

[![](https://storage.googleapis.com/magpie-img/birds-of-blighty/cover.jpg)](https://drive.google.com/file/d/0B_tatgDj9IbyZ2htc01NaDRGalU/view?usp=sharing)

### Scraping Data with Golang

First of all, we need information on birds and I know of a society that can help (no names mentioned). I use Golang for scripting wherever I can, I encourage developers to take the tour&#42;, the Gopher put me off at first but I've learnt to love it since, it's the one nonsensical thing about GoLang.

I split the scrape into stages. Yes, I could have done the scrape in one go with a few loops, I've seen Inception, I understood it, but let's keep things simple.

The scripts I used are here:

__index.go__ - Gets a list of all the bird URLs from the site
__images.go__ - Downloads images for each bird
__content.go__ - Creates a json file for each bird and references any images  already downloaded.

The code is included in "Resources", get in touch if you've any questions. I'm not going through all the code though I'll summarise a few points:

To grab the info from the website, I'm using [GoQuery](https://github.com/PuerkitoBio/goquery) which facilitates pulling info with the use of jQuery selectors. A list of available selectors can be found [here](https://github.com/andybalholm/cascadia/blob/master/selector_test.go). Testing selectors in Chrome DevTools beforehand makes life easy.

![](https://storage.googleapis.com/magpie-img/birds-of-blighty/devtools-selectors.jpg)

__json/encoding__ is used to _marshal_ our data. You'll see a lot of examples of field 'tags' in the code. Tags allow us to change the property names we store, but fields must be "exported" in order to be included which means they must begin with capital letters. `omitempty` acts as it suggests since we don't have data for every field this prevents empty values being written to the json.

``` language-go
  Europe   string  `json:"europe,omitempty"`
```

The reason our data is pulled into separate __.md__ files containing markdown will be understood in the next section.

I'm using panic rather than returning errors in functions, returning errors is a little laborious for such a short script, though I recognise it's good practice. 

``` language-go
    if err != nil {
        panic(err)
    }
```

### HTML with GoHugo

Our bird data is now contained within our __index.md__ files.

``` language-textile
    ├── content
    │   ├── aquaticwarbler
    │   │   └── index.md
    │   ├── arcticskua
    │   │   └── index.md
    │   └── arctictern
    │       └── index.md
```

``` language-js
    {
        "name": "Aquatic warbler",
        "initial": "a",
        "id": "aquaticwarbler",
        "sourceurl": "/birds-and-wildlife/bird-and-wildlife-guides/bird-a-z/a/aquaticwarbler/index.aspx",
        "status": "red",
        "intro": "The aquatic warbler is a regular but scarce autumn migrant to certain areas in southern Britain, visiting on its way between breeding grounds in eastern Europe and its winter home in West Africa.  Its dependence on a specialised and vulnerable breeding habitat means it has become a globally threatened and declining species. It is more yellow-brown and streaked than the simliar sedge warbler.",
        "latin": "Acrocephalus paludicola",
…
```

This data structure was chosen to be used with GoHugo&#42;. A static site generator written in Golang, I prefer it to Jekyll (used for github pages), it's faster and more configurable. SASS is always a preference to CSS but whereas Jekyll supports SASS pre-processing, GoHugo doesn't. I use Takana&#42; to get around this problem.

The body of our PDF is going to be produced from a single page of HTML. "/layouts/section/birds.html" is the only layout we need to produce our page.

A huge selling point of creating PDFs from HTML is we can use "fragment identifiers" to navigate our PDF. As we have some 267 birds to peruse so I'm going to make them navigable by first initial as well as family.

This works exactly as it would on the web, with fragments referenced by a hash prefix `#` linking to an element of the same id.

``` language-html
    <a href="#initial_{{ .Key }}">
       <span>{{ .Key }}</span>
    </a>

    {{ range .Data.Pages.GroupByParam "initial" }}
      <section class="cf initial" id="initial_{{ .Key }}">
      …
      </section>
    {{ end }}
```

I'm defining absolute dimensions for the page (94mm x 151mm), the same size as my Nexus 7. Importantly we define this against the @page selector, our browser now knows what size to save our document.

``` language-css
    @page{
        overflow:hidden;
        // Nexus 2013
        size: 94mm 151mm;
        margin:10mm;
        margin-bottom:15mm;
    }
```

Bird pages will have no overflow, i.e. they won't paginate, so we're specifying the same absolute measurements for their width and height.

``` language-css
    .paper,
    .bird,
    {
        page-break-before: always;
        width: 94mm;
        height: 151mm;
        box-sizing: border-box;
        position: relative;
    }
```

The "page-break-before" property states the top of the page should run over "the fold" of the page above. Effectively you're saying "start this on a new page". "page-break-after" and "page-break-before" are similarly useful and these properties are enormous selling points as no-one wants their content arbitrarily spilling into new pages.

!["Coal Tit"](https://storage.googleapis.com/magpie-img/birds-of-blighty/coaltit.jpg)

!["Initials page"](https://storage.googleapis.com/magpie-img/birds-of-blighty/initials.jpg)

With width and height set we can even stretch our content over the height of the page with `height: 100%;`.

### Finishing Touches

The last thing to do is export our document, in Chrome simply set your destination to "Save as PDF" and ensure that "Background graphics" are ticked and the job's done.

![ "Chrome Print Dialogue" ](https://storage.googleapis.com/magpie-img/birds-of-blighty/print-dialogue.jpg)

To finishing things off, I added a front and back cover to the PDF. This would have proven difficult in the HTML as our page margins are set globally any background graphic would appear in a white frame. Instead I put together the designs in Affinity Designer&#42; and combined the files in OSX's Preview with drag and drop. I feared the process would break the page links but it didn't.

![ "Designing the cover with Affinity Designer"](https://storage.googleapis.com/magpie-img/birds-of-blighty/affinity.jpg)

Coding for print takes a little getting used to, bringing up the print-preview can prove a little arduous but it's a cinch once you're used to it. 

---
#### Resources

- [Github Repo](https://github.com/robstarbuck/blog-birds-of-blighty)
- [Go Tour](https://tour.golang.org/welcome/1)
- [GoQuery](https://github.com/PuerkitoBio/goquery)
- [Takana, SASS preprocessor for sublime users](https://github.com/mechio/takana)
- [Designing for Print with CSS](https://www.smashingmagazine.com/2015/01/designing-for-print-with-css/)
- [Selector for print (@page)](https://developer.mozilla.org/en/docs/Web/CSS/@page) 
- [Affinity Designer](https://affinity.serif.com/en-gb/)