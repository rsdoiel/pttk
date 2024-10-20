// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
import {
	Writer,
	FmtHelp,
	Version,
	ReleaseDate,
	ReleaseHash,
	LicenseText,
} from './deps.ts';
import { parseArgs } from "@std/cli/parse-args";

const help_text = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

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

**gophermap**
: Renders a gophermap file for a Gopher whole

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

function handle_error(eout: Writer, err: string) {
	if (err !== undefined)  {
		eout.writeSync(`${err}`);
		Deno.exit(1);
	}
}

function main() {
	let app_name: string = path.basename(Deno.args[0]);

	const flags = parseArgs(Deno.args, 
		boolean: ["help", "license", "version" ]
	});

	let args: string[] = [...Deno.args],
	    //in = Deno.in,
	    out = Deno.stdout,
	    eout = Deno.stderr;

	if (flags.help) {
        out.writeSync(FmtHelp(help_text, app_name, Version, ReleaseDate, ReleaseHash));
		Deno.exit(0);
	}
	if (flags.version) {
		out.WriteSync(`${app_name} ${Version} ${ReleaseDate} ${ReleaseHash}`);
		Deno.exit(0);
	}
	if (flags.license) {
		out.writeSync(LicenseText);
		Deno.exit(0);
	}

	if (args.length === 0) {
        eout.writeSync(FmtHelp(help_text, app_name, Version, ReleaseDate, ReleaseHash));
		Deno.exit(1);
	}
	const verb: string = args.shift();
	
	switch verb {
	case "help":
        fmt.Fprintln(out, FmtHelp(help_text, app_name, Version, ReleaseDate, ReleaseHash));
		Deno.exit(0);
/*
	case "frontmatter":
		if err := frontmatter.RunFrontmatter(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "ws":
		if err := ws.RunWS(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "gs":
		if err := gs.RunGS(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "blogit":
		if err := blogit.RunBlogIt(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "phlogit":
		if err := phlogit.RunPhlogIt(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "gophermap":
		if err := phlogit.RunGophermap(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "rss":
		src, err := rss.RunRSS(appName, verb, args)
		handle_eerror(eout, err)
		if len(src) > 0 {
			fmt.Fprintf(out, "%s\n", src)
		}
	case "include":
		if err := include.RunInclude(appName, verb, args); err != nil {
			handle_error(eout, err)
		}
	case "sitemap":
		handle_error(eout, fmt.Errorf("%s %s not implemented", appName, verb))
*/
	default:
        eout.writeSync(FmtHelp(help_text, app_name, Version, ReleaseDate, ReleaseHash));
		Deno.exit(1);
	}
}
