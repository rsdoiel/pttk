// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
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
