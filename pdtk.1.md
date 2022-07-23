% pdtk(1) pdtk user manual
% R. S. Doiel
% July, 22, 2022

# NAME

pdtk - pandoc tool kit, a set of tools for website generation using Pandoc


# SYNOPSIS

pdtk [OPTIONS] verb [VERB_OPTIONS] [-- [PANDOC_OPTIONS] ... ]

# DESCRIPTION

pdtk started as a Pandoc preprocessor. It can read JSON 
or YAML from standard input and passes that via an internal 
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" marker. Options before "--" are for
the pdtk preprossor. 

# VERBS

pdtk has grown to include features provide through simple
"verbs". The verbs include the following.

**help**
: Display this help page.

**prep**
: Preprocess JSON or YAML into YAML front matter and run through Pandoc

**ws**
: Runs a simple static web server for checking static site development

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

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool pdtk.

    pdtk prep -- --template example.tmpl < example.json

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

    pdtk prep -- -s -t markdown < example.json

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

    pdtk prep -i codemeta.json -o about.md \
        -- --template codemeta-md.tmpl

Using pdtk to manage blog content with the "blogit"
verb. 

Adding a blog "first-post.md" to "myblog".

    pdtk blogit myblog $HOME/Documents/first-post.md

Adding/Updating the "first-post.md" on "2022-07-22"

    pdtk blogit myblog $HOME/Documents/first-post.md "2022-07-22"

Added additional material for posts on "2022-07-22"

    pdtk blogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"

Refreshing the blogs's blog.json file.

    pdtk blogit myblog

Using pdtk to generate RSS for "myblog"

    pdtk rss myblog

Generating a sitemap in a current directory

    pdtk sitemap .

Running a static web server to view rendering site

    pdtk ws $HOME/Sites/myblog

# SEE ALSO

pdtk website at https://rsdoiel.github.io/pdtk

The source code is avialable from https://github.com/rsdoiel/pdtk


