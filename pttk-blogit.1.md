% pttk-blogit(1) pttk-blogit user manual
% R. S. Doiel
% August 14, 2022

# NAME

pttk

# SYNOPSIS

pttk blogit [OPTIONS]

pttk blogit [OPTIONS] -stn STN_FILENAME

# DESCRIPTION

pttk blogit provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.

__pttk blogit__ also includes an option to extract short (one paragraph) blog posts froom [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file.

# OPTIONS

What follows are the options supported by the blogit verb.

-asset
: Copy asset file to the blog path for provided date (YYYY-MM-DD)

-copyright string
: Set the blog copyright notice.

-description string
: Set the blog description

-ended string
: Set the blog ended date.

-help
: display blogit help

-index-tmpl string
: Set index blog template

-language string
: Set the blog language. (default "en-US")

-license string
: Set the blog language license

-name string
: Set the blog name.

-post-tmpl string
: Set index blog template

-prefix string
: Set the prefix path before YYYY/MM/DD directory structure.

-quip string
: Set the blog quip.

-refresh string
: This will create/refresh the blog.json file for given year(s), if more than one year is to be refresh separate each year with a comma, no spaces.  E.g. "2021,2022,2023"

-save-as-yaml
: save as YAML file instead of blog.yaml file

-started string
: Set the blog started date.

-url string
: Set blog's URL

-stn
: Import short blog posts from an [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file

-author
: Set the "author" string when importing from a simple timesheet notation file.

-verbose
: verbose output

# EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2021.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
use the following command.

~~~shell
   cd Sites/me.example.org
   pttk blogit my-vacation-day.md 2021-07-01
~~~

The *pttk blogit* command will copy "my-vacation-day.md",
creating any necessary file directories to
"Sites/me.example.org/2021/06/01".  It will also update article
lists (index.md) at the year level, month, and day level and month
level of the directory tree and and generate/update a posts.json
in the "Sites/my.example.org" that can be used in your home page
template for listing recent posts.

*pttk blogit* includes an option to set the prefix path to
the blog posting.  In this way you could have separate blogs
structures for things like podcasts or videocasts.

~~~shell
    # Add a landing page for the podcast
    pttk blogit -prefix=podcast my-vacation.md 2021-07-01
    # Add an audio file containing the podcast
    pttk blogit -prefix=podcast my-vacation.wav 2021-07-01
~~~

Where "-prefix" sets the prefix path before the YYYY/MM/DD path.


If you have an existing blog paths in the form of
PREFIX/YYYY/MM/DD you can use blogit to create/update/recreate
the blog.json file.

~~~shell
    pttk blogit -prefix=blog -refresh=2021
~~~

The option "-refresh" is what indicates you want to crawl
for blog posts for that year.


In this final example I am updating blog posts from a [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file called "project-log.txt". I am sending those blog posts to the
prefix directory "blog" and using the author name, "Jane Doe".

~~~
    pttk blogit -prefix=blog -author 'Jane Doe' -stn project-log.txt
~~~

This will create individual, time stamp titled posts for each of the simple timesheet notation entries found in "project-log.txt".


# SEE ALSO

- manual pages for [pttk](pttk.1.html), [pttk-prep](pttk-prep.1.html), [pttk-rss](pttk-rss.1.html)
- pttk website at [https://rsdoiel.github.io/pttk](https://rsdoiel.github.io/pttk)
- The source code is available from [https://github.com/rsdoiel/pttk](https://github.com/rsdoiel/pttk)
- Simple timesheet notation at [https://rsdoiel.github.io/stngo/docs/stn.html](https://rsdoiel.github.io/stngo/docs/stn.html)



