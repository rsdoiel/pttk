// gs.go is a sub-package pdtk. This package for testing a static content
// Gopher site.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package gs

const (
	helpText = `% {app_name}-{verb}(1) {app_name}-{verb} user manual
% R. S. Doiel
% September 23, 2022

# NAME

{app_name} {verb}

# SYNOPSIS

{app_name} {verb} [HTDOC_PATH] [OPTIONS]

# DESCRIPTION

{app_name} {verb} provides a simple static gopher server for
testing the content you're Gopher content.

# EXAMPLE

In the example the htdoc directory is called "myblog"
and you can view the result at gopher://localhost:7000.

  {app_name} {verb} $HOME/Sites/myblog

`
)
