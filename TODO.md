TODO
====


Bugs
----

- [X] RSS Feed isn't passing validation,
    - [X] Dates must comply with http://www.faqs.org/rfcs/rfc822.html
        - NOTE: For Go time.Time what passes validators is RFC1123Z formatted dates.
    - [X] Missing Atom relationship element
    - [X] Links need to point at .html files not .md files
    - Make sure it passes with at least two validators
        - [X] W3C validator: https://validator.w3.org/feed/
        - [X] RSS Board: https://www.rssboard.org/rss-validator/


Next
----

- [ ] Review Gopher and see about adding Gopher support
    - [ ] Look at gophermap and see how it may tranlsate to/from RSS
    - [ ] Look at autogenerating gophermap from blog.json
- [ ] Review Gemini and see about adding Gemini support
- [ ] Review yarn.social and see if it make sense to support in toolkit
- [ ] Review [Micropubs spec](https://micropub.spec.indieweb.org/)
- [ ] Review [JSONFeed spec](https://www.jsonfeed.org/)
- [ ] Review Micro.blog's [Archive Format](https://book.micro.blog/blog-archive-format/)
    - [ ] Review implementation used by [Daring Fireball](https://daringfireball.net/feeds/json)
- [ ] Add support for [JSONFeed](https://www.jsonfeed.org/) 1.1
    - [ ] Evaluate if that could be the cannonical feed used to render RSS 2.0 XML and Atom XML
    - [ ] Figure out how JSONFeed plays with what Dave Winer is doing
- [ ] Make sure my feeds look good in Dave Winer's reallysimple readers
    - [ ] Add support for source namespace so I can do source:markdown element per Dave Winer's reallysimple feeds
      - [ ] rss should be able to understandard a markdown/html file relationship and transform it into RSS 2 with Markdown source elements where appropriate
    - See http://source.scripting.com/#1653758422000
    - Example, see http://feeder.scripting.com/?template=mailbox&feedurl=http://scripting.com/rss.xml
    - Example, see http://feeder.scripting.com/?template=mailbox&feedurl=https://rsdoiel.github.io/index.xml
    - Example, see http://feeder.scripting.com/?template=mailbox&feedurl=https://rsdoiel.github.io/rss.xml
- [x] blog.json needs to contain enough metadata to easily render the RSS feeed. The addtional data could be set via blogit options
- [x] I need to support generating multiple feeds for a website, e.g. site, blog, article series
    - [x] rss should be able to produce a "feed" for all pages in a website using Markdown document's front matter where there is a matching html document
    - [x] rss should be able to produce a "feed" for a selected set of pages driven from YAML front matter elements like "series" name
- [ ] sitemap needs to be implemented and support links to sub-site maps
- [ ] I need to render an index listing pages from Front Matter of content pages
    - [ ] Review how Rmarkdown/RStudio handle inclusion by front matter switches
