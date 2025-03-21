// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package phlogit

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
	phlogAsset     bool
	refreshPhlog   string
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
	setMasthead    string
)

func usage(appName string, verb string, helpText string, exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "%s\n", help.Render(appName, verb, helpText))
	os.Exit(exitCode)
}

type PhlogitConfig struct {
	Author        string `json:"author,omitempty" yaml:"author,omitempty"`
	SaveAsYaml    bool   `json:"save_as_yaml,omitempty" yaml:"save_as_yaml,omitempty"`
	// NOTE: This holds the text of the Masthead for the gophermap, the command line options points to a file holding this value
	Masthead      string `json:"masthead,omitempty" yaml:"masthead,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Quote         string `json:"quote,omitempty" yaml:"quote,omitempty"`
	PrefixPath    string `json:"prefix_path,omitempty" yaml:"prefix_path,omitempty"`
	Copyright     string `json:"copyright,omitempty" yaml:"copyright,omitempty"`
	Language      string `json:"language,omitempty" yaml:"language,omitempty"`
	License       string `json:"license,omitempty" yaml:"license,omitempty"`
	Started       string `json:"started,omitempty" yaml:"started,omitempty"`
	Ended         string `json:"ended,omitempty" yaml:"ended,omitempty"`
	Description   string `json:"description,omitempty" yaml:"description,omitempty"`
	URL           string `json:"url,omitempty" yaml:"url,omitempty"`
	IndexTemplate string `json:"index_template,omitempty" yaml:"index_template,omitempty"`
	PostTemplate  string `json:"post_template,omitempty" yaml:"post_template,omitempty"`
}

func RunGophermap(appName string, verb string, vargs []string) error {
	cfg := new(PhlogitConfig)
	// read in .pttk configuration files if they exist, the setup defaults
	if _, err := os.Stat(".phlogit"); err == nil {
		src, err := os.ReadFile(".phlogit")
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
	flagSet.StringVar(&setMasthead, "masthead", "", "Read in the Masthead from the filename provided")

	flagSet.Parse(vargs)
	args := flagSet.Args()

	// Setup IO
	if showHelp {
		usage(appName, verb, helpTextGophermap, 0)
	}
	if showVerbose {
		quiet = false
	}

	// Make ready to run one of the gophermap command forms
	meta := new(PhlogMeta)

	if setMasthead != "" {
		// Read in the Masthead and assign it to meta.Masthead
		src, err := os.ReadFile(setMasthead)
		if err != nil {
			return err
		}
		meta.Masthead = fmt.Sprintf("%s", src)
	}

	// We have a standard Gophermap command, process args.
	gophermapName, fNames := "", []string{}
	if len(args) > 0 {
		gophermapName = args[0]
	} else if len(args) > 1 {
		gophermapName, fNames = args[0], args[1:]
	} else {
		usage(appName, verb, helpTextGophermap, 1)
	}

	// Now Gophermap it.
	if err := meta.Gophermap(gophermapName, fNames); err != nil {
		return fmt.Errorf("%s\n", err)
	}
	return nil
}

func RunPhlogIt(appName string, verb string, vargs []string) error {
	cfg := new(PhlogitConfig)
	// read in .pttk configuration files if they exist, the setup defaults
	if _, err := os.Stat(".phlogit"); err == nil {
		src, err := os.ReadFile(".phlogit")
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
	flagSet.StringVar(&author, "author", cfg.Author, `Set the author name for use with "Simple Timesheet Notation" file for phlog posts`)
	flagSet.StringVar(&stnImport, "stn", "", `Use a "Simple Timesheet Notation" file for phlog posts`)
	flagSet.BoolVar(&saveAsYAML, "save-as-yaml", cfg.SaveAsYaml, "save as YAML file instead of phlog.yaml file")
	flagSet.StringVar(&prefixPath, "prefix", cfg.PrefixPath, "Set the prefix path before YYYY/MM/DD.")
	flagSet.StringVar(&refreshPhlog, "refresh", "", "Refresh phlog.json for a given year")
	flagSet.StringVar(&setName, "name", cfg.Name, "Set the phlog name.")
	flagSet.StringVar(&setQuote, "quote", cfg.Quote, "Set the phlog quote.")
	flagSet.StringVar(&setCopyright, "copyright", cfg.Copyright, "Set the phlog copyright notice.")
	flagSet.StringVar(&setLanguage, "language", cfg.Language, "Set the phlog language.")
	flagSet.StringVar(&setLicense, "license", cfg.License, "Set the phlog language license.")
	flagSet.StringVar(&setStarted, "started", cfg.Started, "Set the phlog started date.")
	flagSet.StringVar(&setStarted, "ended", cfg.Ended, "Set the phlog ended date.")
	flagSet.StringVar(&setDescription, "description", cfg.Description, "Set the phlog description")
	flagSet.StringVar(&setBaseURL, "url", cfg.URL, "Set phlog's URL")
	flagSet.StringVar(&setIndexTmpl, "index-tmpl", cfg.IndexTemplate, "Set index phlog template")
	flagSet.StringVar(&setPostTmpl, "post-tmpl", cfg.PostTemplate, "Set index phlog template")
	flagSet.BoolVar(&phlogAsset, "asset", false, "Copy asset file to the phlog path for provided date (YYYY-MM-DD)")
	flagSet.StringVar(&setMasthead, "masthead", "", "Read in the Masthead from the filename provided")

	flagSet.Parse(vargs)
	args := flagSet.Args()

	// Setup IO
	if showHelp {
		usage(appName, verb, helpTextPhlog, 0)
	}
	if showVerbose {
		quiet = false
	}

	// Make ready to run one of the PhlogIt command forms
	meta := new(PhlogMeta)

	// Figure out if we are working with phlog.json or phlog.yaml.
	phlogMetadataName := path.Join(prefixPath, "phlog.json")
	// Handle the case where we want to read in JSON but save as YAML.
	if saveAsYAML {
		phlogMetadataName = path.Join(prefixPath, "phlog.yaml")
	}
	loadMetadata := false
	if _, err := os.Stat(phlogMetadataName); err == nil {
		loadMetadata = true
	}
	if loadMetadata {
		if err := LoadPhlogMeta(phlogMetadataName, meta); err != nil {
			return fmt.Errorf("Error reading %q, %s\n", phlogMetadataName, err)
		}
	}

	// handle option cases
	if setName != "" {
		meta.Name = setName
	}
	if setMasthead != "" {
		// Read in the Masthead
		src, err := os.ReadFile(setMasthead)
		if err != nil {
			return fmt.Errorf("failed to read Masthead file %q", setMasthead)
		}
		meta.Masthead = fmt.Sprintf("%s", src)
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
		if setName != "" || setQuote != "" || setDescription != "" ||
			setBaseURL != "" || setIndexTmpl != "" || setPostTmpl != "" {
			if err := meta.Save(phlogMetadataName); err != nil {
				return fmt.Errorf("%s\n", err)
			}
			fmt.Printf("Updated %q completed.\n", phlogMetadataName)
			return nil
		}
		usage(appName, verb, helpTextPhlog, 1)
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
