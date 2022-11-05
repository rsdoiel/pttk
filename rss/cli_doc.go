// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package rss

const (
	helpText = `
NAME

   {app_name} {verb} PATH_TO_SITE

SYNOPSIS

The {verb} renders an RSS file based on the content found in the
directory tree provided. If it encounters a "blog.json" file then
it'll use that file to generate feed content for that directory
and it's content otherwise it'll generate a feed backed on Markdown
front matter encountered in Markdown documents with corresponding
html file.

{app_name} {verb} walks the file system to generate a RSS2 file. It assumes 
that the directory for HTDOCS is is the base directory containing 
subdirectories in the form of /YYYY/MM/DD/ARTICLE_HTML where 
YYYY/MM/DD (Year, Month, Day) corresponds to the publication date 
of ARTICLE_HTML.

If our htdocs folder is our document root and out blog is
htdocs/myblog.

    {app_name} {verb} -channel-title="This Great Beyond" \
        -channel-description="Blog to save the world" \
        -channel-link="http://blog.example.org" \
        htdocs htdocs/rss.xml

This would build an RSS 2 file in htdocs/rss.xml from the
articles in htdocs/myblog/YYYY/MM/DD.

DESCRIPTION

EXAMPLE

ALSO SEE


`
)
