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
)

const (
	helpText = `
USAGE

    {app_name} {verb} [OPTIONS]

SYNOPSIS

{app_name} {verb} provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.

EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2021.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
use the following command.

   cd Sites/me.example.org
   {app_name} {verb} my-vacation-day.md 2021-07-01

The *{app_name} {verb}* command will copy "my-vacation-day.md",
creating any necessary file directories to 
"Sites/me.example.org/2021/06/01".  It will also update article 
lists (index.md) at the year level, month, and day level and month
level of the directory tree and and generate/update a posts.json
in the "Sites/my.example.org" that can be used in your home page
template for listing recent posts.

*{app_name} {verb}* includes an option to set the prefix path to
the blog posting.  In this way you could have separate blogs 
structures for things like podcasts or videocasts.

    # Add a landing page for the podcast
    {app_name} {verb} -prefix=podcast my-vacation.md 2021-07-01
    # Add an audio file containing the podcast
    {app_name} {verb} -prefix=podcast my-vacation.wav 2021-07-01

Where "-p, -prefix" sets the prefix path before the YYYY/MM/DD path.


If you have an existing blog paths in the form of
PREFIX/YYYY/MM/DD you can use blogit to create/update/recreate
the blog.json file.

    {app_name} {verb} -prefix=blog -refresh=2021

The option "-refresh" is what indicates you want to crawl
for blog posts for that year.
`
)

var (
	// Standard Options
	showHelp    bool
	showVerbose bool
	quiet       bool

	// Application Options
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
	return strings.ReplaceAll(strings.ReplaceAll(helpText, "{app_name}", appName), "{verb}", verb)
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
