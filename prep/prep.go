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
	"os/exec"
	"strings"

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
//	out, err := pdtmpl.ReadAll(os.Stdin, "page.tmpl", opt)
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
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, err
	}

	// NOTE: if we convert the JSON to YAML then I can generator
	// YAML front matter and turn that into a Markdown source doc.
	// I can then pipe it to Pandoc's standard input and skip the
	// whole temp file on disk ready with --metadata-file.
	//
	// The "template" option could still be handled via a tmp file
	// or I could just rely on the options passed to Pandoc to find it.

	// Is `src` JSON or YAML, sniff for "{" prefix.
	if bytes.HasPrefix(bytes.TrimSpace(src), []byte(`{`)) {
		// Convert JSON to YAML
		m := map[string]interface{}{}
		if err := json.Unmarshal(src, &m); err != nil {
			return nil, fmt.Errorf("failed to covert JSON to YAML, %s", err)
		}
		src, err = yaml.Marshal(&m)
		if err != nil {
			return nil, fmt.Errorf("failed to generate YAML for front matter, %s", err)
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "%s %s\n", pandoc, strings.Join(options, " "))
	}
	cmd := exec.Command(pandoc, options...)
	if len(bytes.TrimSpace(src)) > 0 {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, fmt.Errorf("could not setup pipe for standard input, %s", err)
		}
		go func() {
			defer stdin.Close()
			fmt.Fprintf(stdin, "---\n%s\n---\n\n", src)
		}()
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	errMsg, _ := io.ReadAll(stderr)
	src, _ = io.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		if len(errMsg) > 0 {
			return nil, fmt.Errorf("%s, %s\n", errMsg, err)
		}
		return nil, err
	}
	return src, nil
}
