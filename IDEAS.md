
Ideas/Musings
=============

This is a document is an exploration of what pttk could be.

A Question and Possibilities
----------------------------

- What would be the minimum additional tools besides GNU Make and Pandoc needed to make it easy to maintain my website?
    - blogit, a tool for generating JSON presentation of blog posts and their metadata
    - An RSS generator would be nice
    - A working sitemap.xml generator would be nice
    - Support for recent innovations that Dave Weiner has done with "reallysimple" feeds
    - A simple static web server for viewing the review in the content rendered with blogit or pandoc
    - A way to integrate fountain manuscripts
    - A nice to have would be a way to generate breadcrumb navigation for both blogs and website
        - Maybe this is driven by the top level index page's front matter? Or an index.yaml file?
    - A nice to have would be a LunrJS compatible index generator for integrate search
- What would be the minimum additional tools besides GNU Make and Pandoc needed to make it easy to maintain my Gopher site?
- What would be the minimum additional tools besides GNU Make and Pandoc needed to make it easy to maintain my Gemini site?
- Part of writing is reading, how can pttk help with that?
    - A feed reader might be nice (e.g. similar to yarnc or twet output or something more rich like newsboat)
    - Need to read from RSS, twtxt.txt, stno
    - Render feed list as OPML to share on feedland
    - Need to support reading Gophermap as feed source (so I can skim phlogs too)
    - Need to support Gemini Atom feeds as feed source (so I can skim Gemlogs too)
    - Integrate following feeds from feedland
- Part of writing includes citing things and looking things up, how can pttk help with that?
    - A personal search engine would be nice, Pagefind shows a very good approach but doesn't search from the command line
    - A way to read/reference freeds like Dave Wiener's feed mailbox with a link to cite it in my blog or read it in Pocket
    - A way to display OPML of interesting sites
- Should your site search be simple (site only) or a "personal search engine"?
    - Pagefind implements a very good website static site
    - Pagefind doesn't support remote resources inless a pagefind bunldle can be referenced
    - Pagefind doesn't work on Gopher and Gemini sites
    - Pagefind search UI is GUI browser based, doesn't work in lynx or as a cli
- A Personal search engine
    - A "personal search engine" would take a set of OPML files, JSONfeeds, twtxt.txt and RSS feeds index their content for easy retrieval
    - A "personal search engine" would return links to the object, citation and link to "read in pocket"
    - A "personal search engine" would also index documentation sites so you could easily look up code questions with out getting the spam and poor quality links of comercial general search engines (e.g. Google, Bing, DuckDuckGo)

Potential Verbs (actions)
-------------------------

- A set of verbs to processing more complex actions
    - "blogit" should be a fresh implementation of the blogit CLI from mkpage
    - "phlogit" should be manage and organizing a Gopher phlog
    - "series" would manage both a curated feed and series table of content based on Markdown frontmatter in files
    - "rss" should generate appropriate, modern RSS 2.0 for syndication
    - "jsonfeed" should generate appropriate JSON feeds for syndication
    - "sitemap" should generate a sitemap for a website and support sub-sitemaps were appropraite
    - "search" would let you search your site's content, if possible via something similar to pagefind
        - "search" also might function as a personal search engine
    - "index" would generate a search index like LunrJS or pagefind indexes used by the "search" verb
        - should be able to index all Markdown content and front matter
        - should be able to index content provided via an OPML file
        - should be able to index content provided by an RSS feed
        - should be able to index content distributed as JSON such as the context synsitive help you have in VSCode from Mozilla developer website, Python or Golang's documentation websites
    - "reader" would provide a feed reading experiences similar to Dave Weiner's mailbox feed reader
        - could be browser based building on a localhost web server
        - should include a citation link
        - should include a "read in pocket" link
        - should include a link to the original content
    - "readit" would take a list of URLs (e.g. .newsboat/url), scan them and gerate a pocket reading list, or a RSS file that could be fed to Dave Weiner's mailbox style feed reader


Current Innovations
-------------------

- Dave Weiner's feeder is a great way of reading a feed. As Dave pointed out my rss.xml doesn't look great, need to fix that. The test URL is [http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Frsdoiel.github.io/rss.xml](http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Frsdoiel.github.io/rss.xml)
    - Compare that with James Fallows [http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Ffallows.substack.com%2Ffeed](http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Ffallows.substack.com%2Ffeed)


Go packages to look at
----------------------

- gofeed package provides a universal reader, appears to be still maintained
- Gorrilla provides a RSS generator, Gorrilla is looking for a new maintainer and appears to have stopped evolving (i.e. is current stable, may not have a community to update it in the future)

