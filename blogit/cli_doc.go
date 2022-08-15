// blogit.go is a sub-package pdtk. A packages for managing static content
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package blogit

const (
	helpText = `% {app_name}-{verb}(1) {app_name}-{verb} user manual
% R. S. Doiel
% August 14, 2022

# NAME

{app_name}

# SYNOPSIS

{app_name} {verb} [OPTIONS]

# DESCRIPTION

{app_name} {verb} provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.

{app_name} {verb} also can ingest a "stn" (Simple Timesheet Notation)
file and render the contents as blog posts.

# OPTIONS

-stn STN_FILENAME
: Import a [Simple Timesheet Notation](https://rsdoiel.github.io/stngo)
file as blog content.

-prefix BLOG_PATH
: Sets the path prefix before the YYYY/MM/DD path created for a blog entry

-refresh YEAR_LIST
: This will refresh the blog.json file for given year(s), if more than
one year is to be refresh separate each year with a comma, no spaces.
E.g. "2021,2022,2023"


# EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2021.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
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
the blog posting.  In this way you could have separate blogs 
structures for things like podcasts or videocasts.

~~~shell
    # Add a landing page for the podcast
    {app_name} {verb} -prefix=podcast my-vacation.md 2021-07-01
    # Add an audio file containing the podcast
    {app_name} {verb} -prefix=podcast my-vacation.wav 2021-07-01
~~~

Where "-p, -prefix" sets the prefix path before the YYYY/MM/DD path.


If you have an existing blog paths in the form of
PREFIX/YYYY/MM/DD you can use blogit to create/update/recreate
the blog.json file.

~~~shell
    {app_name} {verb} -prefix=blog -refresh=2021
~~~

The option "-refresh" is what indicates you want to crawl
for blog posts for that year.


Take "LogBook.txt" which contains [Simple Timesheet Notation](https://rsdoiel.github.io/stngo) and use that to generate blog entries in the 
blog with the prefix path "blog".

~~~
	{app_name} {verb} -prefix=blog -stn=LogBook.txt
~~~

`
)
