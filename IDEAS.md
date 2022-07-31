
Ideas
=====

A Question and Ideas
--------------------

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

Verbs
-----

- A set of verbs to processing more complex actions
    - "prep" should expose the simple pre-processor implemented in pdtmpl experiment
    - "blogit" should be a fresh implementation of the blogit CLI from mkpage
    - "rss" should generate appropriate, modern RSS 2.0 for syndication
    - "sitemap" should generate a sitemap for a website and support sub-sitemaps were appropraite
    


Improvements
------------

Dave Weiner's feeder is a great way of reading a feed. As Dave pointed out
my rss.xml doesn't look great, need to fix that. The test URL is http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Frsdoiel.github.io/rss.xml

Compare that with James Fallows http://feeder.scripting.com/?template=mailbox&feedurl=https%3A%2F%2Ffallows.substack.com%2Ffeed



Feeds
-----

- gofeed package provides a universal reader, appears to be still maintained
- Gorrilla provides a RSS generator, Gorrilla is looking for a new maintainer and appears to have stopped evolving (i.e. is current stable, may not have a community to update it in the future)

