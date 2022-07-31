TODO
====


Bugs
----

- [ ] RSS Feed isn't passing validation, 
    - [X] Dates must comply with http://www.faqs.org/rfcs/rfc822.html
        - NOTE: For Go time.Time what passes validators is RFC1123Z formatted dates.
    - [ ] Missing Atom relationship element
    - [ ] Links need to point at .html files not .md files
    - Make sure it passes with at least two validators
        - [ ] W3C validator: https://validator.w3.org/feed/
        - [ ] RSS Board: https://www.rssboard.org/rss-validator/


Next
----

- [ ] blog.json needs to contain enough metadata to easily render the RSS feeed. The addtional data could be set via blogit options
- [ ] rss should be able to produce a "feed" for all pages in a website using Markdown document's front matter where there is a matching html document
- [ ] rss should be able to understandard a blog.json file and transform it into RSS 2 with Markdown source elements where appropriate
- [ ] sitemap needs to be implemented and support links to sub-site maps
- [ ] need "byline" and "titleline" extractors implemented
