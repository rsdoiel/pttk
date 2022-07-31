% pdtk-ws(1) pdtk-ws user manual
% R. S. Doiel
% July, 31, 2022

NAME

pdtk ws - a static site web server for localhost and content review

SYNOPSIS

pdtk ws [OPTIONS] [HTDOC_PATH]

DESCRIPTION

pdtk ws provides a simple static web server for testing the content
you're rendering with Pandoc (or other static site generator). By default
it displays content at http://localhost:8000 from the current working
directory. It is not intended to be use as a production web server.

EXAMPLE

In the example the htdocs directory is called "myblog"
and you can view the result at http://localhost:8000.

```shell
    pdtk ws $HOME/Sites/myblog
```

# SEE ALSO

- manual pages for [pdtk](pdtk.1.html), [pdtk-prep](pdtk-prep.1.html), [pdtk-blogit](pdtk-blogit.1.html), [pdtk-rss](pdtk-rss.1.html)
- pdtk website at [https://rsdoiel.github.io/pdtk](https://rsdoiel.github.io/pdtk)
- The source code is available from [https://github.com/rsdoiel/pdtk](https://github.com/rsdoiel/pdtk)


