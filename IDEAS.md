
Ideas/Musings
=============

This is a document is an exploration of what pdtk could be. 

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
- Part of writing is reading, how can pdtk help with that?
    - A feed READ might be nice
- Part of writing includes citing things and looking things up, how can pdtk help with that?
    - A personal search would be nice
    - A way to read/reference freeds like Dave Wiener's feed mailbox with a link to cite it my blog or read it in Pocket
    - A way to display OPML of interesting sites
- Should your site search be simple (site only) or a "personal search engine"?
    - A "personal search engine" would take a set of OPML files, RSS feeds index their content for easy retrieval
    - A "personal search engine" would return links to the object, citation and link to "read in pocket"
    - A "personal search engine" would also index documentation sites so you could easily look up code questions with out getting the spam and poor quality links of comercial general search engines (e.g. Google, Bing, DuckDuckGo)

Potential Verbs (actions)
-------------------------

- A set of verbs to processing more complex actions
    - "prep" should expose the simple pre-processor implemented in pdtmpl experiment, this could be used to do things like generate an "about.md" page directly from codemeta.json of a project.
    - "blogit" should be a fresh implementation of the blogit CLI from mkpage
    - "rss" should generate appropriate, modern RSS 2.0 for syndication
    - "sitemap" should generate a sitemap for a website and support sub-sitemaps were appropraite
    - "search" would let you search your site's content perhaps using LunrJS indexes
        - "search" also might function as a personal search engine
    - "index" would generate a search index like LunrJS indexes used by the "search" verb
        - should be able to index all Markdown content and front matter
        - should be able to index content provided via an OPML file
        - should be able to index content provided by an RSS feed
        - should be able to index content distributed as JSON such as the context synsitive help you have in VSCode from Mozilla developer website, Python or Golang's documentation websites
    - "reader" would provide a feed reading experiences similar to Dave Weiner's mailbox feed reader
        - could be browser based building on a localhost web server
        - should include a citation link
        - should include a "read in pocket" link
        - should include a link to the original content


Current Innovations
-------------------

- Dave Weiner's feeder is a great way of reading a feed. As Dave pointed out my rss.xml doesn't look great, need to fix that. The test URL is [http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Frsdoiel.github.io/rss.xml](http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Frsdoiel.github.io/rss.xml)
    - Compare that with James Fallows [http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Ffallows.substack.com%2Ffeed](http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Ffallows.substack.com%2Ffeed)


Go packages to look at
----------------------

- gofeed package provides a universal reader, appears to be still maintained
- Gorrilla provides a RSS generator, Gorrilla is looking for a new maintainer and appears to have stopped evolving (i.e. is current stable, may not have a community to update it in the future)

