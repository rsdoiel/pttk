//
// rss.go is a library for generating RSS 2 documents and is part of the
// ptdk project.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package rss

import (
	"strings"
	"testing"
)

func TestRSS2(t *testing.T) {
	src := []byte(`<?xml version="1.0" encoding="utf-8" ?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:media="http://search.yahoo.com/mrss/">
<channel>
    <title>Such an odd life.</title>
    <link>https://example.com</link>
    <atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="self" href="https://example.com/feeds/index.xml" type="application/rss+xml"></atom:link>
    <description>This is an example of feeds in a personal website. </description>
	<image>
        <url>https://example.com/assets/cool-logo.svg</url>
        <title>The cool logo of example.com</title>
        <link>https://example.com/assets/cool-logo.svg</link>
    </image>
    <pubDate>Mon, 25 Jul 2022 11:47:00 -0800</pubDate>
    <lastBuildDate>Mon, 25 Jul 2022 12:00:00 -0800</lastBuildDate>
    <language>en</language>
    <copyright>2022</copyright>
    <item>
        <pubDate>Mon, 25 Jul 2022 00:08:03 -0800</pubDate>
        <title>The PDTK Tutorial</title>
        <link>https://example.com/blog/2022/07/25/index.html</link>
        <guid>https://example.com/blog/2022/07/25/index.html?v0.0.1</guid>
        <description>This is a tutorial on the speculative pdtk tool and Go package.</description>
    </item>
    <item>
        <pubDate>Mon, 20 Jul 2022 10:26:26 -0800</pubDate>
        <title>Introducing PDTK, a Pandoc Tool Kit</title>
        <link>https://examples.com/blog/2022/07/20/index.html</link>
        <guid>https://examples.com/blog/2022/07/20/index.html?v0.0.2</guid>
        <description>A blog post on creating the pdtk project and updating the tools I use for my blog.</description>
	</item>
</channel>
</rss>`)

	r, err := Parse(src)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if r == nil {
		t.Errorf("expected a populated RSS2 struct, got nil")
		t.FailNow()
	}
	expectedItems := 2
	gotItems := len(r.Items)
	if expectedItems != gotItems {
		t.Errorf("expected %d, got %d items, %+v", expectedItems, gotItems, r.Items)
	}
	expectedCopyright := "2022"
	gotCopyright := r.Copyright
	if strings.Compare(expectedCopyright, gotCopyright) != 0 {
		t.Errorf("expected %q, got %q", expectedCopyright, gotCopyright)
	}
}
