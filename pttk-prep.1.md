% pttk-prep(1) pttk-prep user manual
% R. S. Doiel
% August 18, 2022

# NAME

pttk prep

# SYNOPSIS

pttk prep [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME] [\-\- [PANDOC_OPTIONS] ... ]

# DESCRIPTION

pttk prep is a Pandoc preprocessor. It can read JSON
or YAML from standard input and passes that via an internal
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" (double dash) marker. Options before "--"
are for the pttk preprossor.

NOTE: prep requires the "pandoc-server" to be running.

# OPTIONS

-help
: display usage

-i
: read from a file instead of standard input

-o
: write to a file instead of standard output

# EXAMPLES

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool pttk.

~~~shell
    pttk prep -- --template example.tmpl < example.json
~~~

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

~~~shell
    pttk prep -- -s -t markdown < example.json
~~~

You can specify the input file using the "-i" option or
provide it as the first filename after "prep". These are
equivallent.

~~~shell
    pttk prep -- -s -t markdown < example.json
    pttk prep -i example.json -- i -s -t markdown
    pttk prep example.json -- i -s -t markdown
~~~

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

~~~shell
    pttk prep codemeta.json about.md \
        -- --template codemeta-md.tmpl
~~~

# SEE ALSO

- pttk website at https://rsdoiel.github.io/pttk
- source code is avialable from https://github.com/rsdoiel/pttk


