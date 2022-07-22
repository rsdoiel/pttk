// prep.go is a sub-package pdtk. A packages for managing static content
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package prep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
//```shell
//    // Options passed to Pandoc
//    opt := []string{}
//    out, err := pdtmpl.ReadAll(os.Stdin, "page.tmpl", opt)
//    if err != nil {
//       // ... handle error
//    }
//    fmt.Fprintf(os.Stdout, "%s\n", out)
//```
//
func ReadAll(r io.Reader, options []string) ([]byte, error) {
	// Read the JSON input
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Apply(src, options)
}

// ReadFile reads a JSON or YAML document from a file then uses Apply
// and options returning a slice of bytes and error value.
func ReadFile(name string, options []string) ([]byte, error) {
	// Read the JSON or YAML file
	src, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return Apply(src, options)
}

// ApplyIO reads in JSON from an io.Reader, applies the template
// and parameters via Format() writing the result to the io.Writer.
// returns an error value.
//
//```
//  // Options passed to Pandoc
//  opt := []string{}
//  err := pdtmpl.ApplyIO(os.Stdin, os.Stdout, "example.tmpl", opt)
//  if err != nil {
//     // ... handle error
//  }
//```
//
func ApplyIO(r io.Reader, w io.Writer, options []string) error {
	src, err := ReadAll(r, options)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", src)
	return err
}

// Apply takes a byte array (like you could read from os.Stdin
// containing JSON or YAML. It creates a temp file and passes that to
// Pandoc via `--metadata-file` option along with any additional
// pandoc options provided. Pandoc then renders the output either
// using the template name (if non-empty string) and the
// additional options passed to Pandoc.
//
//```
//  src, err := ioutil.ReadFile("example.json")
//  if err != nil {
//     // ... handle error
//  }
//  // Options passed to Pandoc
//  opt := []string{"--template", "example.tmpl"}
//  src, err := pdtmpl.Apply(src, opt)
//  if err != nil {
//     // ... handle error
//  }
//  fmt.Printf("%s\n", src)
//```
//
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
	errMsg, _ := ioutil.ReadAll(stderr)
	src, _ = ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		if len(errMsg) > 0 {
			return nil, fmt.Errorf("%s, %s\n", errMsg, err)
		}
		return nil, err
	}
	return src, nil
}
