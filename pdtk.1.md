
USAGE:

  ./bin/pdtk [OPTIONS] verb [VERB_OPTIONS] [-- [PANDOC_OPTIONS] ... ]

./bin/pdtk started as a Pandoc preprocessor. It can read JSON 
or YAML from standard input and passes that via an internal 
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" marker. Options before "--" are for
the ./bin/pdtk preprossor. 

./bin/pdtk has grown to include features provide through simple
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


OPTIONS

  -help       display usage
  -license    display license
  -version    display version

BASIC EXAMPLES

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool ./bin/pdtk.

  ./bin/pdtk prep -- --template example.tmpl < example.json

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

  ./bin/pdtk prep -- -s -t markdown < example.json

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

  ./bin/pdtk prep -i codemeta.json -o about.md \
             -- --template codemeta-md.tmpl

Using ./bin/pdtk to manage blog content with the "blogit"
verb. 

Adding a blog "first-post.md" to "myblog".

  ./bin/pdtk blogit myblog $HOME/Documents/first-post.md

Adding/Updating the "first-post.md" on "2022-07-22"

  ./bin/pdtk blogit myblog $HOME/Documents/first-post.md "2022-07-22"

Added additional material for posts on "2022-07-22"

  ./bin/pdtk blogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"

Refreshing the blogs's blog.json file.

  ./bin/pdtk blogit myblog

Using ./bin/pdtk to generate RSS for "myblog"

  ./bin/pdtk rss myblog

Generating a sitemap in a current directory

  ./bin/pdtk sitemap .

Running a static web server to view rendering site

  ./bin/pdtk ws $HOME/Sites/myblog

