// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
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
