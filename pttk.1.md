% pttk(1) pttk user manual
% R. S. Doiel
% July, 22, 2022

# NAME

pttk - a writers tool kit for static site generation using Pandoc

# SYNOPSIS

pttk [OPTIONS] verb [VERB_OPTIONS] [\-\- [PANDOC_OPTIONS] ... ]

# DESCRIPTION

pttk is a toolkit for writing.  The main focus is on static
site generation with Pandoc. The tool kit provides those missing
elements from a deconstructed content management system that
Pandoc does not (possibly should not) provide. Using pttk with
Pandoc should provide the core features expected in producing
a website or blog in the early 21st century. These include
preprocessor called "prep" which lets you take a JSON file
and transform it into Markdown front matter that is
directly passed to Pandoc for processing. You might want to
do this for things like generating a CITATION.cff or an
about page from a codemeta.json. You might want to generate
BibTeX from your Markdown pages front matter you collected.
pttk includes a tool called "blogit" that manages taking a
Markdown source document and placing it in a blog directory
structure while maintaining a blogs metadata in a "blog.json"
file. It includes a tool, "rss", that generates RSS files for a
website or blog.  There is even a localhost web server
for previewing content called "ws".  All these tools are easily
scripted via a Makefile or your favorite programming language
(e.g. Python, Lua, Oberon-07, Go).

"Verbs" are the way you select the tool you want to work
with in the tool kit, e.g. "prep", "blogit", "rss" or "ws".

## Meet the VERBS

pttk tools are accessed through "verbs". These
"verbs" currently include the following.

**help**
: Display this help page.

**blogit**
: Renders a blog directory structure by "importing" Markdown documents
or updating existing ones. It maintains a blog.json document collecting
metadata and supporting RSS rendering.

**include**
: A preprocessor of doing recursive includes using an include directive like `#include(myfile.md);`

**prep**
: Preprocess JSON or YAML into YAML front matter and run through Pandoc

**rss**
: Renders RSS feeds from the contents of a blog.json document

**sitemap**
: Renders sitemap.xml files for a static website

**ws**
: Runs a simple static web server for previewing content in your web browser


**gs**
: Run a simlpe static gopher server for previewing content in your gopher client

# OPTIONS

-help
: display usage

-license
: display license

-version
: display version

# EXAMPLES


## blogit

Adding a blog "first-post.md" to "myblog".

```shell
    pttk blogit myblog $HOME/Documents/first-post.md
```

Adding/Updating the "first-post.md" on "2022-07-22"

```shell
    pttk blogit myblog $HOME/Documents/first-post.md "2022-07-22"
```

Added additional material for posts on "2022-07-22"

```shell
    pttk blogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"
```

Refreshing the blog's blog.json file.

```shell
    pttk blogit myblog
```


## prep

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "`<`" is used to pipe the content of "example.json"
into the command line tool pttk.

```shell
    pttk prep -- --template example.tmpl < example.json
```

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

```shell
    pttk prep -- -s -t markdown < example.json
```

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandoc's template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

```shell
    pttk prep -i codemeta.json -o about.md \
        -- --template codemeta-md.tmpl
```

Using pttk to manage blog content with the "blogit"
verb.


## rss

Using pttk to generate RSS for "myblog"

```shell
    pttk rss myblog
```

## sitemap

Generating a sitemap in a current directory

```shell
    pttk sitemap .
```


## ws

Running a static web server to view rendering site

```shell
    pttk ws $HOME/Sites/myblog
```

## gs

Running a static gopher server to view rendering site


# SEE ALSO

- manual pages for [pttk-prep](pttk-prep.1.html), [pttk-blogit](pttk-blogit.1.html), [pttk-rss](pttk-rss.1.html), [pttk-ws](pttk-ws.1.html)
- pttk website at [https://rsdoiel.github.io/pttk](https://rsdoiel.github.io/pttk)
- The source code is available from [https://github.com/rsdoiel/pttk](https://github.com/rsdoiel/pttk)


