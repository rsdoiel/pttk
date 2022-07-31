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
elements from a deconstricted content management system that
Pandoc does not (possibly should not) provide. Using pdtk with
Pandoc should provide the core features expected in producing
a website or blog in the early 21st century. These include
a prep processor called "prep" which lets you take JSON
and transform it into Markdown frontmatter that is
directly passed to Pandoc for processing. You might want to
do this for things like generating a CITATION.cff from a
codemeta.json file or an about page from a codemeta.json file.
Included is "blogit" tool that manages taking a Markdown source
document and placing it in a blog directory struction while
maintaining a "blog.json" metadata file describing the whole blog.
Another tool is "rss" that generates RSS files for a website or
blog. All these tools are easily scripted via Makefile or your
favorite programming language (e.g. Python, Lua, Go).

Finally there is a tool called "ws" which provides you with
a "localhost" web server for preview your website or blog
before you publish it.

"Verbs" are the way you select the tool you want to work
with in the tool kit, e.g. prep, blogit, rss.

# VERBS

pdtk has grown to include features provide through simple
"verbs". The verbs include the following.

**help**
: Display this help page.

**prep**
: Preprocess JSON or YAML into YAML front matter and run through Pandoc

**ws**
: Runs a simple static web server for previewing content in your web browser

**blogit**
: Renders a blog directory structure by "importing" Markdown documents
or updating existing ones. It maintains a blog.json document collecting
metadata and supportting RSS rendering.

**rss**
: Renders RSS feeds from the contents of a blog.json document

**sitemap**
: Renders sitemap.xml files for a static website


# OPTIONS

-help
: display usage

-license
: display license

-version
: display version

# EXAMPLES

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
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

```shell
    pdtk prep -i codemeta.json -o about.md \
        -- --template codemeta-md.tmpl
```

Using pdtk to manage blog content with the "blogit"
verb. 

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

Refreshing the blogs's blog.json file.

```shell
    pdtk blogit myblog
```

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
- The source code is avialable from [https://github.com/rsdoiel/pdtk](https://github.com/rsdoiel/pdtk)


