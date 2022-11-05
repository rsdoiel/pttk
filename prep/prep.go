// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package prep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	// Project packages
	"github.com/rsdoiel/pttk/pandoc"

	// 3rd Party Packages
	"gopkg.in/yaml.v3"
)

var verbose bool

// SetVerbose when set true will show the Pandoc command
// envocation before running Pandoc to process the JSON document
// and template. Mainly useful for debugging.
func SetVerbose(onoff bool) {
	verbose = onoff
}

// ReadAll reads JSON from as a stream using an io.Reader.
// Buffers it. Then uses Apply and options return
// a slice of bytes and error value.
//
// ```shell
//
//	// Options passed to Pandoc
//	opt := []string{}
//	out, err := prep.ReadAll(os.Stdin, "page.tmpl", opt)
//	if err != nil {
//	   // ... handle error
//	}
//	fmt.Fprintf(os.Stdout, "%s\n", out)
//
// ```
func ReadAll(r io.Reader, options []string) ([]byte, error) {
	// Read the JSON input
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Apply(src, options)
}

// ReadFile reads a JSON or YAML document from a file then uses Apply
// and options returning a slice of bytes and error value.
func ReadFile(name string, options []string) ([]byte, error) {
	// Read the JSON or YAML file
	src, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return Apply(src, options)
}

// ApplyIO reads in JSON from an io.Reader, passes options
// such on to pandoc.
//
// ```
//
//	// Options passed to Pandoc
//	opt := []string{"-s", "--template", "example.tmpl"}
//	err := pdtmpl.ApplyIO(os.Stdin, os.Stdout, opt)
//	if err != nil {
//	   // ... handle error
//	}
//
// ```
func ApplyIO(r io.Reader, w io.Writer, options []string) error {
	buf, err := io.ReadAll(r)
	src, err := Apply(buf, options)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", src)
	return err
}

// Apply takes a byte array (like you could read from os.Stdin
// containing JSON or YAML. It creates a in memory representation
// of Pandoc YAML front matter and pipes it to Pandoc via Pandoc's
// standard input. Pandoc then renders the output. Apply's options
// are simply passed as is to Pandoc when invoked.
//
// ```
//
//	src, err := os.ReadFile("example.json")
//	if err != nil {
//	   // ... handle error
//	}
//	// Options passed to Pandoc
//	opt := []string{"--template", "example.tmpl"}
//	src, err := pdtmpl.Apply(src, opt)
//	if err != nil {
//	   // ... handle error
//	}
//	fmt.Printf("%s\n", src)
//
// ```
func Apply(src []byte, options []string) ([]byte, error) {
	// FIXME: Make sure we can connect to the pandoc server
	api := new(pandoc.API)
	if api.Settings == nil {
		api.Settings = new(pandoc.Settings)
	}
	if verbose {
		api.Verbose = true
	}
	version, err := api.Version()
	if err != nil {
		return nil, fmt.Errorf("pandoc-server unreachable")
	}
	// NOTE: I need to convert the pandoc command line paremeters into a Pandoc Server all format
	api.PandocOptions(options)

	// NOTE: if we convert the JSON to YAML then I can generator
	// YAML front matter and turn that into a Markdown source doc.
	// I can then pipe it to Pandoc's standard input and skip the
	// whole temp file on disk ready with --metadata-file.
	//
	// The "template" option could still be handled via a tmp file
	// or I could just rely on the options passed to Pandoc to find it.

	// Is `src` JSON or YAML, sniff for "{" prefix.
	m := map[string]interface{}{}
	if bytes.HasPrefix(bytes.TrimSpace(src), []byte(`{`)) {
		// Convert JSON to YAML
		if err := json.Unmarshal(src, &m); err != nil {
			return nil, fmt.Errorf("failed to covert JSON to YAML, %s", err)
		}
	} else {
		src, err = yaml.Marshal(&m)
		if err != nil {
			return nil, fmt.Errorf("failed to generate YAML for front matter, %s", err)
		}
	}
	if len(m) > 0 {
		api.Settings.Metadata = fmt.Sprintf("%s", src)
	}
	src, err = json.Marshal(api.Settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall settings, %s", err)
	}
	if api.Verbose {
		fmt.Fprintf(os.Stderr, "Contacting pandoc server %s\n", version)
	}
	src, err = api.Convert(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("pandoc conversion failed, %s", err)
	}
	return src, nil
}
