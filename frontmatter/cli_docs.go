// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package frontmatter

var (
	helpText = `%% pttk-frontmatter(1) user manual
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


`
)
