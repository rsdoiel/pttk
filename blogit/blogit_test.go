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
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func TestPrivateFuncs(t *testing.T) {
	rsdoiel := "R. S. Doiel"
	obj := map[string]string{}
	obj["name"] = rsdoiel
	objs := append([]interface{}{}, obj)
	creators := unpackCreators(objs)
	if len(creators) != 1 {
		t.Errorf("expected one creator got %d", len(creators))
	} else {
		if fmt.Sprintf("%T", creators[0].Name) != "string" {
			t.Errorf("Expected crestors[0].Name to ba a string, got %T", creators[0].Name)
		}
		if strings.Compare(creators[0].Name, rsdoiel) != 0 {
			t.Errorf("expected creators[0].Name = %q, got %q", rsdoiel, creators[0].Name)
		}
	}
	ymd, err := calcYMD("2003-01-02")
	if err != nil {
		t.Errorf("Unexpected err %s, calcYMD()", err)
		t.FailNow()
	}
	if len(ymd) != 3 {
		t.Errorf("expected len(ymd) = 3, got %+v", ymd)
	} else {
		for i, val := range []string{"2003", "01", "02"} {
			if ymd[i] != val {
				t.Errorf("expected %q, got %q", val, ymd[i])
			}
		}
		expectedS := path.Join("blog", "2003", "01", "02")
		resultS, err := calcPath("blog", ymd)
		if err != nil {
			t.Errorf("Unexpected err %s, calcPath()", err)
		}

		if expectedS != resultS {
			t.Errorf("expected %q, got %q", expectedS, resultS)
		}
	}
}

func TestExportedFuncs(t *testing.T) {
	var (
		pName      string
		prefix     string
		blogPrefix string
		blogJSON   string
		src        []byte
		blogMeta   *BlogMeta
	)
	prefix = "test"
	blogPrefix = path.Join(prefix, "blog")
	blogJSON = path.Join(blogPrefix, "blog.json")

	// Start with an empty blog ...
	os.RemoveAll(path.Join("test", "blog"))
	os.MkdirAll("test", 0777)
	blogMeta = new(BlogMeta)

	// Generate and write test data for BlogIt()
	for i := 1; i <= 10; i = i + 1 {
		src = []byte(fmt.Sprintf(`{
	"title": "Hello No. %d",
	"subtitle": "This is the %d(th) test blog post",
	"date": "2021-05-%02d",
	"keywords": [ "test" ],
	"creators": [ "R. S. Doiel" ],
	"byline": "By R. S. Doiel"
}


# Hello World!

Test Blog post.
`, i, i, i))
		pName = path.Join(prefix, fmt.Sprintf("hello_%d.md", i))
		if err := os.WriteFile(pName, src, 0666); err != nil {
			t.Errorf("Can't created %q, %s", pName, err)
			t.FailNow()
		}

		dateString := fmt.Sprintf("2021-05-%02d", i)
		if err := blogMeta.BlogIt(blogPrefix, pName, dateString); err != nil {
			t.Errorf("BlogIt(%q, %q, %q) failed, %s", blogPrefix, pName, dateString, err)
			t.FailNow()
		}
		if err := blogMeta.Save(blogJSON); err != nil {
			t.Errorf("Failed to write %q, %s", blogJSON, err)
			t.FailNow()
		}
	}
}

func TestRefreshFromPath(t *testing.T) {
	meta := new(BlogMeta)
	year := "2021"
	prefix := "test"
	blogPrefix := path.Join(prefix, "blog")
	blogJSON := path.Join(blogPrefix, "blog.json")
	os.Remove(blogJSON)
	if err := meta.RefreshFromPath(blogPrefix, year); err != nil {
		t.Errorf("expected nil, got %s", err)
		t.FailNow()
	}
	meta.Save(blogJSON)
}
