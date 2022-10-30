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
	LF = byte(10)
)

// extractFrontmatter
func extractFrontmatter(src []byte) ([]byte, error) {
	out := bytes.Buffer{}
	if bytes.HasPrefix(src, []byte("---\n")) {
		src = bytes.TrimPrefix(src, []byte("---\n"))
		buf := bytes.NewBuffer(src)
		inFrontmatter := true
		line, err := buf.ReadString(LF)
		for err != io.EOF {
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
