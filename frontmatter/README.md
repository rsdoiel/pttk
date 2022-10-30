
frontmatter
===========

This go package will read in a plain text file and look for a starting delimiter of "---"
in the first to identify the start of a YAML frontmatter block like those supported by
[Pandoc](https://pandoc.org).  If fount the YAML text will be read in and converted
to JSON and return an a byte slice.


