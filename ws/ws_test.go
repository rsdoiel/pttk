// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
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
