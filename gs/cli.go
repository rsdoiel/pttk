// gs.go is a sub-package pdtk. This package is for testing a static content
// Gopher site.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
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
