---
title: pdtk
section: 1
header: User Manual
footer: pdtk 0.0.0
date: July 22, 2022
---

# NAME

pdtk - a pandoc preprocessor for JSON and YAML content.

# SYNOPSIS

  pdtk [OPTIONS] [-- [PANDOC_OPTIONS] ...]

# DESCRIPTION

pdtk is a Pandoc preprocessor. It reads JSON from standard
input, transforms it into YAML front matter suitable for
processing via Pandoc. The common usecase would be to render
via a pandoc template. By default pdtk reads from standard
input and writes standard out.

# OPTIONS

**-help**
: display usage

**-license**
: display license

**-version**
: display version

**-i FILENAME**
: read JSON or YAML file

**-o FILENAME**
: write Pandoc output to file

# EXAMPLE

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect `"<"` is used to pipe the content of "example.json"
into the command line tool pdtk.

```shell
  pdtk prep -- --template example.tmpl < example.json
```

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

```shell
  pdtk prep -- -s -t markdown < example.json
```

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

```shell
  pdtk prep -i codemeta.json -o about.md \
       -- --template codemeta-md.tmpl
```

