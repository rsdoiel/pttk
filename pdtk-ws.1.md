% pdtk-ws(1) pdtk-ws user manual
% R. S. Doiel
% September 23, 2022

# NAME

pdtk ws

# USAGE

pdtk ws [HTDOC_PATH] [URL_TO_LISTEN_FOR] [OPTIONS]

# SYNOPSIS

pdtk ws provides a simple static web server for
testing the content you're rendering with Pandoc (or
other static site generator).

# EXAMPLE

In the example the htdoc directory is called "myblog"
and you can view the result at http://localhost:8000.

  pdtk ws $HOME/Sites/myblog


