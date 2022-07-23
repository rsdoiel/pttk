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
	"io/ioutil"
	"os"
)

func RunPrep(appName string, verb string, args []string) ([]byte, error) {
	var (
		input  string
		output string
		err    error
	)
	flagSet := flag.NewFlagSet(appName, flag.ExitOnError)
	flagSet.StringVar(&input, "i", "", "read JSON or YAML from file")
	flagSet.StringVar(&output, "o", "", "write Pandoc output to file")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)

	args = flagSet.Args()
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
