
Ideas
=====

Questions
---------

- What would be the minimum additional tools besides GNU Make and Pandoc needed to make it easy to maintain my website?
    - blogit for generating JSON presentation of blog post and rendering to the appropriate location of blog
    - An RSS generator would be nice
    - A working sitemap.xml generator would be nice
    - feeder for generating reallysimple feeds and RSS 2.0 feeds in each directory
    - A simple static web server for viewing the rendered content would ne nice
    - A navigation file generator based on bread crumbs


Improvements
------------

- A set of verbs to processing more complex actions
    - a "walker" verb that would walk the file system looking for JSON document to turn into Markdown then into HTML so that the blog is always in a data format
    - a "rss" verb might render the RSS for the site/blog
        - should support Dave's recent innovations for including Markdown via a "source" name space
    - a "sitemap" verb might render a series of sitemap files for the site
        - these would change together to avoid overflowing the limit on files in a sitemap.xml file.
    - a "reallysimple" would produce JSON feed documents based on Dave's recent innovations
    



