package ws

const (
	helpText = `

USAGE

    {app_name} {verb} [HTDOC_PATH] [OPTIONS]

SYNOPSIS

{app_name} {verb} provides a simple static web server for
testing the content you're rendering with Pandoc (or
other static site generator).

EXAMPLE

In the example the htdoc directory is called "myblog"
and you can view the result at http://localhost:8000.

  {app_name} {verb} $HOME/Sites/myblog

`
)
