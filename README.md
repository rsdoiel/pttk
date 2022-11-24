
[![Project Status: Concept – Minimal or no implementation has been done yet, or the repository is only intended to be a limited example, demo, or proof-of-concept.](https://www.repostatus.org/badges/latest/concept.svg)](https://www.repostatus.org/#concept)

pttk
====

**pttk** is a plain text tool kit for writing. The main focus is on static site generation supplimenting tools like Pandoc.  The metaphor behind the tool kit is a deconstructed content management system. It is easily scripted from your favorite POSIX shell or Makefile. It provides a number of functions including a blogging tool called blogit, phlogit, JSONfeed, RSS generations and rudementry support for sitemap.xml. Combines with Pandoc and Pagefind you can easily build rich websites and blogs.

A command line tool kit
-----------------------

**pttk** is a program that works on the command line or shell.
**pttk** usage is structured around the idea "verbs" or actions.
Each "verb" can have it's own set of options and command syntax.

The basic usage is as follows

```
   pttk VERB [OPTIONS]
```

Currently there are four verbs supported by **pttk**.

__blogit__
: A tool for manage a blog directory structure and a
"blog.json" metadata file

__rss__
: A tool for generating RSS files from blogit

__sitemap__
: A tool for generating a sitemap.xml file.

__include__
: A "include" text preprocessor including files with via an "#include();" directive.

__blogit__ is a tool to make it easy to separate website generation
from where you might want to write your blog posts. It will generate
and maintain a blog style directory structure. A blog directory structure
is usually in the form of `/YYYY/MM/DD/` where "YYYY" is a year, "MM" is
a two digit month and "DD" is a two digit day representation. It also
maintains a "blog.json" document that describes the metadata and layout for
your blog. __blogit__ uses the front matter in your Markdown documents to
know things like titles, post dates and authorship.  The two **pttk**
verbs "rss" and "sitemap" know how to interpret the blog.json to generate
RSS and sitemap.xml respectively.

The form of the __blogit__ command is

```shell
    pttk blogit PATH_TO_DOCUMENT_TO_IMPORT [YYYY_MM_DD]
```

In this example I have a Markdown document I want to use as a blog post
in `$HOME/Documents/pttk-tutorial.md`.  I'm generating my blog in a
directory called `$HOME/Sites/my-website/blog`.  If I want to "blog" the
document I first change to "my-website" directory and use __blogit__
to update my blog.

```shell
   cd $HOME/Sites/my-website/blog
   pttk blogit $HOME/Documents/pttk-tutorial.md
```

The __blogit__ verb assumes you are in the current working directory
where you have your blog.


By default __blogit__ will use the current date in "YYYY-MM-DD" format
for the blog post. If you want to have the post on a specific day then
you include the date for the post in "YYYY-MM-DD" format. Here's an
example of posting the tutorial on 2022-08-01 (August 8th, 2022).

```shell
   cd $HOME/Sites/my-website/blog
   pttk blogit $HOME/Documents/pttk-tutorial.md 2022-08-08
```

__rss__ is the verb used to generate an RSS feed from a __blogit__
blog.json file.  The format of the command is

```shell
    pttk rss PATH_TO_BLOG_JSON PATH_TO_RSS_FILE
```

If I want my blog feed to be `feeds/index.xml` in the Wordpress style
for my blog in the `blog` directory I would change to `my-website`
directory and then use the __rss__ as follows.

```shell
    cd $HOME/Sites/my-website
    pttk rss blog/blog.json feeds/index.xml
```

This will generate our `feeds/index.xml` document. If the feeds directory
doesn't exist it'll get created. Updating the RSS picking up new post
is just a matter of invoking `pttk rss` the command again.

__sitemap__ generates a "sitemap.xml" file that describes the site layout
to searching crawlers.  The specification for sitemap.xml stipulates a
maximum number of entries in the sitemap.xml. For large websites this used
to be a problem but the specification allows for multiple sitemaps to be
used.  The __sitemap__ verb will generate a sitemap.xml in the root
website directory and in any sub-directories of the website.  If Markdown
documents are found then it'll use front matter for the matching HTML files
and "blog.json" file for the blog content.

The form for __sitemap__ is simple.

```
   ptdk sitemap [ROOT_WEBSITE_DIRECTORY]
```

Here's an example for our "my-website" directory.

```
    pttk sitemap $HOME/Sites/my-website
```

This wold generate a sitemap file of `$HOME/Sites/my-website/sitemap.xml`
and if necessary ones in the sub directories like `blog`.


Go package
----------

Here's some simple use examples of the three functions supplied
in the pttk package.

Given a JSON Object document  as a slice of bytes render formatted
output based on the Pandoc template `example.tmpl`

```go
    src, err := ioutil.ReadFile("example.json")
    if err != nil {
        // ... handle error
    }
    // options passed to Pandoc
    opt := []string{"--template", "example.tmpl"}
    src, err = pttk.Apply(src, opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
```

Using an `io.Reader` to retrieve the JSON content, process with the
`example.tmpl` template and write standard output

```go
    f, err := Open("example.json")
    if err != nil {
        // ... handle error
    }
    defer f.Close()
    // options passed to Pandoc
    opt := []string{"--template", "example.tmpl"}
    src, err := pttk.ReadAll(f, opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
```

Using an `io.Reader` and `io.Writer` to read JSON source from standard
input and write the processed Pandoc templated standard output.

```go
    // options passed to Pandoc
    opt := []string{"--template", "example.tmpl"}
    err := pttk.ApplyIO(os.Stdin, os.Stdout, opt)
    if err != nil {
        // ... handle error
    }
```

Requirements
------------

- [Pandoc](https://pandoc.org) 2.18 or better
- [Go](https://golang.org) 1.18.4 or better to compile from source
- [GNU Make](https://www.gnu.org/software/make/) (optional) to automated compilation
- [Git](https://git-scm.com/) or other Git client to retrieve this repository

Installation
------------

1. Clone https://github.com/rsdoiel/pttk to your local machine
2. Change directory into the git repository (i.e. `pttk`
3. Compile using `go build`
4. Install using `go install`

```shell
    git clone https://github.com/rsdoiel/pttk
    cd pttk
    git fetch origin
    git pull origin main
    go build -o bin/pttk cmd/pttk/pttk.go
    go install
```

NOTE: This recipe assumes' you are familiar with setting up a
Go development environment (e.g. You've set GOPATH environment
appropriately). See the [go website](https://golang.org) for
details about setting up and compiler programs.

License
-------

BSD 3-Clause License

Copyright (c) 2022, R. S. Doiel
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

