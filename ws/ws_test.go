// ws.go is a sub-package pdtk. A packages for managing static content
// blogs and documentation via Pandoc.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// Copyright (c) 2022, R. S. Doiel
//
package ws

import (
	"testing"
)

func TestIsDotPath(t *testing.T) {
	boolTests := map[string]bool{
		"":                        false,
		".":                       false,
		"./":                      false,
		"..":                      false,
		"recent/articles":         false,
		"./something else":        false,
		"./something/.git/config": true,
		"../../../":               false,
		".git":                    true,
		".ssh":                    true,
		"../../reoirwepoiewr/../poierwer/../.git/ewrpoiewrrwe/../../": false,
		"../../reoirwepoiewr/../poierwer/../.git/ewrpoiewrrwe/..":     true,
	}

	for p, expected := range boolTests {
		r := IsDotPath(p)
		if r != expected {
			t.Errorf("expected %t, got %t for %s", expected, r, p)
		}
	}
}
