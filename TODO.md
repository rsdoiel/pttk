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

- [ ] Evaluate re-implenting in TypeScript as compiled `deno` executable
- [ ] Evaluate if [flatlake](https://flatlake.app) can replace pttk's blog.json/phlog.json generation, if so then I need to think about integrate that and to render my RSS and sitemaps from the flatlake api directory
  - [ ] If I adopt flatlake then I should integrate sitemap generation for the static api directory
  - [ ] I should be able to generate RSS, JSONfeed and Atom feeds from the static api directory
- [ ] Decide if pttk should be rewritten in Rust or TypeScript rather than Go, is it essentially is a "runner" of sorts and integrating with Rust projects like flatlake and pagefind could make pttk a little more compelling.
- [ ] Sitemap should include sub-sitemaps.xml when needed, the one I'm generating is pretty much useless
- [ ] A gophermap generated to work in `pttk gs` doesn't work like the Gopherserver on SDF, they need to match so gophermaps are portable.
- [ ] Think about leverage skimmer's database of harvested feed items and how that might integrate into blogging (e.g. quoting articles)
- [ ] I need a way I can read gophermaps, twtxt, JSONfeed, RSS/Atom feeds in a single reader, preferrably a console app, it should be driven from an OPML file or simple text file like newsboat (wish newsboat supported subscriptions to gophermaps, twtxt then I could just translate JSONfeed to twtxt...)
- [ ] Review git.mills.io/prologic/go-gopher and understand what was implemented
    - [ ] Evaluate using as is and what I would need need to write to replace gophermap handling with how it works on sdf.org's Gopher deployment
    - [ ] Evaluate writing my own gopher server for previewing content easily is it turns out I need to fork go-gopher to get the behavior I want
    - [ ] Look at serving table content out of SQLite3 database files
    - [ ] Look at integrations with dataset
    - [ ] Gopher service needs to minic Gophernicus in how it handles Gophermaps and menus
- [ ] Review Gopher and see about adding Gopher support
    - [ ] Look at gophermap and see how it may tranlsate to/from RSS
    - [ ] Look at autogenerating gophermap from blog.json
    - [ ] Update .editorconfig to fixed issues about converting tabs in gophermap
    - [ ] Look at modifying go-gopher to support "+" operator in the selector for hostname and port, or just auto create them in not supplied.
        - gopher.Dir() is probably what I need to look at to modify.
        - See http://gopherinfo.somnolescent.net/servers/customizing-menus/ for description of server side processing of gophermap
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
- [ ] sitemap needs to be implemented and support links to sub-site maps
- [ ] I need to render an index listing pages from Front Matter of content pages
    - [ ] Review how Rmarkdown/RStudio handle inclusion by front matter switches
- [x] blog.json needs to contain enough metadata to easily render the RSS feeed. The addtional data could be set via blogit options
- [x] I need to support generating multiple feeds for a website, e.g. site, blog, article series
    - [x] rss should be able to produce a "feed" for all pages in a website using Markdown document's front matter where there is a matching html document
    - [x] rss should be able to produce a "feed" for a selected set of pages driven from YAML front matter elements like "series" name
- [X] Remove prep/pandoc as it is not needed, recent versions of Pandoc include `--metadata-file` to ingest JSON as metadata
