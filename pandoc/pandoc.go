// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package pandoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type API struct {
	// Hostname is the host to contact, normally this is localhost unless
	// the Pandoc service is being proxied
	Hostname string `json:"hostname,omitempty"`
	// Port defaults to 3030, it is the port number that pandoc-server listens on
	Port string `json:"port,omitempty"`
	// Verbose, if true log success as well as errors
	Verbose bool
	// Settings holds the settings to use to submit to Pandoc Server's root end point
	Settings *Settings `json:"settings,omitempty"`
}

// Settings defines the data possible to send to Pandoc Server's
// root level end point
type Settings struct {
	// From is the doc type you are converting from, e.g. markdown
	From string `json:"from,omitempty"`
	// To is the doc type you are converting to, e.g. html5
	To string `json:"to,omitempty"`
	//
	// For the following fields see https://pandoc.org/pandoc-server.html#root-endpoint
	//
	ShiftHeadingLevel     int                    `json:"shift-heading-level-by,omitempty"`
	IdentedCodeClasses    []string               `json:"indented-code-classes,omitempty"`
	DefaultImageExtension string                 `json:"default-image-extension,omitempty"`
	Metadata              string                 `json:"metadata,omitempty"`
	TabStop               int                    `json:"tab-stop,omitempty"`
	TrackChanges          string                 `json:"track-changes,omitempty"`
	Abbreviations         []string               `json:"abbreviations,omitempty"`
	Standalone            bool                   `json:"standalone,omitempty"`
	Text                  string                 `json:"text,omitempty"`
	Template              string                 `json:"template,omitempty"`
	Variables             map[string]interface{} `json:"variables,omitempty"`
	DPI                   int                    `json:"dpi,omitemtpy"`
	Wrap                  string                 `json:"wrap,omitempty"`
	Columns               int                    `json:"columns,omitempty"`
	TableOfContents       bool                   `json:"table-of-contents,omitempty"`
	TOCDepth              int                    `json:"toc-depth,omitempty"`
	StripComments         bool                   `json:"strip-comments,omitempty"`
	HighlightStyle        string                 `json:"highlight-style,omitempty"`
	EmbedResources        string                 `json:"embed-resources,omitempty"`
	HTMLQTags             bool                   `json:"html-q-tags,omitempty"`
	Ascii                 bool                   `json:"ascii,omitempty"`
	ReferenceLinks        bool                   `json:"reference-links,omitempty"`
	ReferenceLocation     string                 `json:"reference-location,omitempty"`
	SetExtHeaders         string                 `json:"setext-headers,omitempty"`
	TopLevelDivision      string                 `json:"top-level-division,omitempty"`
	NumberSections        string                 `json:"number-sections,omitempty"`
	NumberOffset          []int                  `json:"number-offset,omitempty"`
	HTMLMathMethod        string                 `json:"html-math-method,omitempty"`
	Listings              bool                   `json:"listings,omitempty"`
	Incremental           bool                   `json:"incremental,omitempty"`
	SideLevel             int                    `json:"slide-level,omitempty"`
	SectionDivs           bool                   `json:"section-divs,omitempty"`
	EmailObfuscation      string                 `json:"email-obfuscation,omitempty"`
	IdentifierPrefix      string                 `json:"identifier-prefix,omitempty"`
	TitlePrefix           string                 `json:"title-prefix,omitempty"`
	ReferenceDoc          string                 `json:"reference-doc,omitempty"`
	EPubCoverImage        string                 `json:"epub-cover-image,omitempty"`
	EPubMetadata          string                 `json:"epub-metadata,omitempty"`
	EPubChapterLevel      int                    `json:"epub-chapter-level,omitempty"`
	EPubSubdirectory      string                 `json:"epub-subdirectory,omitempty"`
	EPubFonts             string                 `json:"epub-fonts,omitempty"`
	IpynbOutput           string                 `json:"ipynb-output,omitempty"`
	Citeproc              bool                   `json:"citeproc,omitempty"`
	Bibliography          []string               `json:"bibliography,omitempty"`
	Csl                   string                 `json:"csl,omitempty"`
	CiteMethod            string                 `json:"cite-method,omitempty"`
	Files                 []string               `json:files,omitempty"`

	// Verbose if set true then include logging on success as well as error
	Verbose bool

	// ExtTypes holds a mapping of extension to file type, e.d. ".html" to "html5"
	//ExtTypes map[string]string `json:"ext-types,omitempty"`
}

var (
	// DefaultExtTypes maps file extensions to document types. This allows the "to", "from"
	// Pandoc options to be set based on file extension. This can be overwritten by setting
	// `.ext_types` in the JSON configuraiton file.
	DefaultExtTypes = map[string]string{
		".md":   "markdown",
		".html": "html5",
	}
)

func inStringList(val string, list []string) bool {
	for _, expected := range list {
		if val == expected {
			return true
		}
	}
	return false
}

// Load will read a JSON file containing configuration
// and return an API struct and error. The API structure
// is used to interact with a running Pandoc Server.
//
// ```
//
//	api, err := pandoc.Load("config.json")
//	if err != nil {
//	    // ... handle errror
//	}
//	src, err := os.ReadFile("helloworld.md")
//	if err != nil {
//	    // ... handle errror
//	}
//	buf := bytes.NewReader(src)
//	src, err = api.Convert(buf)
//	if err != nil {
//	    // ... handle errror
//	}
//	err = os.WriteFile("helloworld.html", src, 0664)
//	if err != nil {
//	    // ... handle errror
//	}
//
// ```
func Load(fName string) (*API, error) {
	src, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	api := new(API)
	if err := json.Unmarshal(src, api); err != nil {
		return nil, err
	}
	if api.Port == "" {
		api.Port = ":3030"
	} else if !strings.HasPrefix(api.Port, ":") {
		api.Port = fmt.Sprintf(":%s", api.Port)
	}
	if api.Settings == nil {
		api.Settings = Settings{}
	}
	if api.Settings.Template != "" {
		// See if we have file path, if so read in the template
		if !strings.Contains(api.Settings.Template, "\n") {
			if err, _ := os.Stat(api.Settings.Template); err == nil {
				src, err := os.ReadFile(api.Settings.Template)
				if err == nil {
					api.Settings.Template = fmt.Sprint("%s", src)
				}
			}
		}
	}

	if !inStringList(api.Settings.TrackChanges, []string{"accept", "reject", "all", ""}) {
		return api, fmt.Errorf("tract-changes: %q is not supported", api.TrackChanges)
	}
	if !inStringList(api.Settings.Wrap, []string{"auto", "preserve", "none", ""}) {
		return api, fmt.Errorf("wrap: %q is not supported", api.Wrap)
	}
	if !inStringList(api.Settings.HighlightStyle, []string{"pygments", "kate", "monochrome", "breezeDark", "espresso", "zenburn", "haddock", "tango", ""}) {
		return api, fmt.Errorf("highlight-style: %q is not supported", api.HighlightStyle)
	}
	if !inStringList(api.Settings.ReferenceLocation, []string{"document", "section", "block", ""}) {
		return api, fmt.Errorf("wrap: %q is not supported", api.ReferenceLocation)
	}
	if !inStringList(api.Settings.TopLevelDivision, []string{"default", "part", "chapter", "section", ""}) {
		return api, fmt.Errorf("top-level-division: %q is not supported", api.TopLevelDivision)
	}
	if !inStringList(api.Settings.HTMLMathMethod, []string{"plain", "webtex", "gladtex", "mathml", "mathjax", "katex", ""}) {
		return api, fmt.Errorf("html-math-method: %q is not supported", api.HTMLMathMethod)
	}
	if !inStringList(api.Settings.EmailObfuscation, []string{"none", "references", "javascript", ""}) {
		return api, fmt.Errorf("email-obfuscation: %q is not supported", api.EmailObfuscation)
	}
	if !inStringList(api.Settings.IpynbOutput, []string{"best", "all", "none", ""}) {
		return api, fmt.Errorf("ipynb-output: %q is not supported", api.IpynbOutput)
	}
	if !inStringList(api.Settings.CiteMethod, []string{"citeproc", "natbib", "biblatex", ""}) {
		return api, fmt.Errorf("cite-method: %q is not supported", api.CiteMethod)
	}
	return api, nil
}

// sendToRootEndpoint takes content type and sends the request to the Pandoc Server
// Root end point based on the state of configuration struct used.
func (api *API) sendToRootEndpoint() ([]byte, error) {
	// NOTE: Pandoc Server API want JSON in POST not urlencoded form data
	if api.Settings == nil {
		return nil, fmt.Errorf("api.Settings not configured")
	}
	if api.Settings.Text == "" {
		return nil, fmt.Errorf("expected to have a source text to convert, %+v", api)
	}
	src, err := json.Marshal(api.Settings)
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		log.Printf("Nothing to convert")
		return nil, fmt.Errorf("nothing to convert")
	}
	// Setup out our JSON post request.
	u := fmt.Sprintf("http://localhost%s/", api.Port)
	body := bytes.NewReader(src)
	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s POST failed, %s", u, err)
		return nil, err
	}
	defer resp.Body.Close()
	// Process response
	src, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s POST read body failed, %s", u, err)
		return nil, err
	}
	if len(src) == 0 {
		log.Printf("zero bytes returned from Root Endpoint")
		return nil, fmt.Errorf("zero bytes returned by pandoc")
	}
	if api.Verbose {
		log.Printf("%d bytes returned successful from Root Endpoint", len(src))
	}
	return src, nil
}

// Pandoc a takes the configuration settings and sends a request
// to the Pandoc server with contents read from the io.Reader
// and returns a slice of bytes and error.
//
// ```
//
// // Setup our client configuration
// api := pandoc.API{}
// api.Port = ":3030"
// api.Hostname = "localhost"
//
//	api.Settings = Settings{
//	               Standalone: true,
//	               From: "markdown",
//	               To: "html5",
//	}
//
// src, err := os.ReadFile("htdocs/index.md")
// // ... handle error
// txt, err :=  api.Convert(bytes.NewReader(src))
//
//	if err := os.WriteFile("htdocs/index.html", src, 0664); err != nil {
//	    // ... handle error
//	}
//
// ```
func (api *API) Convert(input io.Reader) ([]byte, error) {
	var src []byte

	src, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}
	// NOTE: The source needs to already be converted to bytes, if necessary base64 encoded.
	api.Text = fmt.Sprintf("%s", src)
	defer func() {
		api.Text = ""
	}()
	src, err = api.sendToRootEndpoint()
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Walk takes a path and walks the directories converting the files that map
// to the From values in the configuration.
func (api *API) Walk(startPath string, fromExt string, toExt string) error {
	err := filepath.Walk(startPath,
		func(fName string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := path.Ext(fName)
				if ext == fromExt {
					toFName := strings.TrimSuffix(fName, ext) + toExt
					src, err := os.ReadFile(fName)
					if err != nil {
						log.Printf("%s", err)
						return err
					}
					txt, err := api.Convert(bytes.NewReader(src))
					if err != nil {
						log.Printf("%s", err)
						return err
					}
					err = os.WriteFile(toFName, txt, 0664)
					if err != nil {
						log.Printf("%s", err)
						return err
					}
					if api.Verbose {
						log.Printf("covert %q to %q\n", fName, toFName)
					}
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
