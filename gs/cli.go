// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package gs

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

var (
	showHelp bool
	// local app options
	uri     string
	docRoot string
)

func usage(appName string, verb string, exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "%s\n", strings.ReplaceAll(strings.ReplaceAll(helpText, "{verb}", verb), "{app_name}", appName))
	os.Exit(exitCode)
}

func exitOnError(w io.Writer, err error, exitCode int) {
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		os.Exit(exitCode)
	}
}

func RunGS(appName string, verb string, vargs []string) error {
	var (
		u   *url.URL
		err error
	)
	flagSet := flag.NewFlagSet(verb, flag.ExitOnError)
	flagSet.BoolVar(&showHelp, "help", false, fmt.Sprintf("display help for %s", verb))
	flagSet.StringVar(&docRoot, "htdoc", "", "set the document root")
	flagSet.StringVar(&uri, "url", "", "set the URL to listen on")
	flagSet.Parse(vargs)
	args := flagSet.Args()

	eout := os.Stderr

	if showHelp {
		usage(appName, verb, 0)
	}

	if uri != "" {
		u, err = url.Parse(uri)
		exitOnError(eout, err, 1)
	}
	for _, arg := range args {
		switch {
		case uri == "" && strings.HasPrefix(arg, "gopher://"):
			u, err = url.Parse(arg)
			exitOnError(eout, err, 1)
		case docRoot == "":
			docRoot = arg
		}
	}
	gs := DefaultGopherService()
	if docRoot != "" {
		if _, err := os.Stat(docRoot); err != nil {
			exitOnError(eout, err, 1)
		}
		gs.DocRoot = docRoot
	}
	if u != nil {
		if u.Scheme != "gopher" {
			return fmt.Errorf("%q not supported by gopher service", u.Scheme)
		}
		gs.Gopher.Scheme = u.Scheme
		gs.Gopher.Host = u.Hostname()
		gs.Gopher.Port = u.Port()
	}
	return gs.Run()
}
