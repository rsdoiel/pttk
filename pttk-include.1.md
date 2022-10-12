% pttk-include(1) pttk-include user manual
% R. S. Doiel
% August 26, 2022

# NAME

pttk include

# SYNOPSIS

pttk include [INPUT FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

For each line that starts with the include directive will cause
the included file to be written to the output stream. If the
included file itself has include directives those will be
rendering int he output stream.

The include directive starts with "#include(" and is closed
with ");". Between the start and end of the directive the
text is considered a filename. Any text in the line with the
directive after ");" will not be included in the output stream.

pttk normally reads from standard input and writes to
standard output but may read from an optional input file name
and write to an optional output filename. The input filename is
presumed to come before the output filename.

# OPTIONS

-i FILENAME
: Using FILENAME for input stream

-o FILENAME
: Write output stream to FILENAME

-help
: Display this help

# EXAMPLES

## running the command

In this example we run an "include" operation on the file
"myfile.txt" and each of the include directives encountered
will be processed recursively.

~~~shell
    pttk include myfile.txt
~~~

Read "myfile.txt" and write to "yourfile.txt"

~~~shell
    pttk include myfile.txt yourfile.txt
~~~

To piece together a Markdown file from parts named in "book.txt"
which include the following Markdown documents.

~~~shell
    pttk include book.txt book.md
~~~

## exampel of directives

The "book.md" could look like

~~~markdown
    ---
    title: "Mybook"
    author: "jane.doe@example.org (Jane Doe)"
    pubDate: 2022-08-26
    ---

    Mybook
    ======

    #include(toc.md);

    #include(chapter1.md);

    #include(chapter2.md);
~~~

The line "#include(toc.md);" tells the include operation
to include the "toc.md" file. Similary the lines for "chapter1.md"
and "chapter2.md" work the say way.  If "toc.md", "chapter1.md"
or "chapter2.md" also has include directives there content would
also be included in the output stream.


