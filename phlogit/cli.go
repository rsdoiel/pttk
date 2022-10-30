// phlogit.go is a sub-package pttk. A packages for managing static content
// phlogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package phlogit

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rsdoiel/pttk/help"
)

var (
	// Standard Options
	showHelp    bool
	showVerbose bool
	quiet       bool

	// Application Options
	saveAsYAML bool

	stnImport      string
	author         string
	prefixPath     string
	docName        string
	dateString     string
	phlogAsset     bool
	refreshPhlog   string
	setName        string
	setStarted     string
	setEnded       string
	setQuip        string
	setDescription string
	setBaseURL     string
	setIndexTmpl   string
	setPostTmpl    string
	setCopyright   string
	setLicense     string
	setLanguage    string
)

func usage(appName string, verb string, exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "%s\n", help.Render(appName, verb, helpText))
	os.Exit(exitCode)
}

func RunPhlogIt(appName string, verb string, vargs []string) error {
	flagSet := flag.NewFlagSet(appName+":"+verb, flag.ExitOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showVerbose, "verbose", false, "verbose output")

	// Application specific options
	flagSet.StringVar(&author, "author", "", `Set the author name for use with "Simple Timesheet Notation" file for phlog posts`)
	flagSet.StringVar(&stnImport, "stn", "", `Use a "Simple Timesheet Notation" file for phlog posts`)
	flagSet.BoolVar(&saveAsYAML, "save-as-yaml", false, "save as YAML file instead of phlog.yaml file")
	flagSet.StringVar(&prefixPath, "prefix", "", "Set the prefix path before YYYY/MM/DD.")
	flagSet.StringVar(&refreshPhlog, "refresh", "", "Refresh phlog.json for a given year")
	flagSet.StringVar(&setName, "name", "", "Set the phlog name.")
	flagSet.StringVar(&setQuip, "quip", "", "Set the phlog quip.")
	flagSet.StringVar(&setCopyright, "copyright", "", "Set the phlog copyright notice.")
	flagSet.StringVar(&setLanguage, "language", "en-US", "Set the phlog language.")
	flagSet.StringVar(&setLicense, "license", "", "Set the phlog language license.")
	flagSet.StringVar(&setStarted, "started", "", "Set the phlog started date.")
	flagSet.StringVar(&setStarted, "ended", "", "Set the phlog ended date.")
	flagSet.StringVar(&setDescription, "description", "", "Set the phlog description")
	flagSet.StringVar(&setBaseURL, "url", "", "Set phlog's URL")
	flagSet.StringVar(&setIndexTmpl, "index-tmpl", "", "Set index phlog template")
	flagSet.StringVar(&setPostTmpl, "post-tmpl", "", "Set index phlog template")
	flagSet.BoolVar(&phlogAsset, "asset", false, "Copy asset file to the phlog path for provided date (YYYY-MM-DD)")

	flagSet.Parse(vargs)
	args := flagSet.Args()

	// Setup IO
	if showHelp {
		usage(appName, verb, 0)
	}
	if showVerbose {
		quiet = false
	}

	// Make ready to run one of the PhlogIt command forms
	meta := new(PhlogMeta)

	phlogMetadataName := path.Join(prefixPath, "phlog.json")
if _, err := os.Stat(phlogMetadataName); err != nil {
	if _, err := os.Stat(phlogMetadataName); err == nil {
		phlogMetadataName = path.Join(prefixPath, "phlog.yaml")
	} else {
	phlogMetadataName = ""

	}
}
	if phlogMetadataName != "" {
		if err := LoadPhlogMeta(phlogMetadataName, meta); err != nil {
			return fmt.Errorf("Error reading %q, %s\n", phlogMetadataName, err)
		}
	}

	// handle option cases
	if saveAsYAML {
		phlogMetadataName = path.Join(prefixPath, "phlog.yaml")
	}
	if setName != "" {
		meta.Name = setName
	}
	if setQuip != "" {
		meta.Quip = setQuip
	}
	if setDescription != "" {
		meta.Description = setDescription
	}
	if setCopyright != "" {
		meta.Copyright = setCopyright
	}
	if setLicense != "" {
		meta.License = setLicense
	}
	if setStarted != "" {
		meta.Started = setStarted
	}
	if setEnded != "" {
		meta.Ended = setEnded
	}
	if setBaseURL != "" {
		meta.BaseURL = setBaseURL
	}
	if setIndexTmpl != "" {
		meta.IndexTmpl = setIndexTmpl
	}
	if setPostTmpl != "" {
		meta.PostTmpl = setPostTmpl
	}

	// Handle Import of STN for phlog posts
	if stnImport != "" {
		if err := meta.PhlogSTN(prefixPath, stnImport, author); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		if err := meta.Save(phlogMetadataName); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	}

	// handle option terminating case of refreshPhlog
	if refreshPhlog != "" {
		years := []string{}
		if strings.Contains(refreshPhlog, ",") {
			years = strings.Split(refreshPhlog, ",")
		} else {
			years = []string{refreshPhlog}
		}
		for i, year := range years {
			year = strings.TrimSpace(year)
			fmt.Printf("Refreshing (%d/%d) %q from %q\n", i+1, len(years), phlogMetadataName, path.Join(prefixPath, year))
			if err := meta.RefreshFromPath(prefixPath, year); err != nil {
				return fmt.Errorf("%s\n", err)
			}
		}
		if err := meta.Save(phlogMetadataName); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		fmt.Printf("Refresh completed.\n")
		return nil
	}

	// We have a standard PhlogIt command, process args.
	switch len(args) {
	case 1:
		docName, dateString = args[0], time.Now().Format(DateFmt)
	case 2:
		docName, dateString = args[0], args[1]
		if _, err := time.Parse(DateFmt, dateString); err != nil {
			return fmt.Errorf("Date error %q, %s", dateString, err)
		}
	default:
		if setName != "" || setQuip != "" || setDescription != "" ||
			setBaseURL != "" || setIndexTmpl != "" || setPostTmpl != "" {
			if err := meta.Save(phlogMetadataName); err != nil {
				return fmt.Errorf("%s\n", err)
			}
			fmt.Printf("Updated %q completed.\n", phlogMetadataName)
			return nil
		}
		usage(appName, verb, 1)
	}

	// Handle Copy Asset terminating case
	if phlogAsset {
		fmt.Printf("Adding asset %q to posts for %q\n", docName, dateString)
		if err := meta.PhlogAsset(prefixPath, docName, dateString); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	}

	// Now phlog it.
	if err := meta.PhlogIt(prefixPath, docName, dateString); err != nil {
		return fmt.Errorf("%s\n", err)
	}
	if err := meta.Save(phlogMetadataName); err != nil {
		return fmt.Errorf("%s\n", err)
	}
	return nil
}
