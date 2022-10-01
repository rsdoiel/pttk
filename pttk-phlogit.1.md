% pdtk-phlogit(1) pdtk-phlogit user manual
% R. S. Doiel
% August 14, 2022

# NAME

pdtk

# SYNOPSIS

pdtk phlogit [OPTIONS]

pdtk phlogit [OPTIONS] -stn STN_FILENAME

# DESCRIPTION

pdtk phlogit provides a quick tool to add or replace phlog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.

__pdtk phlogit__ also includes an option to extract short (one paragraph) phlog posts froom [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file.

# OPTIONS

What follows are the options supported by the phlogit verb.

-asset
: Copy asset file to the phlog path for provided date (YYYY-MM-DD)

-copyright string
: Set the phlog copyright notice.

-description string
: Set the phlog description

-ended string
: Set the phlog ended date.

-help
: display phlogit help

-index-tmpl string
: Set index phlog template

-language string
: Set the phlog language. (default "en-US")

-license string
: Set the phlog language license

-name string
: Set the phlog name.

-post-tmpl string
: Set index phlog template

-prefix string
: Set the prefix path before YYYY/MM/DD.

-quip string
: Set the phlog quip.

-refresh string
: This will create/refresh the phlog.json file for given year(s), if more than one year is to be refresh separate each year with a comma, no spaces.  E.g. "2021,2022,2023"

-save-as-yaml
: save as YAML file instead of phlog.yaml file

-started string
: Set the phlog started date.

-url string
: Set phlog's URL

-stn
: Import short phlog posts from an [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file

-author
: Set the "author" string when importing from a simple timesheet notation file.

-verbose
: verbose output

# EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my phlog for the date July 1, 2021.  I've written
"my-vacation-day.md" in my home "Documents" folder and my phlog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the phlog me.example.org would
use the following command.

~~~shell
   cd Sites/me.example.org
   pdtk phlogit my-vacation-day.md 2021-07-01
~~~

The *pdtk phlogit* command will copy "my-vacation-day.md",
creating any necessary file directories to
"Sites/me.example.org/2021/06/01".  It will also update article
lists (index.md) at the year level, month, and day level and month
level of the directory tree and and generate/update a posts.json
in the "Sites/my.example.org" that can be used in your home page
template for listing recent posts.

*pdtk phlogit* includes an option to set the prefix path to
the phlog posting.  In this way you could have separate phlogs
structures for things like podcasts or videocasts.

~~~shell
    # Add a landing page for the podcast
    pdtk phlogit -prefix=podcast my-vacation.md 2021-07-01
    # Add an audio file containing the podcast
    pdtk phlogit -prefix=podcast my-vacation.wav 2021-07-01
~~~

Where "-prefix" sets the prefix path before the YYYY/MM/DD path.


If you have an existing phlog paths in the form of
PREFIX/YYYY/MM/DD you can use phlogit to create/update/recreate
the phlog.json file.

~~~shell
    pdtk phlogit -prefix=phlog -refresh=2021
~~~

The option "-refresh" is what indicates you want to crawl
for phlog posts for that year.


In this final example I am updating phlog posts from a [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file called "project-log.txt". I am sending those phlog posts to the
prefix directory "phlog" and using the author name, "Jane Doe".

~~~
    pdtk phlogit -prefix=phlog -author 'Jane Doe' -stn project-log.txt
~~~

This will create individual, time stamp titled posts for each of the simple timesheet notation entries found in "project-log.txt".


# SEE ALSO

- manual pages for [pdtk](pdtk.1.html), [pdtk-prep](pdtk-prep.1.html), [pdtk-rss](pdtk-rss.1.html)
- pdtk website at [https://rsdoiel.github.io/pdtk](https://rsdoiel.github.io/pdtk)
- The source code is available from [https://github.com/rsdoiel/pdtk](https://github.com/rsdoiel/pdtk)
- Simple timesheet notation at [https://rsdoiel.github.io/stngo/docs/stn.html](https://rsdoiel.github.io/stngo/docs/stn.html)



