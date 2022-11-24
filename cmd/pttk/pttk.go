// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rsdoiel/pttk"
	"github.com/rsdoiel/pttk/blogit"
	"github.com/rsdoiel/pttk/frontmatter"
	"github.com/rsdoiel/pttk/gs"
	"github.com/rsdoiel/pttk/include"
	"github.com/rsdoiel/pttk/phlogit"
	"github.com/rsdoiel/pttk/rss"
	"github.com/rsdoiel/pttk/ws"
)

const (
	helpText = `% {app_name}(1) {app_name} user manual
% R. S. Doiel
% August 18, 2022

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] verb [VERB_OPTIONS] [\-\- [PANDOC_OPTIONS] ... ]

# DESCRIPTION


{app_name} implements a deconstructed content management system suitable for
working with plain text. It intended as a compliment to Pandoc focusing
on collections of documents and structed text.  Currently {app_name} provides
tools to layout blog directories, generate Gophermap files for Gopher
distribution.  The ideas is to provide the tooling that will allow
publication and distribution both on the world wide web as well as
the "small internet".

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


# EXAMPLES

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

## phlogit verb

Using {app_name} to manage phlog content with the "phlogit"
verb.

Adding a blog "first-post.md" to "myphlog".

~~~shell
  {app_name} phlogit myphlog $HOME/Documents/first-post.md
~~~

Adding/Updating the "first-post.md" on "2022-07-22"

~~~shell
  {app_name} phlogit myblog $HOME/Documents/first-post.md "2022-07-22"
~~~

Added additional material for posts on "2022-07-22"

~~~shell
  {app_name} phlogit myblog $HOME/Documents/charts/my-graph.svg "2022-07-22"
~~~

Refreshing the phlogs's phlog.json file.

~~~shell
  {app_name} phlogit myblog
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
	appName := path.Base(os.Args[0])

	flag.BoolVar(&showHelp, "help", false, "display usage")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.Parse()
	args := flag.Args()

	if showHelp {
		fmt.Print(usage(appName))
		os.Exit(0)
	}
	if showVersion {
		fmt.Printf("%s %s\n", appName, pttk.Version)
		os.Exit(0)
	}
	if showLicense {
		fmt.Printf("%s\n", pttk.LicenseText)
		fmt.Printf("%s %s\n", appName, pttk.Version)
		os.Exit(0)
	}

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
	case "frontmatter":
		if err := frontmatter.RunFrontmatter(appName, verb, args); err != nil {
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
	case "phlogit":
		if err := phlogit.RunPhlogIt(appName, verb, args); err != nil {
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
