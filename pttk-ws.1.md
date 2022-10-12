% pttk-ws(1) pttk-ws user manual
% R. S. Doiel
% September 23, 2022

# NAME

pttk ws

# USAGE

pttk ws [HTDOC_PATH] [URL_TO_LISTEN_FOR] [OPTIONS]

# SYNOPSIS

pttk ws provides a simple static web server for
testing the content you're rendering with Pandoc (or
other static site generator).

# EXAMPLE

In the example the htdoc directory is called "myblog"
and you can view the result at http://localhost:8000.

  pttk ws $HOME/Sites/myblog


