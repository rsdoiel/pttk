// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package pttk

import (
	"io"

	"github.com/rsdoiel/pttk/prep"
)

//
// NOTE: Below are a mapping of the prep preprocessor functions
// to the pttk package level.
//

func SetVerbose(onoff bool) {
	prep.SetVerbose(onoff)
}

func ReadAll(r io.Reader, options []string) ([]byte, error) {
	return prep.ReadAll(r, options)
}

func ReadFile(name string, options []string) ([]byte, error) {
	return prep.ReadFile(name, options)
}

func ApplyIO(r io.Reader, w io.Writer, options []string) error {
	return prep.ApplyIO(r, w, options)
}

func Apply(src []byte, options []string) ([]byte, error) {
	return prep.Apply(src, options)
}
