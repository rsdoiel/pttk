// ws.go is a sub-package pdtk. A packages for managing static content
// blogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package ws

const (
	helpText = `

USAGE

    {app_name} {verb} [HTDOC_PATH] [OPTIONS]

SYNOPSIS

{app_name} {verb} provides a simple static web server for
testing the content you're rendering with Pandoc (or
other static site generator).

EXAMPLE

In the example the htdoc directory is called "myblog"
and you can view the result at http://localhost:8000.

  {app_name} {verb} $HOME/Sites/myblog

`
)
