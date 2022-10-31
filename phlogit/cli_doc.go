// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package phlogit

const (
	helpText = `% {app_name}-{verb}(1) {app_name}-{verb} user manual
% R. S. Doiel
% August 14, 2022

# NAME

{app_name}

# SYNOPSIS

{app_name} {verb} [OPTIONS]

{app_name} {verb} [OPTIONS] -stn STN_FILENAME

# DESCRIPTION

{app_name} {verb} provides a quick tool to add or replace phlog content
organized around a date oriented file path. In addition to
placing documents it also will simple gophermap documents
for inclusion in navigation.

__{app_name} {verb}__ also includes an option to extract short (one paragraph) phlog posts from [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file.

# OPTIONS

What follows are the options supported by the phlogit verb.

-asset
: Copy asset file to the phlog path for provided date (YYYY-MM-DD)

-copyright string
: Set the phlog copyright notice.

-description string
: Set the phlog description

-ended string
: Set the phlog ended date.

-help
: display phlogit help

-index-tmpl string
: Set index phlog template

-language string
: Set the phlog language. (default "en-US")

-license string
: Set the phlog language license

-name string
: Set the phlog name.

-post-tmpl string
: Set index phlog template

-prefix string
: Set the prefix path before YYYY/MM/DD.

-quip string
: Set the phlog quip.

-refresh string
: This will create/refresh the phlog.json file for given year(s), if more than one year is to be refresh separate each year with a comma, no spaces.  E.g. "2021,2022,2023"

-save-as-yaml
: save as YAML file instead of phlog.yaml file

-started string
: Set the phlog started date.

-url string
: Set phlog's URL

-stn
: Import short phlog posts from an [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file

-author
: Set the "author" string when importing from a simple timesheet notation file.

-verbose
: verbose output

# EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my phlog for the date July 1, 2021.  I've written
"my-vacation-day.md" in my home "Documents" folder and my phlog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the phlog me.example.org would
use the following command.

~~~shell
   cd Sites/me.example.org
   {app_name} {verb} my-vacation-day.md 2021-07-01
~~~

The *{app_name} {verb}* command will copy "my-vacation-day.md",
creating any necessary file directories to
"Sites/me.example.org/2021/06/01".  It will also update article
lists (index.md) at the year level, month, and day level and month
level of the directory tree and and generate/update a posts.json
in the "Sites/my.example.org" that can be used in your home page
template for listing recent posts.

*{app_name} {verb}* includes an option to set the prefix path to
the phlog posting.  In this way you could have separate phlogs
structures for things like podcasts or videocasts.

~~~shell
    # Add a landing page for the podcast
    {app_name} {verb} -prefix=podcast my-vacation.md 2021-07-01
    # Add an audio file containing the podcast
    {app_name} {verb} -prefix=podcast my-vacation.wav 2021-07-01
~~~

Where "-prefix" sets the prefix path before the YYYY/MM/DD path.


If you have an existing phlog paths in the form of
PREFIX/YYYY/MM/DD you can use phlogit to create/update/recreate
the phlog.json file.

~~~shell
    {app_name} {verb} -prefix=phlog -refresh=2021
~~~

The option "-refresh" is what indicates you want to crawl
for phlog posts for that year.


In this final example I am updating phlog posts from a [simple timesheet notation](https://rsdoiel.github.io/stngo/docs/stn.html) file called "project-log.txt". I am sending those phlog posts to the
prefix directory "phlog" and using the author name, "Jane Doe".

~~~
    pttk phlogit -prefix=phlog -author 'Jane Doe' -stn project-log.txt
~~~

This will create individual, time stamp titled posts for each of the simple timesheet notation entries found in "project-log.txt".


# SEE ALSO

- manual pages for [pttk](pttk.1.html), [pttk-prep](pttk-prep.1.html), [pttk-rss](pttk-rss.1.html)
- pttk website at [https://rsdoiel.github.io/pttk](https://rsdoiel.github.io/pttk)
- The source code is available from [https://github.com/rsdoiel/pttk](https://github.com/rsdoiel/pttk)
- Simple timesheet notation at [https://rsdoiel.github.io/stngo/docs/stn.html](https://rsdoiel.github.io/stngo/docs/stn.html)


`
)
