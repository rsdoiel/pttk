// pttk.go is a package (with sub-packages) for managing static web content,
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
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
