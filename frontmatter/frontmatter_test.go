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
	"strings"
	"testing"
)

func TestFrontmatter(t *testing.T) {
	// t.Fatal("not implemented")
	expected_object := map[string]interface{}{
		"title":   "A title of a post",
		"byline":  "R. S. Doiel",
		"pubDate": "2022-10-31T00:00:00Z",
	}
	txt := []byte(`---
title: A title of a post
byline: R. S. Doiel
pubDate: 2022-10-31
---

A title of a post
=================

By R. S. Doiel, 2022-10-30

This is a blog post body.

Bye.
`)

	// Check if ReadAll(buf) works
	buf := bytes.NewBuffer(txt)
	src, err := ReadAll(buf)
	if err != nil {
		t.Errorf("ReadAll error %s", err)
		t.FailNow()
	}
	obj := map[string]interface{}{}
	if err := json.Unmarshal(src, &obj); err != nil {
		t.Errorf("Returnd src is not a JSON object, %s", err)
		t.FailNow()
	}
	for k, v := range expected_object {
		expected := fmt.Sprintf("%s", v)
		got, ok := obj[k]
		if !ok {
			t.Errorf("expected attribute %q, not found", k)
		} else if strings.Compare(expected, got.(string)) != 0 {
			t.Errorf("expected: %s got: %s", expected, got)
			t.FailNow()
		}
	}
	// Now see if TrimFrontmatter works.
	expected_doc := []byte(`
A title of a post
=================

By R. S. Doiel, 2022-10-30

This is a blog post body.

Bye.
`)
	buf = bytes.NewBuffer(txt)
	src, err = TrimFrontmatter(buf)
	if err != nil {
		t.Errorf("TrimFrontmatter(buf) error %s", err)
		t.FailNow()
	}
	if bytes.Compare(expected_doc, src) != 0 {
		t.Errorf("expected doc:\n%s\nGot doc:\n%s\n", expected_doc, src)
		t.FailNow()
	}
}
