// include.go - an simple file include preprocessor based example in
// Software Tools in Pascal by Kernigan & Plauger
package include

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	includeText = `% {app_name}-{verb}(1) {app_name}-{verb} user manual
% R. S. Doiel
% August 26, 2022

# NAME

{app_name} {verb}

# SYNOPSIS

{app_name} {verb} [INPUT FILENAME] [OUTPUT_FILENAME]

# DESCRIPTION

For each line that starts with the include directive will cause
the included file to be written to the output stream. If the
included file itself has include directives those will be
rendering int he output stream.

The include directive starts with "#include(" and is closed
with ");". Between the start and end of the directive the
text is considered a filename. Any text in the line with the
directive after ");" will not be included in the output stream.

{app_name} normally reads from standard input and writes to
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
    {app_name} {verb} myfile.txt
~~~

Read "myfile.txt" and write to "yourfile.txt"

~~~shell
    {app_name} {verb} myfile.txt yourfile.txt
~~~

To piece together a Markdown file from parts named in "book.txt"
which include the following Markdown documents.

~~~shell
    {app_name} {verb} book.txt book.md
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

The line "#include(toc.md);" tells the {verb} operation
to include the "toc.md" file. Similary the lines for "chapter1.md"
and "chapter2.md" work the say way.  If "toc.md", "chapter1.md"
or "chapter2.md" also has include directives there content would
also be included in the output stream.

`

	prefix = `#include(`
	suffix = `);`
)

func fmtText(s string, appName string, verb string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "{app_name}", appName), "{verb}", verb)
}

func hasIncludeName(s string) bool {
	if strings.HasPrefix(s, prefix) && strings.Contains(s, suffix) {
		return true
	}
	return false
}

func getIncludeName(s string) (string, error) {
	out := strings.TrimPrefix(s, prefix)
	if strings.Compare(out, s) == 0 {
		return "", fmt.Errorf("could not find %q in %q", prefix, s)
	}
	pos := strings.Index(out, suffix)
	if pos < 0 {
		return "", fmt.Errorf("could not find %q in %q", suffix, s)
	}
	return out[0:pos], nil
}

func Apply(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		// Check to see if we have prefix and suffix
		if hasIncludeName(line) {
			fName, err := getIncludeName(line)
			if err != nil {
				return err
			}
			// If we have an include, then include the file by openning
			// it and recursively calling Apply before closing the file.
			includeFile, err := os.Open(fName)
			if err != nil {
				return err
			}
			if err := Apply(includeFile, out); err != nil {
				return err
			}
			includeFile.Close()
		} else {
			fmt.Fprintf(out, "%s\n", line)
		}
	}
	return nil
}

func RunInclude(appName string, verb string, args []string) error {
	var (
		showHelp bool
		input    string
		output   string
		err      error
	)
	flagSet := flag.NewFlagSet(appName+"-"+verb, flag.ExitOnError)
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&input, "i", "", "read input from filename")
	flagSet.StringVar(&output, "o", "", "write output to filename")
	flagSet.Parse(args)
	args = flagSet.Args()

	if showHelp {
		fmt.Fprintln(os.Stdout, fmtText(includeText, appName, verb))
		return nil
	}

	if input == "" && len(args) > 0 {
		input = args[0]
	}
	if output == "" && len(args) > 1 {
		output = args[1]
	}

	in := os.Stdin
	out := os.Stdout

	if input != "" && input != "-" {
		in, err = os.Open(input)
		if err != nil {
			return err
		}
		defer in.Close()
	}
	if output != "" && output != "-" {
		out, err = os.Create(output)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return Apply(in, out)
}
