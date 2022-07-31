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

- [ ] Make sure my feeds look good in Dave Weiner's reallysimple readers
    - See http://source.scripting.com/#1653758422000
    - Example, see http://feeder.scripting.com/?template=mailbox&feedurl=http://scripting.com/rss.xml
    - Example, see http://feeder.scripting.com/?template=mailbox&feedurl=https://rsdoiel.github.io/index.xml
    - Example, see http://feeder.scripting.com/?template=mailbox&feedurl=https://rsdoiel.github.io/rss.xml
- [ ] Add support for source namespace so I can do source:markdown element per Dave Weiner's reallysimple feeds
    - [ ] rss should be able to understandard a markdown/html file relationship and transform it into RSS 2 with Markdown source elements where appropriate
- [x] blog.json needs to contain enough metadata to easily render the RSS feeed. The addtional data could be set via blogit options
- [x] I need to support generating multiple feeds for a website, e.g. site, blog, article series
    - [x] rss should be able to produce a "feed" for all pages in a website using Markdown document's front matter where there is a matching html document
    - [x] rss should be able to produce a "feed" for a selected set of pages driven from YAML front matter elements like "series" name
- [ ] sitemap needs to be implemented and support links to sub-site maps
- [ ] I need to render an index listing pages from Front Matter of content pages
    - [ ] Review how Rmarkdown/RStudio handle inclusion by front matter switches
