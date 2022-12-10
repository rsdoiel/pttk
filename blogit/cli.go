// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package blogit

import (
	"encoding/json"
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
	blogAsset      bool
	refreshBlog    string
	setName        string
	setStarted     string
	setEnded       string
	setQuote       string
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

type BlogitConfig struct {
	Author        string `json:"author,omitempty"`
	SaveAsYaml    bool   `json:"save_as_yaml,omitempty"`
	Name          string `json:"name,omitempty"`
	Quote         string `json:"quote,omitempty"`
	PrefixPath    string `json:"prefix_path,omitempty"`
	Copyright     string `json:"copyright,omitempty"`
	Language      string `json:"language,omitempty"`
	License       string `json:"license,omitempty"`
	Started       string `json:"started,omitempty"`
	Ended         string `json:"ended,omitempty"`
	Description   string `json:"description,omitempty"`
	URL           string `json:"url,omitempty"`
	IndexTemplate string `json:"index_template,omitempty"`
	PostTemplate  string `json:"post_template,omitempty"`
}

func RunBlogIt(appName string, verb string, vargs []string) error {
	cfg := new(BlogitConfig)
	// read in .pttk configuration files if they exist, the setup defaults
	if _, err := os.Stat(".blogit"); err == nil {
		src, err := os.ReadFile(".blogit")
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		if err := json.Unmarshal(src, &cfg); err != nil {
			return err
		}
	}
	if cfg.Language == "" {
		cfg.Language = "en-US"
	}
	flagSet := flag.NewFlagSet(appName+":"+verb, flag.ExitOnError)

	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showVerbose, "verbose", false, "verbose output")

	// Application specific options
	flagSet.StringVar(&author, "author", cfg.Author, `Set the author name for use with "Simple Timesheet Notation" file for blog posts`)
	flagSet.StringVar(&stnImport, "stn", "", `Use a "Simple Timesheet Notation" file for blog posts`)
	flagSet.BoolVar(&saveAsYAML, "save-as-yaml", cfg.SaveAsYaml, "save as YAML file instead of blog.yaml file")
	flagSet.StringVar(&prefixPath, "prefix", cfg.PrefixPath, "Set the prefix path before YYYY/MM/DD.")
	flagSet.StringVar(&refreshBlog, "refresh", "", "Refresh blog.json for a given year")
	flagSet.StringVar(&setName, "name", cfg.Name, "Set the blog name.")
	flagSet.StringVar(&setQuote, "quote", cfg.Quote, "Set the blog quote.")
	flagSet.StringVar(&setCopyright, "copyright", cfg.Copyright, "Set the blog copyright notice.")
	flagSet.StringVar(&setLanguage, "language", cfg.Language, "Set the blog language.")
	flagSet.StringVar(&setLicense, "license", cfg.License, "Set the blog language license.")
	flagSet.StringVar(&setStarted, "started", cfg.Started, "Set the blog started date.")
	flagSet.StringVar(&setStarted, "ended", cfg.Ended, "Set the blog ended date.")
	flagSet.StringVar(&setDescription, "description", cfg.Description, "Set the blog description")
	flagSet.StringVar(&setBaseURL, "url", cfg.URL, "Set blog's URL")
	flagSet.StringVar(&setIndexTmpl, "index-tmpl", cfg.IndexTemplate, "Set index blog template")
	flagSet.StringVar(&setPostTmpl, "post-tmpl", cfg.PostTemplate, "Set index blog template")
	flagSet.BoolVar(&blogAsset, "asset", false, "Copy asset file to the blog path for provided date (YYYY-MM-DD)")

	flagSet.Parse(vargs)
	args := flagSet.Args()

	// Setup IO
	if showHelp {
		usage(appName, verb, 0)
	}
	if showVerbose {
		quiet = false
	}

	// Make ready to run one of the BlogIt command forms
	meta := new(BlogMeta)

	blogMetadataName := path.Join(prefixPath, "blog.json")

	// See if we have metadata to read in.
	if _, err := os.Stat(blogMetadataName); os.IsNotExist(err) {
		blogMetadataName = path.Join(prefixPath, "blog.yaml")
		if _, err := os.Stat(blogMetadataName); err != nil {
			blogMetadataName = ""
		}

	}

	if blogMetadataName != "" {
		if err := LoadBlogMeta(blogMetadataName, meta); err != nil {
			return fmt.Errorf("Error reading %q, %s\n", blogMetadataName, err)
		}
	}

	// handle option cases
	if saveAsYAML {
		blogMetadataName = path.Join(prefixPath, "blog.yaml")
	} else {
		blogMetadataName = path.Join(prefixPath, "blog.json")
	}
	if setName != "" {
		meta.Name = setName
	}
	if setQuote != "" {
		meta.Quip = setQuote
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

	// Handle Import of STN for blog posts
	if stnImport != "" {
		if err := meta.BlogSTN(prefixPath, stnImport, author); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		if err := meta.Save(blogMetadataName); err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
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
			fmt.Printf("Refreshing (%d/%d) %q from %q\n", i+1, len(years), blogMetadataName, path.Join(prefixPath, year))
			if err := meta.RefreshFromPath(prefixPath, year); err != nil {
				return fmt.Errorf("%s\n", err)
			}
		}
		if err := meta.Save(blogMetadataName); err != nil {
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
		if setName != "" || setQuote != "" || setDescription != "" ||
			setBaseURL != "" || setIndexTmpl != "" || setPostTmpl != "" {
			if err := meta.Save(blogMetadataName); err != nil {
				return fmt.Errorf("%s\n", err)
			}
			fmt.Printf("Updated %q completed.\n", blogMetadataName)
			return nil
		}
		usage(appName, verb, 1)
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
	if err := meta.Save(blogMetadataName); err != nil {
		return fmt.Errorf("%s\n", err)
	}
	return nil
}
