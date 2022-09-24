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

func RunPrep(appName string, verb string, vargs []string) error {
	var (
		showHelp bool
		input    string
		output   string
		err      error
	)
	// Copy out pandoc options after "--"
	options := []string{}
	if len(vargs) > 0 {
		copyOpt := false
		cutPos := -1
		for i, arg := range vargs {
			if copyOpt {
				options = append(options, arg)
			}
			if arg == "--" {
				copyOpt = true
				cutPos = i
			}
		}
		if copyOpt && cutPos > -1 {
			vargs = vargs[0:cutPos]
		}
	}

	flagSet := flag.NewFlagSet(appName, flag.ExitOnError)
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&input, "i", "", "read JSON or YAML from file")
	flagSet.StringVar(&output, "o", "", "write Pandoc output to file")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")

	flagSet.Parse(vargs)
	args := flagSet.Args()
	if showHelp {
		usage(appName, verb, 0)
	}
	SetVerbose(verbose)

	//fmt.Printf("DEBUG vargs -> %+v\n", vargs)
	//fmt.Printf("DEBUG args -> %+v\n", args)

	if input == "" && len(args) > 0 && !strings.HasPrefix(args[0], "--") {
		input = args[0]
		args = args[1:]
	}
	if output == "" && len(args) > 0 && !strings.HasPrefix(args[0], "--") {
		output = args[0]
		args = args[1:]
	}
	if len(args) > 0 {
		return fmt.Errorf("did not expect remaining args: %+v\n", args)
	}

	//fmt.Printf("DEBUG input %q, output %q\n", input, output)
	//fmt.Printf("DEBUG options %+v\n", options)

	// The default action is to just processing JSON/YAML
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
	//return fmt.Errorf("DEBUG test option processing")
	return ApplyIO(in, out, options)
}
