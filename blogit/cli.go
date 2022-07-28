// blogit.go is a sub-package pdtk. A packages for managing static content
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package blogit

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rsdoiel/pdtk/help"
)

var (
	// Standard Options
	showHelp    bool
	showVerbose bool
	quiet       bool

	// Application Options
	channelFile        string
	channelDescription string
	channelCopyright   string

	prefixPath     string
	docName        string
	dateString     string
	blogAsset      bool
	refreshBlog    string
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

func usage(appName string, verb string) string {
	return help.Render(appName, verb, helpText)
}

func RunBlogIt(appName string, verb string, vargs []string) error {
	flagSet := flag.NewFlagSet(appName+":"+verb, flag.ExitOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showVerbose, "verbose", false, "verbose output")

	// Application specific options
	flagSet.StringVar(&prefixPath, "prefix", "", "Set the prefix path before YYYY/MM/DD.")
	flagSet.StringVar(&refreshBlog, "refresh", "", "Refresh blog.json for a given year")
	flagSet.StringVar(&setName, "name", "", "Set the blog name.")
	flagSet.StringVar(&setQuip, "quip", "", "Set the blog quip.")
	flagSet.StringVar(&setCopyright, "copyright", "", "Set the blog copyright notice.")
	flagSet.StringVar(&setLanguage, "language", "en-US", "Set the blog language.")
	flagSet.StringVar(&setLicense, "license", "", "Set the blog language license.")
	flagSet.StringVar(&setStarted, "started", "", "Set the blog started date.")
	flagSet.StringVar(&setStarted, "ended", "", "Set the blog ended date.")
	flagSet.StringVar(&setDescription, "description", "", "Set the blog description")
	flagSet.StringVar(&setBaseURL, "url", "", "Set blog's URL")
	flagSet.StringVar(&setIndexTmpl, "index-tmpl", "", "Set index blog template")
	flagSet.StringVar(&setPostTmpl, "post-tmpl", "", "Set index blog template")
	flagSet.BoolVar(&blogAsset, "asset", false, "Copy asset file to the blog path for provided date (YYYY-MM-DD)")

	flagSet.Parse(vargs)
	args := flagSet.Args()

	// Setup IO
	if showHelp {
		fmt.Fprintf(os.Stdout, "%s\n", usage(appName, verb))
		return nil
	}
	if showVerbose {
		quiet = false
	}

	// Make ready to run one of the BlogIt command forms
	meta := new(BlogMeta)

	blogJSON := path.Join(prefixPath, "blog.json")

	// See if we have data to read in.
	if _, err := os.Stat(blogJSON); os.IsNotExist(err) {
	} else {
		if err := LoadBlogMeta(blogJSON, meta); err != nil {
			return fmt.Errorf("Error reading %q, %s\n", blogJSON, err)
		}
	}

	// handle option cases
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

	// handle option terminating case of refreshBlog
	if refreshBlog != "" {
		years := []string{}
		if strings.Contains(refreshBlog, ",") {
			years = strings.Split(refreshBlog, ",")
		} else {
			years = []string{refreshBlog}
		}
		for i, year := range years {
			year = strings.TrimSpace(year)
			fmt.Printf("Refreshing (%d/%d) %q from %q\n", i+1, len(years), blogJSON, path.Join(prefixPath, year))
			if err := meta.RefreshFromPath(prefixPath, year); err != nil {
				return fmt.Errorf("%s\n", err)
			}
		}
		if err := meta.Save(blogJSON); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		fmt.Printf("Refresh completed.\n")
		return nil
	}

	// We have a standard BlogIt command, process args.
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
			if err := meta.Save(blogJSON); err != nil {
				return fmt.Errorf("%s\n", err)
			}
			fmt.Printf("Updated blog.json completed.\n")
			return nil
		}
		return fmt.Errorf("%s\n", usage(appName, verb))
	}
	// Handle Copy Asset terminating case
	if blogAsset {
		fmt.Printf("Adding asset %q to posts for %q\n", docName, dateString)
		if err := meta.BlogAsset(prefixPath, docName, dateString); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	}

	// Now blog it.
	if err := meta.BlogIt(prefixPath, docName, dateString); err != nil {
		return fmt.Errorf("%s\n", err)
	}
	if err := meta.Save(blogJSON); err != nil {
		return fmt.Errorf("%s\n", err)
	}
	return nil
}
