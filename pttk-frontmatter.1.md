%% pttk-frontmatter(1) user manual
% R. S. Doiel
% 2022-10-30

# NAME

pttk frontmatter

# SYNOPSIS

pttk frontmatter [OPTIONS] [INPUT_FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

The frontmatter action allows you to extract metadata from a
text document that uses the Markdown style frontmatter. By default
it will return the documents frontmatter as a JSON structure.
If you include a root level attribute name then you can extract
the value as a unquoted string. By default frontmatter can read
from standard input but if you provide a filename it'll ready
from the file instead. Likewise it writes to standard out but
if an output filename is provided it write to that file instead.

# OPTIONS

-help
: display this help page

# EXAMPLES

Here is an example of extracting a frontmatter from "mypost.md"
markdown file.

~~~
pttk frontmatter mypost.md
~~~

Here is an example of extracting a "title" from the frontmatter
when combined with Stephen Dolan's jq.

~~~
pttk frontmatter mypost.md | jq .title
~~~



