package prep

const (
	helpText = `

NAME

   pdtk prep

SYNOPSIS

   pdtk [OPTIONS] prep [VERB_OPTIONS] [-- [PANDOC_OPTIONS] ... ]

DESCRIPTION

pdtk prep is a Pandoc preprocessor. It can read JSON 
or YAML from standard input and passes that via an internal 
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" marker. Options before "--" are for
the pdtk preprossor. 

OPTIONS

-help
: display usage

-i
: read from a file instead of standard input

-o
: write to a file instead of standard output

EXAMPLES

In this example we have a JSON object document called
"example.json" and a Pandoc template called "example.tmpl".
A redirect "<" is used to pipe the content of "example.json"
into the command line tool pdtk.

    pdtk prep -- --template example.tmpl < example.json

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

    pdtk prep -- -s -t markdown < example.json

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

    pdtk prep -i codemeta.json -o about.md \
        -- --template codemeta-md.tmpl

Using pdtk to manage blog content with the "blogit"
verb. 

SEE ALSO

pdtk website at https://rsdoiel.github.io/pdtk

The source code is avialable from https://github.com/rsdoiel/pdtk


`
)
