// pttk.go is a package (with sub-packages) for managing static content
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rsdoiel/pttk"
	"github.com/rsdoiel/pttk/blogit"
	"github.com/rsdoiel/pttk/gs"
	"github.com/rsdoiel/pttk/include"
	"github.com/rsdoiel/pttk/prep"
	"github.com/rsdoiel/pttk/rss"
	"github.com/rsdoiel/pttk/ws"
)

const (
	licenseText = `
{app_name} {version}

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
`

	helpText = `% {app_name}(1) {app_name} user manual
% R. S. Doiel
% August 18, 2022

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] verb [VERB_OPTIONS] [-- [PANDOC_OPTIONS] ... ]

# DESCRIPTION

{app_name} started as a Pandoc preprocessor. It can read JSON
or YAML from standard input and passes that via an internal
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" marker. Options before "--" are for
the {app_name} preprossor.

{app_name} has grown to include features provide through simple
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
with the form ` + "`" + `{app_name} VERB -h` + "`" + `

**help**
: Display this help page.

**prep**
: Preprocess JSON or YAML into YAML front matter and run through Pandoc

**ws**
: Runs a simple static web server for checking static site development

**gs**
: Runs a simple Gopher service for static site development

**blogit**
: Renders a blog directory structure by "importing" Markdown documents
or updating existing ones. It maintains a blog.json document collecting
metadata and supportting RSS rendering.

**include**
: Include any files indicated by an include directive (e.g. "#include(toc.md);"). Include operates recursively so included files can also include other files.

**rss**
: Renders RSS feeds from the contents of a blog.json document

**sitemap**
: Renders sitemap.xml files for a static website


# EXAMPLES

## prep verb

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool {app_name}.

~~~shell
  {app_name} prep -- --template example.tmpl < example.json
~~~

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

~~~shell
  {app_name} prep -- -s -t markdown < example.json
~~~

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

~~~shell
  {app_name} prep -i codemeta.json -o about.md \
             -- --template codemeta-md.tmpl
~~~

## blogit verb

Using {app_name} to manage blog content with the "blogit"
verb.

Adding a blog "first-post.md" to "myblog".

~~~shell
  {app_name} blogit myblog $HOME/Documents/first-post.md
~~~

Adding/Updating the "first-post.md" on "2022-07-22"

~~~shell
  {app_name} blogit myblog $HOME/Documents/first-post.md "2022-07-22"
~~~

Added additional material for posts on "2022-07-22"

~~~shell
  {app_name} blogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"
~~~

Refreshing the blogs's blog.json file.

~~~shell
  {app_name} blogit myblog
~~~

## rss verb

Using {app_name} to generate RSS for "myblog"

~~~shell
  {app_name} rss myblog
~~~

## sitemap verb

Generating a sitemap in a current directory (i.e. the "." directory)

~~~shell
  {app_name} sitemap .
~~~

## ws verb

Running a static web server to view rendering site

~~~shell
  {app_name} ws $HOME/Sites/myblog
~~~

## gs verb

Running a static gopher server to view rendering site

~~~
  {app_name} gs $HOME/Sites/myblog
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
	{app_name} {verb} book.txt book.md
~~~

`
)

func fmtText(src string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "{app_name}", appName), "{version}", version)
}

func version(appName string) string {
	return fmt.Sprintf("%s %s\n", path.Base(appName), pttk.Version)
}

func license(appName string) string {
	return fmtText(licenseText, appName, pttk.Version)
}

func usage(appName string) string {
	return fmtText(helpText, appName, pttk.Version)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func main() {
	var (
		showHelp    bool
		showLicense bool
		showVersion bool
	)
	flag.BoolVar(&showHelp, "help", false, "display usage")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.Parse()
	args := flag.Args()

	if showHelp {
		fmt.Print(usage(os.Args[0]))
		os.Exit(0)
	}
	if showVersion {
		fmt.Print(version(os.Args[0]))
		os.Exit(0)
	}
	if showLicense {
		fmt.Print(license(os.Args[0]))
		os.Exit(0)
	}

	appName := path.Base(os.Args[0])
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", usage(appName))
		os.Exit(1)
	}
	verb := args[0]
	if len(args) == 1 {
		args = []string{}
	} else {
		args = args[1:]
	}

	switch verb {
	case "help":
		fmt.Printf("%s\n", usage(appName))
		os.Exit(0)
	case "prep":
		if err := prep.RunPrep(appName, verb, args); err != nil {
			handleError(err)
		}
	case "ws":
		if err := ws.RunWS(appName, verb, args); err != nil {
			handleError(err)
		}
	case "gs":
		if err := gs.RunGS(appName, verb, args); err != nil {
			handleError(err)
		}
	case "blogit":
		if err := blogit.RunBlogIt(appName, verb, args); err != nil {
			handleError(err)
		}
	case "rss":
		src, err := rss.RunRSS(appName, verb, args)
		handleError(err)
		if len(src) > 0 {
			fmt.Printf("%s\n", src)
		}
	case "include":
		if err := include.RunInclude(appName, verb, args); err != nil {
			handleError(err)
		}
	case "sitemap":
		handleError(fmt.Errorf("%s %s not implemented", appName, verb))
	default:
		fmt.Fprintf(os.Stderr, "%s\n", usage(appName))
		os.Exit(1)
	}
}
