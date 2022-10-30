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
}
