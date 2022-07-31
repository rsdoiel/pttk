% pdtk-blogit(1) pdtk user manual
% R. S. Doiel
% July, 31, 2022

# NAME

pdtk blogit

# SYNOPSIS

pdtk [OPTIONS] blogit [BLOGIT_OPTIONS] POST_MARKDOWN_FILE [YYYY_MM_DD]

# DESCRIPTION

pdtk blogit provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.

# EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2021.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
use the following command.

   cd Sites/me.example.org
   pdtk blogit my-vacation-day.md 2021-07-01

The *pdtk blogit* command will copy "my-vacation-day.md",
creating any necessary file directories to 
"Sites/me.example.org/2021/06/01".  It will also update article 
lists (index.md) at the year level, month, and day level and month
level of the directory tree and and generate/update a posts.json
in the "Sites/my.example.org" that can be used in your home page
template for listing recent posts.

*pdtk blogit* includes an option to set the prefix path to
the blog posting.  In this way you could have separate blogs 
structures for things like podcasts or videocasts.

    # Add a landing page for the podcast
    pdtk blogit -prefix=podcast my-vacation.md 2021-07-01
    # Add an audio file containing the podcast
    pdtk blogit -prefix=podcast my-vacation.wav 2021-07-01

Where "-p, -prefix" sets the prefix path before the YYYY/MM/DD path.


If you have an existing blog paths in the form of
PREFIX/YYYY/MM/DD you can use blogit to create/update/recreate
the blog.json file.

    pdtk blogit -prefix=blog -refresh=2021

The option "-refresh" is what indicates you want to crawl
for blog posts for that year.

# SEE ALSO

pdtk website at https://rsdoiel.github.io/pdtk

The source code is avialable from https://github.com/rsdoiel/pdtk


