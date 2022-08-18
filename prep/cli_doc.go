package prep

const (
	helpText = `% {app_name}-{verb}(1) {app_name}-{verb} user manual
% R. S. Doiel
% August 18, 2022

# NAME

{app_name} {verb}

# SYNOPSIS

{app_name} {verb} [OPTIONS] [INPUT_FILENAME] [-- [PANDOC_OPTIONS] ... ]

DESCRIPTION

{app_name} {verb} is a Pandoc preprocessor. It can read JSON 
or YAML from standard input and passes that via an internal 
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" (double dash) marker. Options before "--" 
are for the {app_name} preprossor. 

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
into the command line tool {app_name}.

~~~shell
    {app_name} {verb} -- --template example.tmpl < example.json
~~~

Render example.json as Markdown document. We need to use
Pandoc's own options of "-s" (stand alone) and "-t" (to
tell Pandoc the output format)

~~~shell
    {app_name} {verb} -- -s -t markdown < example.json
~~~

You can specify the input file using the "-i" option or
provide it as the first filename after "{verb}". These are
equivallent.

~~~shell
    {app_name} {verb} -- -s -t markdown < example.json
    {app_name} {verb} -i example.json -- i -s -t markdown
    {app_name} {verb} example.json -- i -s -t markdown
~~~

Process a "codemeta.json" file with "codemeta-md.tmpl" to
produce an about page in Markdown via Pandocs template
processing (the "codemeta-md.tmpl" is a Pandoc template
marked up to produce Markdown output).

~~~shell
    {app_name} {verb} codemeta.json about.md \
        -- --template codemeta-md.tmpl
~~~

# SEE ALSO

- pdtk website at https://rsdoiel.github.io/pdtk
- source code is avialable from https://github.com/rsdoiel/pdtk

`
)
