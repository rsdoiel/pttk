% pdtk(1) pdtk user manual
% R. S. Doiel
% July, 22, 2022

# NAME

pdtk - a writers tool kit for static site generation using Pandoc

# SYNOPSIS

pdtk [OPTIONS] verb [VERB_OPTIONS] [-- [PANDOC_OPTIONS] ... ]

# DESCRIPTION

pdtk is a toolkit for writing.  The main focus is on static
site generation with Pandoc. The tool kit provides those missing
elements from a deconstructed content management system that
Pandoc does not (possibly should not) provide. Using pdtk with
Pandoc should provide the core features expected in producing
a website or blog in the early 21st century. These include
preprocessor called "prep" which lets you take a JSON file
and transform it into Markdown front matter that is
directly passed to Pandoc for processing. You might want to
do this for things like generating a CITATION.cff or an 
about page from a codemeta.json. You might want to generate
BibTeX from your Markdown pages front matter you collected.
pdtk includes a tool called "blogit" that manages taking a
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

pdtk tools are accessed through "verbs". These
"verbs" currently include the following.

**help**
: Display this help page.

**blogit**
: Renders a blog directory structure by "importing" Markdown documents
or updating existing ones. It maintains a blog.json document collecting
metadata and supporting RSS rendering.

**prep**
: Preprocess JSON or YAML into YAML front matter and run through Pandoc

**rss**
: Renders RSS feeds from the contents of a blog.json document

**sitemap**
: Renders sitemap.xml files for a static website

**ws**
: Runs a simple static web server for previewing content in your web browser


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
    pdtk blogit myblog $HOME/Documents/first-post.md
```

Adding/Updating the "first-post.md" on "2022-07-22"

```shell
    pdtk blogit myblog $HOME/Documents/first-post.md "2022-07-22"
```

Added additional material for posts on "2022-07-22"

```shell
    pdtk blogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"
```

Refreshing the blog's blog.json file.

```shell
    pdtk blogit myblog
```

## prep

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "`<`" is used to pipe the content of "example.json"
into the command line tool pdtk.

```shell
    pdtk prep -- --template example.tmpl < example.json
```

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

```shell
    pdtk prep -- -s -t markdown < example.json
```

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandoc's template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

```shell
    pdtk prep -i codemeta.json -o about.md \
        -- --template codemeta-md.tmpl
```

Using pdtk to manage blog content with the "blogit"
verb. 

## rss

Using pdtk to generate RSS for "myblog"

```shell
    pdtk rss myblog
```

## sitemap

Generating a sitemap in a current directory

```shell
    pdtk sitemap .
```

## ws

Running a static web server to view rendering site

```shell
    pdtk ws $HOME/Sites/myblog
```

# SEE ALSO

- manual pages for [pdtk-prep](pdtk-prep.1.html), [pdtk-blogit](pdtk-blogit.1.html), [pdtk-rss](pdtk-rss.1.html), [pdtk-ws](pdtk-ws.1.html)
- pdtk website at [https://rsdoiel.github.io/pdtk](https://rsdoiel.github.io/pdtk)
- The source code is available from [https://github.com/rsdoiel/pdtk](https://github.com/rsdoiel/pdtk)


