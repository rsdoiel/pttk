// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package rss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	// My packages
	"github.com/rsdoiel/fountain"
	"github.com/rsdoiel/pttk/blogit"
	"github.com/rsdoiel/pttk/frontmatter"
	// 3rd Part support (e.g. YAML)
)

const (
	// DateExp is the default format used by mkpage utilities for date exp
	DateExp = `[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]`
	// BylineExp is the default format used by mkpage utilities
	BylineExp = (`^[B|b]y\s+(\w|\s|.)+` + DateExp + "$")
	// TitleExp is the default format used by mkpage utilities
	TitleExp = `^#\s+(\w|\s|.)+$`
)

var (
	// Config holds a global config.
	// Uses the same structure as Front Matter in that it is
	// the result of parsing TOML, YAML or JSON into a
	// map[string]interface{} tree
	Config map[string]interface{}
)

// normalizeEOL takes a []byte and normalizes the end of line
// to a `\n' from `\r\n` and `\r`
func normalizeEOL(input []byte) []byte {
	if bytes.Contains(input, []byte("\r\n")) {
		input = bytes.Replace(input, []byte("\r\n"), []byte("\n"), -1)
	}
	return input
}

// ProcessorConfig takes front matter and returns
// a map[string]interface{} containing configuration
func ProcessorConfig(fmSrc []byte) (map[string]interface{}, error) {
	//FIXME: Need to merge with .Config and return the merged result.
	m := map[string]interface{}{}
	// Do nothing is we have zero front matter to process.
	if len(fmSrc) == 0 {
		return m, nil
	}
	src, err := frontmatter.ReadAll(bytes.NewBuffer(fmSrc))
	if err != nil {
		return nil, err
	}
	// JSON Front Matter
	if err := json.Unmarshal(src, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConfigFountain sets the fountain defaults then applies
// the map[string]interface{} overwriting the defaults
// returns error necessary.
func ConfigFountain(config map[string]interface{}) error {
	if thing, ok := config["fountain"]; ok == true {
		cfg := thing.(map[string]interface{})
		for k, v := range cfg {
			switch v.(type) {
			case bool:
				onoff := v.(bool)
				switch k {
				case "AsHTMLPage":
					fountain.AsHTMLPage = onoff
				case "InlineCSS":
					fountain.InlineCSS = onoff
				case "LinkCSS":
					fountain.LinkCSS = onoff
				}
			case string:
				if k == "IncludeCSS" {
					fountain.CSS = v.(string)
				}
			default:
				return fmt.Errorf("Unknown fountain option %q", k)
			}
		}
	}
	return nil
}

// fountainProcessor wraps fountain.Run() splitting off the front
// matter if present.
func fountainProcessor(input []byte) ([]byte, error) {
	var err error

	fmSrc, err := frontmatter.ReadAll(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	config, err := ProcessorConfig(fmSrc)
	if err != nil {
		return nil, err
	}
	if err := ConfigFountain(config); err != nil {
		return nil, err
	}
	fountainSrc, err := frontmatter.TrimFrontmatter(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	src, err := fountain.Run(fountainSrc)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// RelativeDocPath calculate the relative path from source to target based on
// implied common base.
//
// Example:
//
//	docPath := "docs/chapter-01/lesson-02.html"
//	cssPath := "css/site.css"
//	fmt.Printf("<link href=%q>\n", MakeRelativePath(docPath, cssPath))
//
// Output:
//
//	<link href="../../css/site.css">
func RelativeDocPath(source, target string) string {
	var result []string

	sep := string(os.PathSeparator)
	dname, _ := path.Split(source)
	for i := 0; i < strings.Count(dname, sep); i++ {
		result = append(result, "..")
	}
	result = append(result, target)
	p := strings.Join(result, sep)
	if strings.HasSuffix(p, "/.") {
		return strings.TrimSuffix(p, ".")
	}
	return p
}

// NormalizeDate takes a MySQL like date string and returns a time.Time or error
func NormalizeDate(s string) (time.Time, error) {
	switch len(s) {
	case len(`2006-01-02 15:04:05 -0700`):
		dt, err := time.Parse(`2006-01-02 15:04:05 -0700`, s)
		return dt, err
	case len(`2006-01-02 15:04:05`):
		dt, err := time.Parse(`2006-01-02 15:04:05`, s)
		return dt, err
	case len(`2006-01-02`):
		dt, err := time.Parse(`2006-01-02`, s)
		return dt, err
	default:
		return time.Time{}, fmt.Errorf("Can't format %s, expected format like 2006-01-02 15:04:05 -0700", s)
	}
}

// Walk takes a start path and walks the file system to process Markdown files f or useful elements.
func Walk(startPath string, filterFn func(p string, info os.FileInfo) bool, outputFn func(s string, info os.FileInfo) error) error {
	err := filepath.Walk(startPath, func(p string, info os.FileInfo, err error) error {
		// Are we interested in this path?
		if filterFn(p, info) == true {
			// Yes, so send to output function.
			if err := outputFn(p, info); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Grep looks for the first line matching the expression
// in src.
func Grep(exp string, src string) string {
	re, err := regexp.Compile(exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q is not a valid, %s\n", exp, err)
		return ""
	}
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		s := re.FindString(line)
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

// Generate a Feed from walking the blogit.BlogMeta structure
func BlogMetaToRSS(blog *blogit.BlogMeta, feed *RSS2) error {
	if len(blog.Name) > 0 {
		feed.Title = blog.Name
	}
	if len(blog.BaseURL) > 0 {
		feed.Link = blog.BaseURL
	}
	if len(blog.Quip) > 0 {
		feed.Description = "> " + blog.Quip + "\n\n"
	}
	if len(blog.Description) > 0 {
		feed.Description += blog.Description
	}
	if len(blog.Updated) > 0 {
		dt, err := time.Parse("2006-01-02", blog.Updated)
		if err != nil {
			return err
		}
		feed.PubDate = dt.Format(time.RFC1123Z)
		feed.LastBuildDate = dt.Format(time.RFC1123Z)
	}
	if len(blog.Language) > 0 {
		feed.Language = blog.Language
	}
	if len(blog.Copyright) > 0 {
		feed.Copyright = blog.Copyright
	}
	//FIXME: Need to iterate over years, months, days and build our
	// blog items.
	for _, years := range blog.Years {
		yr := years.Year
		for _, months := range years.Months {
			mn := months.Month
			for _, days := range months.Days {
				dy := days.Day
				for _, post := range days.Posts {
					pubDate, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", yr, mn, dy))
					if err != nil {
						return err
					}
					// NOTE: We only want to process Markdown documents.
					// We look for Markdown related file extensions.
					includeDescription := false
					linkPath := post.Document
					for _, ext := range []string{".md", ".markdown", ".txt", ".asciidoc"} {
						if strings.HasSuffix(post.Document, ext) {
							includeDescription = true
							linkPath = strings.TrimSuffix(post.Document, ext) + ".html"
						}
					}
					item := new(Item)
					item.Title = post.Title
					if strings.Contains(blog.BaseURL, "://") {
						item.Link = strings.Join([]string{blog.BaseURL, linkPath}, "/")
					} else {
						item.Link = strings.TrimSuffix(feed.Link, "/") + "/" + strings.TrimPrefix(linkPath, "/")
					}
					if strings.Contains(item.Link, "://") {
						item.GUID = item.Link
					} else {
						item.GUID = strings.TrimSuffix(feed.Link, "/") + "/" + strings.TrimPrefix(item.Link, "/")
					}
					item.PubDate = pubDate.Format(time.RFC1123Z)
					if len(post.Description) == 0 && len(post.Document) > 0 {
						// Read the article, extract a description
						buf, err := os.ReadFile(post.Document)
						if err != nil {
							return err
						}
						fMatter := map[string]interface{}{}
						fSrc, err := frontmatter.ReadAll(bytes.NewBuffer(buf))
						if err != nil {
							return err
						}
						tSrc, err := frontmatter.TrimFrontmatter(bytes.NewBuffer(buf))
						if err != nil {
							return err
						}
						if len(fSrc) > 0 {
							if err := json.Unmarshal(fSrc, &fMatter); err != nil {
								fMatter = map[string]interface{}{}
							}
						}
						if val, ok := fMatter["description"]; ok {
							post.Description = val.(string)
						} else if includeDescription {
							post.Description = OpeningParagraphs(fmt.Sprintf("%s", tSrc), 5, "\n\n")
							if len(post.Description) < len(tSrc) {
								post.Description += " ..."
							}
						}
					}
					if len(post.Abstract) > 0 {
						item.Description = post.Abstract
					}
					if len(post.Description) > 0 {
						item.Description = post.Description
					}
					if item.Title != "" || item.Description != "" {
						feed.ItemList = append(feed.ItemList, *item)
					}
				}
			}
		}
	}
	return nil
}

// Generate a Feed by walking the file system.
func WalkRSS(feed *RSS2, htdocs string, baseURL string, excludeList string, titleExp string, bylineExp string, dateExp string) error {
	validBlogPath := regexp.MustCompile("/[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]/")
	err := Walk(htdocs, func(p string, info os.FileInfo) bool {
		fname := path.Base(p)
		if validBlogPath.MatchString(p) == true &&
			strings.HasSuffix(fname, ".md") == true {
			// NOTE: We have a possible published markdown article.
			// Make sure we have a HTML version before adding it
			// to the feed.
			if _, err := os.Stat(path.Join(p, path.Base(fname)+".html")); os.IsNotExist(err) {
				return false
			}
			return true
		}
		return false
	}, func(p string, info os.FileInfo) error {
		// Read the article
		buf, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		fMatter := map[string]interface{}{}
		fSrc, err := frontmatter.ReadAll(bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
		tSrc, err := frontmatter.TrimFrontmatter(bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
		if len(fSrc) > 0 {
			if err := json.Unmarshal(fSrc, &fMatter); err != nil {
				fMatter = map[string]interface{}{}
			}
		}

		// Calc URL path
		pname := strings.TrimPrefix(p, htdocs)
		if strings.HasPrefix(pname, "/") {
			pname = strings.TrimPrefix(pname, "/")
		}
		dname := path.Dir(pname)
		bname := strings.TrimSuffix(path.Base(pname), ".md") + ".html"
		articleURL := fmt.Sprintf("%s/%s", baseURL, path.Join(dname, bname))
		u, err := url.Parse(articleURL)
		if err != nil {
			return err
		}
		// Collect metadata
		//NOTE: Use front matter if available otherwise
		var (
			title, byline, author, description, pubDate string
		)
		src := fmt.Sprintf("%s", buf)
		if val, ok := fMatter["title"]; ok {
			title = val.(string)
		} else {
			title = strings.TrimPrefix(Grep(titleExp, src), "# ")
		}
		if val, ok := fMatter["byline"]; ok {
			byline = val.(string)
		} else {
			byline = Grep(bylineExp, src)
		}
		if val, ok := fMatter["pubDate"]; ok {
			switch val.(type) {
			case string:
				pubDate = val.(string)
			case time.Time:
				dt := val.(time.Time)
				pubDate = dt.Format(blogit.DateFmt)
			}
		} else {
			pubDate = Grep(dateExp, byline)
		}
		if val, ok := fMatter["description"]; ok {
			description = val.(string)
		} else {
			description = OpeningParagraphs(fmt.Sprintf("%s", tSrc), 5, "\n\n")
			if len(description) < len(tSrc) {
				description += " ..."
			}
		}
		if val, ok := fMatter["creator"]; ok {
			author = val.(string)
		} else if val, ok = fMatter["author"]; ok {
			author = val.(string)
		} else {
			author = byline
			if len(byline) > 2 {
				author = strings.TrimSpace(strings.TrimSuffix(byline[2:], pubDate))
			}
		}
		// Reformat pubDate to conform to RSS2 date formats
		var (
			dt time.Time
		)
		if pubDate == "" {
			dt = time.Now()
		} else {
			dt, err = time.Parse(`2006-01-02`, pubDate)
			if err != nil {
				return err
			}
		}
		pubDate = dt.Format(time.RFC1123Z)
		item := new(Item)
		if strings.Contains(articleURL, "://") {
			item.GUID = articleURL
		} else {
			item.GUID = strings.TrimSuffix(feed.Link, "/") + "/" + strings.TrimPrefix(articleURL, "/")
		}
		item.Title = title
		item.Author = author
		item.PubDate = pubDate
		item.Link = u.String()
		item.Description = description
		feed.ItemList = append(feed.ItemList, *item)
		return nil
	})
	return err
}

var (
	// Set a default -f (from) value used by Pandoc
	PandocFrom string
	// Set a default -t (to) value used by Pandoc
	PandocTo string
)

// MetadataBlock holds the Pandoc style Metadata block delimited
// by start '%' at the being of the line in the start of a text file.
type MetadataBlock struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Date    string   `json:"date"`
}

func (block *MetadataBlock) String() string {
	return fmt.Sprintf("%% %s\n%% %s\n%% %s", block.Title, strings.Join(block.Authors, "; "), block.Date)
}

func (block *MetadataBlock) Unmarshal(src []byte) error {
	lines := bytes.Split(src, []byte("\n"))
	fieldCnt := 0
	key := ""
	block.Title = ""
	block.Authors = []string{}
	block.Date = ""
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if bytes.HasPrefix(line, []byte("% ")) {
			fieldCnt += 1
			switch fieldCnt {
			case 1:
				key = "title"
			case 2:
				key = "authors"
			case 3:
				key = "date"
			default:
				key = ""
			}
			line = bytes.TrimPrefix(line, []byte("% "))
		}
		if (len(key) > 0) && (fieldCnt <= 3) {
			switch key {
			case "title":
				if len(block.Title) > 0 {
					block.Title = fmt.Sprintf("%s\n%s", block.Title, bytes.TrimSpace(line))
				} else {
					block.Title = fmt.Sprintf("%s", bytes.TrimSpace(line))
				}
			case "authors":
				if bytes.Contains(line, []byte(";")) {
					parts := bytes.Split(line, []byte(";"))
					for _, part := range parts {
						block.Authors = append(block.Authors, fmt.Sprintf("%s", bytes.TrimSpace(part)))
					}
				} else {
					block.Authors = append(block.Authors, fmt.Sprintf("%s", bytes.TrimSpace(line)))
				}
			case "date":
				block.Date = fmt.Sprintf("%s", bytes.TrimSpace(line))
				key = ""
			}
		}
	}
	if fieldCnt != 3 {
		return fmt.Errorf("Missing or ill formed metablock, expecting title, author(s), date")
	}
	return nil
}

func (block *MetadataBlock) Marshal() ([]byte, error) {
	return json.Marshal(block)
}

func scanArgs(s string) (string, []string) {
	var (
		tok       string
		generator string
		params    []string
		i, j      int
	)
	for i = 0; i < len(s) && (tok != " "); i++ {
		tok = string(s[i])
	}
	generator = strings.TrimSpace(string(s[0:i]))
	params = []string{}
	j = len(generator) + 1
	for ; i < len(s); i++ {
		tok = string(s[i])
		switch tok {
		case "'":
			for ; i < len(s) && tok != "'"; i++ {
				// advance to next single quote.
				tok = string(s[i])
				if tok == "\\" {
					i += 1
					tok = string(s[i])
				}
			}
		case `"`:
			for ; i < len(s) && tok != `"`; i++ {
				// advance to next double quote.
				tok = string(s[i])
				if tok == "\\" {
					i += 1
					tok = string(s[i])
				}
			}
		case " ":
			params = append(params, strings.TrimSpace(string(s[j:i])))
			j = i
		}
	}
	if j < i {
		params = append(params, strings.TrimSpace(string(s[j:i])))
	}
	return generator, params
}

// JSONGenerator accepts  command line string and executes it.
// It take command's output, validates that it is JSON and returns it.
func JSONGenerator(cmdExpr string, obj interface{}) error {
	var (
		out, eOut bytes.Buffer
		generator string
		params    []string
		err       error
	)
	//NOTE: We use the scanner because we want to treat quote strings
	// as one parameter.
	generator, params = scanArgs(cmdExpr)
	cmd := exec.Command(generator, params...)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", cmdExpr, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", cmdExpr, err)
		}
		return err
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", cmdExpr, eOut.String())
		return err
	}
	src := out.Bytes()
	//NOTE: Validate our JSON by trying to unmarshaling it
	err = json.Unmarshal(src, &obj)
	if err != nil {
		err = fmt.Errorf("Invalid JSON from %q exit error, %s", cmdExpr, err)
	}
	return err
}

// OpeningParagraphs scans a Markdown file and attempts to copy the
// first `cnt` paragraphs ignoring GitHub image inserts. Paragraphs are
// scanned text followed by a delimiting character, typically that's a two
// new line sequence.
//
// Example:
//
//	src, _ := os.ReadFile("post.md")
//	opening := mkpage.OpenParagraphs(fmt.Sprintf("%s", src, "\n\n"), 2)
func OpeningParagraphs(src string, cnt int, para string) string {
	blocks := strings.Split(strings.ReplaceAll(src, "\r", ""), para)
	n := 0
	txt := []string{}
	// Find blocks
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		//NOTE: It would be good to handle more than Markdown
		// But this is a simple test for now.
		skipBlock := false
		// Normalize text so to make it easier to detect context
		tblock := strings.TrimSpace(block)
		// Skip Titles and badges
		if strings.HasSuffix(tblock, "====") ||
			strings.HasSuffix(tblock, "----") ||
			strings.HasPrefix(tblock, "#") ||
			strings.HasPrefix(tblock, "[!") {
			skipBlock = true
		}
		if byline := Grep(bylineExp, tblock); byline != "" {
			skipBlock = true
		}
		if !skipBlock {
			n += 1
			txt = append(txt, block)
		}
		if n >= cnt {
			break
		}
	}
	return strings.Join(txt, para)
}
