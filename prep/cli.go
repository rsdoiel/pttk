// prep.go is a sub-package pdtk. A packages for managing static content
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package prep

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func usage(appName string, verb string, exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "%s\n", strings.ReplaceAll(strings.ReplaceAll(helpText, "{app_name}", appName), "{verb}", verb))
	os.Exit(exitCode)
}

func RunPrep(appName string, verb string, args []string) ([]byte, error) {
	var (
		showHelp bool
		input    string
		output   string
		err      error
	)
	flagSet := flag.NewFlagSet(appName, flag.ExitOnError)
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&input, "i", "", "read JSON or YAML from file")
	flagSet.StringVar(&output, "o", "", "write Pandoc output to file")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)

	args = flagSet.Args()
	if showHelp {
		usage(appName, verb, 0)
	}
	SetVerbose(verbose)
	// The default action is to just processing JSON/YAML
	in := os.Stdin
	out := os.Stderr
	if input != "" && input != "-" {
		in, err = os.Open(input)
		if err != nil {
			return nil, err
		}
		defer in.Close()
	}
	if output != "" && output != "-" {
		out, err = os.Create(output)
		if err != nil {
			return nil, err
		}
		defer out.Close()
	}
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	return Apply(buf, args)
}
