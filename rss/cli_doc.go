package rss

const (
	helpText = `
NAME

   {app_name} {verb} PATH_TO_SITE

SYNOPSIS

The {verb} renders an RSS file based on the content found in the
directory tree provided. If it encounters a "blog.json" file then
it'll use that file to generate feed content for that directory
and it's content otherwise it'll generate a feed backed on Markdown
front matter encountered in Markdown documents with corresponding
html file.

DESCRIPTION

EXAMPLE

ALSO SEE

`
)
