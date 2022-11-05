// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package pandoc

import (
	"bytes"
	"os"
	"testing"
)

func TestPandoc(t *testing.T) {
	mdText = `---
title: "Hello World!"
author: "jane.doe@example.org (Jane Doe)"
byline: "Jane Doe"
pubDate: 2022-11-05
---

Hello World!
============

By Jane Doe, 2022-11-05

I exist in the virtual space and with that I say "Hello World!".

`

	// FIXME: make sure we have pandoc-server running

	setupSource := []byte(`{
  "port": ":3030",
  "host": "localhost",
  "settings": {
	"from": "markdown",
	"to": "html5",
	"standalone": false
  }
}`)
	if _, err := os.Stat("testout"); os.IsNotExist(err) {
		os.MkdirAll("testout", 0775)
	}
	if err := os.WriteFile("testout/pandoc-setup.json", setupSource, 0664); err != nil {
		t.Error(err)
		t.FailNow()
	}
	api, err := Load("testout/pandoc-setup.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// Use verbose logging for tests.
	api.Verbose = true
	src, err := api.Convert(bytes.NewReader(mdText))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(src) == 0 {
		t.Errorf("Expected content returned from cfg.Convert(), got none")
		t.FailNow()
	}
}
