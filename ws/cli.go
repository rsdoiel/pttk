// The ws package provides static web services for pdtk.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package ws

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	showHelp bool
	// local app options
	uri          string
	docRoot      string
	sslKey       string
	sslCert      string
	CORSOrigin   string
	redirectsCSV string
)

func helpPage(appName string, verb string, text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "{verb}", verb), "{app_name}", appName)
}

func exitOnError(w io.Writer, err error, exitCode int) {
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		os.Exit(exitCode)
	}
}

func RunWS(appName string, verb string, vargs []string) error {
	flagSet := flag.NewFlagSet(verb, flag.ExitOnError)
	flagSet.BoolVar(&showHelp, "help", false, fmt.Sprintf("display help for %s", verb))
	flagSet.StringVar(&docRoot, "htdoc", ".", "set the document root")
	flagSet.StringVar(&sslKey, "ssl-key", "", "set the path to the SSL key")
	flagSet.StringVar(&sslCert, "ssl-cert", "", "set the path to the SSL cert")
	flagSet.StringVar(&CORSOrigin, "cors", "*.*", "set the path for CORS Origin")
	flagSet.StringVar(&uri, "url", "http://localhost:8000", "set the URL to listen on")
	flagSet.Parse(vargs)

	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", helpPage(appName, verb, helpText))
		os.Exit(0)
	}

	log.Printf("DocRoot %s", docRoot)

	u, err := url.Parse(uri)
	exitOnError(eout, err, 1)

	if u.Scheme == "https" {
		log.Printf("SSL Key %s", sslKey)
		log.Printf("SSL Cert %s", sslCert)
	}
	log.Printf("Listening for %s", uri)
	cors := CORSPolicy{
		Origin: CORSOrigin,
	}
	// Setup redirects defined the redirects CSV
	var rService *RedirectService
	if redirectsCSV != "" {
		src, err := ioutil.ReadFile(redirectsCSV)
		exitOnError(eout, fmt.Errorf("Can't read %s, %s", redirectsCSV, err), 1)
		r := csv.NewReader(bytes.NewReader(src))
		// Allow support for comment rows
		r.Comment = '#'
		// Make a redirect map[string]string
		rmap := map[string]string{}
		for {
			row, err := r.Read()
			if err == io.EOF {
				break
			}
			exitOnError(eout, fmt.Errorf("Can't read %s, %s", redirectsCSV, err), 1)

			if len(row) == 2 {
				// Define direct here.
				target := ""
				destination := ""
				if strings.HasPrefix(row[0], "/") {
					target = row[0]
				} else {
					target = "/" + row[0]
				}
				if strings.HasPrefix(row[1], "/") {
					destination = row[1]
				} else {
					destination = "/" + row[1]
				}
				rmap[target] = destination
			}
		}
		rService, err = MakeRedirectService(rmap)
		exitOnError(eout, fmt.Errorf("Can't make redirect service, %s", err), 1)
	}
	http.Handle("/", cors.Handler(http.FileServer(http.Dir(docRoot))))
	if u.Scheme == "https" {
		if rService != nil {
			err = http.ListenAndServeTLS(u.Host, sslCert, sslKey, RequestLogger(StaticRouter(rService.RedirectRouter(http.DefaultServeMux))))
		} else {
			err = http.ListenAndServeTLS(u.Host, sslCert, sslKey, RequestLogger(StaticRouter(http.DefaultServeMux)))
		}
		exitOnError(eout, err, 1)
	} else {
		if rService != nil {
			err = http.ListenAndServe(u.Host, RequestLogger(StaticRouter(rService.RedirectRouter(http.DefaultServeMux))))
		} else {
			err = http.ListenAndServe(u.Host, RequestLogger(StaticRouter(http.DefaultServeMux)))
		}
		exitOnError(eout, err, 1)
	}
	return nil
}
