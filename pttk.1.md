%pttk(1) skimmer user manual | version 0.0.17 5a5db9c
% R. S. Doiel
% 2024-10-21

# NAME

pttk

# SYNOPSIS

pttk [OPTIONS] verb [VERB_OPTIONS] [\-\- [PANDOC_OPTIONS] ... ]

# DESCRIPTION


pttk implements a deconstructed content management system suitable for
working with plain text. It intended as a compliment to Pandoc focusing
on collections of documents and structed text.  Currently pttk provides
tools to layout blog directories, generate Gophermap files for Gopher
distribution.  The ideas is to provide the tooling that will allow
publication and distribution both on the world wide web as well as
the "small internet".

pttk has grown to include features provide through simple
"verbs". The verbs include the following.

# OPTIONS

-help
: Display help

-license
: Display license

-version
: Display version

# VERBS

Verbs have their one options. You can see a list of them
with the form `pttk VERB -h`

**help**
: Display this help page.

**ws**
: Runs a simple static web server for checking static site development

**gs**
: Runs a simple Gopher service for static site development

**frontmatter**
: Reads a Pandoc markdown file with frontmatter and write out JSON

**blogit**
: Renders a blog directory structure by "importing" Markdown documents
or updating existing ones. It maintains a blog.json document collecting
metadata and supportting RSS rendering.

**phlogit**
: Renders a Phlog directory structure by "importing" text files
or updating existing ones. It maintains a phlog.json document collecting
metadata and supporting RSS rendering as well as generating gophermap files.

**include**
: Include any files indicated by an include directive (e.g. "#include(toc.md);"). Include operates recursively so included files can also include other files.

**rss**
: Renders RSS feeds from the contents of a blog.json document

**sitemap**
: Renders sitemap.xml files for a static website

**gophermap**
: Renders a gophermap file for a Gopher whole

# EXAMPLES

## blogit verb

Using pttk to manage blog content with the "blogit"
verb.

Adding a blog "first-post.md" to "myblog".

~~~shell
  pttk blogit myblog $HOME/Documents/first-post.md
~~~

Adding/Updating the "first-post.md" on "2022-07-22"

~~~shell
  pttk blogit myblog $HOME/Documents/first-post.md "2022-07-22"
~~~

Added additional material for posts on "2022-07-22"

~~~shell
  pttk blogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"
~~~

Refreshing the blogs's blog.json file.

~~~shell
  pttk blogit myblog
~~~

## phlogit verb

Using pttk to manage phlog content with the "phlogit"
verb.

Adding a blog "first-post.md" to "myphlog".

~~~shell
  pttk phlogit myphlog $HOME/Documents/first-post.md
~~~

Adding/Updating the "first-post.md" on "2022-07-22"

~~~shell
  pttk phlogit myblog $HOME/Documents/first-post.md "2022-07-22"
~~~

Added additional material for posts on "2022-07-22"

~~~shell
  pttk phlogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"
~~~

Refreshing the phlogs's phlog.json file.

~~~shell
  pttk phlogit myblog
~~~

## rss verb

Using pttk to generate RSS for "myblog"

~~~shell
  pttk rss myblog
~~~

## sitemap verb

Generating a sitemap in a current directory (i.e. the "." directory)

~~~shell
  pttk sitemap .
~~~

## ws verb

Running a static web server to view rendering site

~~~shell
  pttk ws $HOME/Sites/myblog
~~~

## gs verb

Running a static gopher server to view rendering site

~~~
  pttk gs $HOME/Sites/myblog
~~~

## include verb

Including a table of contents "toc.md", and "chapters1.md"
and "chapters2.md" in a file called "book.txt" and writing
the result to "book.md".

The "book.txt" file would look like

~~~
   # My Book

   #include(toc.md);

   #include(chapter1.md);

   #include(chapter2.md);
~~~

Putting the "book" together as on file.

~~~shell
	pttk {verb} book.txt book.md
~~~


