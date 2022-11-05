// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package frontmatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	LF      = byte(10)
	sMarker = "---\n"
	bMarker = []byte(sMarker)
)

// extractFrontmatter
func extractFrontmatter(src []byte) ([]byte, error) {
	out := bytes.Buffer{}
	if bytes.HasPrefix(src, bMarker) {
		src = bytes.TrimPrefix(src, bMarker)
		buf := bytes.NewBuffer(src)
		inFrontmatter := true
		line, err := buf.ReadString(LF)
		for (err != io.EOF) && inFrontmatter {
			if strings.Compare(line, "---\n") == 0 {
				inFrontmatter = false
			}
			if inFrontmatter {
				if _, err := out.Write([]byte(line)); err != nil {
					return nil, fmt.Errorf("save buffer failed, %s\n", err)
				}
			}
			line, err = buf.ReadString(LF)
		}
	}
	return out.Bytes(), nil
}

// hasFrontmatter returns true if frontmatter if found at
// start of buffer, false otherwise.
func hasFrontmatter(src []byte) bool {
	// Make a copy of buffer contents, so we can see if
	// we have frontmatter or not.
	return bytes.HasPrefix(src, bMarker)
}

// TrimFrontmatter returns the contents of an io buffer removing
// any FrontMatter found.
func TrimFrontmatter(buf io.Reader) ([]byte, error) {
	src, err := io.ReadAll(buf)
	// Check to see if we have frontmatter
	if err == nil && hasFrontmatter(src) {
		in := bytes.NewBuffer(src)
		// Read rest of frontmatter until closing
		inFrontmatter := true
		// Skip our open frontmatter marker
		line, err := in.ReadString(LF)
		// Find the closing frontmatter marker
		for (err != io.EOF) && inFrontmatter {
			line, err = in.ReadString(LF)
			if err == nil && strings.Compare(line, sMarker) == 0 {
				inFrontmatter = false
			}
		}
		// Now copy the rest of the document to out buffer
		out := bytes.Buffer{}
		for err != io.EOF {
			line, err = in.ReadString(LF)
			if err == nil {
				out.Write([]byte(line))
			}
		}
		return out.Bytes(), nil
	}
	// No frontmatter found, return buf as slice of bytes
	return src, err
}

// ReadFile reads a file an extracts front matter converting it to an
// JSON document.
func ReadFile(fName string) ([]byte, error) {
	txt, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	fmText, err := extractFrontmatter(txt)
	if err != nil {
		return nil, err
	}
	// Convert fmText YAML to JSON
	if len(fmText) > 0 {
		obj := map[string]*interface{}{}
		err := yaml.Unmarshal(fmText, &obj)
		if err != nil {
			return nil, err
		}
		return json.MarshalIndent(obj, "", "    ")
	}
	return nil, nil
}

// ReadAll reads a io buffer and extracts the frontmatter from it.
func ReadAll(buf io.Reader) ([]byte, error) {
	txt, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	fmText, err := extractFrontmatter(txt)
	if err != nil {
		return nil, err
	}
	// Convert fmText YAML to JSON
	if len(fmText) > 0 {
		obj := map[string]*interface{}{}
		err := yaml.Unmarshal(fmText, &obj)
		if err != nil {
			return nil, err
		}
		return json.MarshalIndent(obj, "", "    ")
	}
	return nil, nil
}
