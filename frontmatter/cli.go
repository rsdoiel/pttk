// frontmatter.go is a sub-package pttk used to extract YAML frontmatter from a Markdown document and return endcoded JSON.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package frontmatter

import (
	"flag"
	"fmt"
	"os"

	"github.com/rsdoiel/pttk/help"
)

var (
	// Standard Options
	showHelp bool

	// Verb options
	inFName  string
	outFName string
)

func usage(appName string, verb string, exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "%s\n", help.Render(appName, verb, helpText))
	os.Exit(exitCode)
}

func RunFrontmatter(appName string, verb string, vargs []string) error {
	flagSet := flag.NewFlagSet(appName+":"+verb, flag.ExitOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display help")

	flagSet.Parse(vargs)
	args := flagSet.Args()

	// Setup IO
	if showHelp {
		usage(appName, verb, 0)
	}
	in := os.Stdin
	out := os.Stdout

	// We have a standard PhlogIt command, process args.
	switch len(args) {
	case 1:
		inFName, outFName = args[0], ""
	case 2:
		inFName, outFName = args[0], args[1]
	default:
		inFName, outFName = "", ""
	}
	var (
		src []byte
		err error
	)

	if inFName == "" || inFName == "-" {
		src, err = ReadAll(in)
	} else {
		src, err = ReadFile(inFName)
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}
