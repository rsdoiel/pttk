// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package prep

const (
	helpText = `% {app_name}-{verb}(1) {app_name}-{verb} user manual
% R. S. Doiel
% August 18, 2022

# NAME

{app_name} {verb}

# SYNOPSIS

{app_name} {verb} [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME] [-- [PANDOC_OPTIONS] ... ]

# DESCRIPTION

{app_name} {verb} is a Pandoc preprocessor. It can read JSON
or YAML from standard input and passes that via an internal
pipe to Pandoc as YAML front matter. Pandoc can then process it
accordingly Pandoc options. Pandoc options are those options
coming after a "--" (double dash) marker. Options before "--"
are for the {app_name} preprossor.

NOTE: prep requires that "pandoc-server" be running in order to work.

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

- pttk website at https://rsdoiel.github.io/pttk
- source code is avialable from https://github.com/rsdoiel/pttk

`
)
