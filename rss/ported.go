package rss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	// My packages
	"github.com/rsdoiel/pdtk/blogit"

	// Caltech Library packages
	"github.com/caltechlibrary/rss2"

	// 3rd Part support (e.g. YAML)
	"gopkg.in/yaml.v3"

	// Fountain support for scripts, interviews and narration
	"github.com/rsdoiel/fountain"
)

const (
	// Prefix for explicit string types

	// JSONPrefix designates a string as JSON formatted content
	JSONPrefix = "json:"
	// CommonMarkPrefix designates a string as Common Mark
	// (a rich markdown dialect) content
	CommonMarkPrefix = "commonmark:"
	// MarkdownPrefix designates a string as Markdown (pandoc's dialect)
	// content
	MarkdownPrefix = "markdown:"
	// MarkdownStrict designates a strnig as John Gruber's Markdown content
	MarkdownStrictPrefix = "markdown_strict:"
	// GfmMarkdownPrefix designates a string as GitHub Flavored Markdown
	GfmMarkdownPrefix = "gfm:"
	// MMarkPrefix designates MMark format, for now this will just be passed to pandoc.
	MMarkPrefix = "mmark:"
	// TextPrefix designates a string as text/plain not needed processing
	TextPrefix = "text:"
	// FountainPrefix designates a string as Fountain formatted content
	FountainPrefix = "fountain:"
	// TextilePrefix designates source as Textile for processing by pandoc.
	TextilePrefix = "textile:"
	// ReStructureText designates source as ReStructureText for processing by pandoc
	ReStructureTextPrefix = "rst:"
	// JiraPrefix markup designates source as Jire text for processing by pandoc
	JiraPrefix = "jira:"
	// JSONGeneratorPrefix evaluates the value as a command line that
	// returns JSON.
	JSONGeneratorPrefix = "json-generator:"

	// DateExp is the default format used by mkpage utilities for date exp
	DateExp = `[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]`
	// BylineExp is the default format used by mkpage utilities
	BylineExp = (`^[B|b]y\s+(\w|\s|.)+` + DateExp + "$")
	// TitleExp is the default format used by mkpage utilities
	TitleExp = `^#\s+(\w|\s|.)+$`

	//
	// Supported types for Front Matter
	//

	// FrontMatterIsUnknown means front matter and we can't parse it
	FrontMatterIsUnknown = iota
	// FrontMatterIsJSON means we have detected JSON front matter
	FrontMatterIsJSON
	// FrontMatterIsPandocMetadata means we have detected a Pandoc
	// style metadata block, e.g. opening lines start with
	// '%' attribute name followed by value(s)
	// E.g.
	//      % title
	//      % author(s)
	//      % date
	FrontMatterIsPandocMetadata
	// FrontMatterIsYAML means we have detected a Pandoc YAML
	// front matter block.
	FrontMatterIsYAML
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
	/*
		if bytes.Contains(input, []byte("\r")) {
			input = bytes.Replace(input, []byte("\r"), []byte("\n"), -1)
		}
	*/
	return input
}

// SplitFrontMatter takes a []byte input splits it into front matter type,
// front matter source and Markdown source. If either is missing an
// empty []byte is returned for the missing element.
// NOTE: removed yaml, toml support as of v0.2.4
// NOTE: Added support for Pandoc title blocks v0.2.5
func SplitFrontMatter(input []byte) (int, []byte, []byte) {
	// JSON front matter, most Markdown processors.
	if bytes.HasPrefix(input, []byte("{\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("{\n")), []byte("\n}\n"), 2)
		src := []byte(fmt.Sprintf("{\n%s\n}\n", parts[0]))
		if len(parts) > 1 {
			return FrontMatterIsJSON, src, parts[1]
		}
		return FrontMatterIsJSON, src, []byte("")
	}
	if bytes.HasPrefix(input, []byte("---\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("---\n")), []byte("\n---\n"), 2)
		src := []byte(fmt.Sprintf("---\n%s\n---\n", parts[0]))
		if len(parts) > 1 {
			return FrontMatterIsYAML, src, parts[1]
		}
		return FrontMatterIsYAML, src, []byte("")
	}
	if bytes.HasPrefix(input, []byte("% ")) {
		lines := bytes.Split(input, []byte("\n"))
		i := 0
		fieldCnt := 0
		src := []byte{}
		for ; (i < len(lines)) && (fieldCnt < 3); i++ {
			if bytes.HasPrefix(lines[i], []byte("% ")) {
				fieldCnt += 1
				src = append(append(src, lines[i]...), []byte("\n")...)
			} else if fieldCnt < 3 {
				//NOTE: Dates can only one line, so we stop extra
				// line consumption with authors.
				src = append(append(src, lines[i]...), []byte("\n")...)
			}
		}
		if fieldCnt == 3 {
			return FrontMatterIsPandocMetadata, src, input[len(src):]
		}
	}
	// Handle case of no front matter
	return FrontMatterIsUnknown, []byte(""), input
}

// UnmarshalFrontMatter takes a []byte of front matter source
// and unmarshalls using only JSON frontmatter
// NOTE: removed yaml, toml support as of v0.2.4
// NOTE: Added support for Pandoc title blocks as of v0.2.5
func UnmarshalFrontMatter(configType int, src []byte, obj *map[string]interface{}) error {
	var (
		txt []byte
		err error
	)
	switch configType {
	case FrontMatterIsPandocMetadata:
		block := MetadataBlock{}
		if err = block.Unmarshal(txt); err != nil {
			return err
		}
		if txt, err = block.Marshal(); err != nil {
			return nil
		}
		if err = json.Unmarshal(txt, &obj); err != nil {
			return err
		}
	case FrontMatterIsJSON:
		// Make sure we have valid JSON
		if err = json.Unmarshal(src, &obj); err != nil {
			return err
		}
	case FrontMatterIsYAML:
		if err = yaml.Unmarshal(src, &obj); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported Front matter format")
	}
	return nil
}

// ProcessorConfig takes front matter and returns
// a map[string]interface{} containing configuration
// NOTE: removed yaml, toml support as of v0.2.4
// NOTE: added Pandoc Metadata block as of v0.2.5
func ProcessorConfig(configType int, frontMatterSrc []byte) (map[string]interface{}, error) {
	//FIXME: Need to merge with .Config and return the merged result.
	m := map[string]interface{}{}
	// Do nothing is we have zero front matter to process.
	if len(frontMatterSrc) == 0 {
		return m, nil
	}
	// Convert Front Matter to JSON
	switch configType {
	case FrontMatterIsPandocMetadata:
		block := MetadataBlock{}
		if err := block.Unmarshal(frontMatterSrc); err != nil {
			return nil, err
		}
		m["title"] = block.Title
		m["authors"] = block.Authors
		m["date"] = block.Date
	case FrontMatterIsJSON:
		// JSON Front Matter
		if err := json.Unmarshal(frontMatterSrc, &m); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown supported front matter format")
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

	configType, frontMatterSrc, fountainSrc := SplitFrontMatter(input)
	config, err := ProcessorConfig(configType, frontMatterSrc)
	if err != nil {
		return nil, err
	}
	if err := ConfigFountain(config); err != nil {
		return nil, err
	}
	src, err := fountain.Run(fountainSrc)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// ResolveData takes a data map and reads in the files and URL sources
// as needed turning the data into strings to be applied to the template.
func ResolveData(data map[string]string) (map[string]interface{}, error) {
	var (
		out map[string]interface{}
	)

	isContentType := func(vals []string, target string) bool {
		for _, h := range vals {
			if strings.Contains(h, target) == true {
				return true
			}
		}
		return false
	}

	out = make(map[string]interface{})
	for key, val := range data {
		switch {
		case strings.HasPrefix(val, TextPrefix) == true:
			out[key] = strings.TrimPrefix(val, TextPrefix)
		case strings.HasPrefix(val, MMarkPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MMarkPrefix)), "markdown_mmd", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, CommonMarkPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, CommonMarkPrefix)), "commonmark_x", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, MarkdownPrefix) == true:
			//NOTE: We're using pandoc's flavor Markdown as our processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MarkdownPrefix)), "markdown", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, MarkdownStrictPrefix) == true:
			//NOTE: We're using origanal John Gruber Markdown
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MarkdownStrictPrefix)), "markdown_strict", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, GfmMarkdownPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, GfmMarkdownPrefix)), "gfm", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, JiraPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, JiraPrefix)), "jira", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, TextilePrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, TextilePrefix)), "textile", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, ReStructureTextPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, ReStructureTextPrefix)), "rst", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, FountainPrefix) == true:
			src, err := fountainProcessor([]byte(strings.TrimPrefix(val, FountainPrefix)))
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, JSONPrefix) == true:
			var o interface{}
			err := json.Unmarshal(bytes.TrimPrefix([]byte(val), []byte(JSONPrefix)), &o)
			if err != nil {
				return out, fmt.Errorf("Can't JSON decode (%s) %s, %s", key, val, err)
			}
			out[key] = o
		case strings.HasPrefix(val, JSONGeneratorPrefix) == true:
			//NOTE: JSONGenerator expects a command line that results
			// in JSON written to stdout. It then passes this back to
			// be processed by pandoc in the metadata file.
			var o interface{}
			cmd := strings.TrimPrefix(val, JSONGeneratorPrefix)
			err := JSONGenerator(cmd, &o)
			if err != nil {
				return out, fmt.Errorf("(key: %q) %q failed, %s", key, cmd, err)
			}
			out[key] = o
		case strings.HasPrefix(val, "http://") == true || strings.HasPrefix(val, "https://") == true:
			resp, err := http.Get(val)
			if err != nil {
				return out, fmt.Errorf("Error from (%s) %s, %s", key, val, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				buf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return out, err
				}
				fmType, fmSrc, docSrc := SplitFrontMatter(buf)
				if len(fmSrc) > 0 {
					buf = docSrc
					fmData := map[string]interface{}{}
					if err := UnmarshalFrontMatter(fmType, fmSrc, &fmData); err != nil {
						return out, fmt.Errorf("Can't process front matter (%s), %q, %q", key, val, err)
					}
					// Update, Overwrite `out` with front matter values
					for k, v := range fmData {
						out[k] = v
					}
				}
				if contentTypes, ok := resp.Header["Content-Type"]; ok == true {
					switch {
					case isContentType(contentTypes, "application/json") == true:
						var o interface{}
						err := json.Unmarshal(buf, &o)
						if err != nil {
							return out, fmt.Errorf("Can't JSON decode (%s) %s, %s", key, val, err)
						}
						out[key] = o
					case isContentType(contentTypes, "text/markdown") == true:
						src, err := pandocProcessor(buf, "", "html")
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					case isContentType(contentTypes, "text/commonmark") == true:
						src, err := pandocProcessor(buf, "commonmark_x", "html")
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					case isContentType(contentTypes, "text/mmark") == true:
						src, err := pandocProcessor(buf, "mmark", "html")
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					case isContentType(contentTypes, "text/fountain") == true:
						src, err := fountainProcessor(buf)
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					default:
						out[key] = string(buf)
					}
				} else {
					out[key] = string(buf)
				}
			}
		default:
			ext := path.Ext(val)
			buf, err := ioutil.ReadFile(val)
			if err != nil {
				return out, fmt.Errorf("Can't read (%s) %q, %s", key, val, err)
			}
			//NOTE: We only split front matter for supported markup
			// formats, e.g. MultiMarkdown, CommonMark, Markdown, Textile,
			// ReStructureText, JiraText, Fountain
			if strings.Compare(ext, ".json") != 0 {
				fmType, fmSrc, docSrc := SplitFrontMatter(buf)
				if len(fmSrc) > 0 {
					buf = docSrc
					fmData := map[string]interface{}{}
					if err := UnmarshalFrontMatter(fmType, fmSrc, &fmData); err != nil {
						return out, fmt.Errorf("Can't process front matter (%s), %q, %q", key, val, err)
					}
					// Update, Overwrite `out` with front matter values
					for k, v := range fmData {
						out[k] = v
					}
				}
			}
			switch {
			case strings.Compare(ext, ".fountain") == 0 ||
				strings.Compare(ext, ".spmd") == 0:
				src, err := fountainProcessor(buf)
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".md") == 0:
				src, err := pandocProcessor(buf, "", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".mmd") == 0:
				src, err := pandocProcessor(buf, "markdown_mmd", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".rst") == 0:
				src, err := pandocProcessor(buf, "rst", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".textile") == 0:
				src, err := pandocProcessor(buf, "textile", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".jira") == 0:
				src, err := pandocProcessor(buf, "jira", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".json") == 0:
				var o interface{}
				err := json.Unmarshal(buf, &o)
				if err != nil {
					return out, fmt.Errorf("Can't JSON decode (%s) %s, %s", key, val, err)
				}
				out[key] = o
			default:
				out[key] = string(buf)
			}
		}
	}
	return out, nil
}

//
// RelativeDocPath calculate the relative path from source to target based on
// implied common base.
//
// Example:
//
//     docPath := "docs/chapter-01/lesson-02.html"
//     cssPath := "css/site.css"
//     fmt.Printf("<link href=%q>\n", MakeRelativePath(docPath, cssPath))
//
// Output:
//
//     <link href="../../css/site.css">
//
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
func BlogMetaToRSS(blog *blogit.BlogMeta, feed *rss2.RSS2) error {
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
		feed.PubDate = dt.Format(time.RFC1123)
		feed.LastBuildDate = dt.Format(time.RFC1123)
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
					for _, ext := range []string{".md", ".markdown", ".txt", ".asciidoc"} {
						if strings.HasSuffix(post.Document, ext) {
							includeDescription = true
						}
					}
					item := new(rss2.Item)
					item.Title = post.Title
					item.Link = strings.Join([]string{blog.BaseURL, post.Document}, "/")
					item.GUID = item.Link
					item.PubDate = pubDate.Format(time.RFC1123)
					if len(post.Description) == 0 && len(post.Document) > 0 {
						// Read the article, extract a description
						buf, err := ioutil.ReadFile(post.Document)
						if err != nil {
							return err
						}
						fMatter := map[string]interface{}{}
						fType, fSrc, tSrc := SplitFrontMatter(buf)
						if len(fSrc) > 0 {
							if err := UnmarshalFrontMatter(fType, fSrc, &fMatter); err != nil {
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
func WalkRSS(feed *rss2.RSS2, htdocs string, excludeList string, titleExp string, bylineExp string, dateExp string) error {
	// Required
	channelLink := feed.Link

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
		buf, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		fMatter := map[string]interface{}{}
		fType, fSrc, tSrc := SplitFrontMatter(buf)
		if len(fSrc) > 0 {
			if err := UnmarshalFrontMatter(fType, fSrc, &fMatter); err != nil {
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
		articleURL := fmt.Sprintf("%s/%s", channelLink, path.Join(dname, bname))
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
			pubDate = val.(string)
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
		pubDate = dt.Format(time.RFC1123)
		item := new(rss2.Item)
		item.GUID = articleURL
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

func fmtPandocError(err error) error {
	if err == nil {
		return nil
	}
	parts := []string{
		"Pandoc error (see https://pandoc.org), ",
		fmt.Sprintf("%s", err),
	}

	// Provide context where error is related to the PATH
	if strings.Contains(parts[1], "PATH") {
		parts = append(parts, fmt.Sprintf("PATH is %q", os.Getenv("PATH")))
	}
	return fmt.Errorf("%s", strings.Join(parts, "\n  "))
}

// Return the Pandoc version that will be used when calling Pandoc.
func GetPandocVersion() (string, error) {
	var (
		out, eOut bytes.Buffer
	)
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return "", fmtPandocError(err)
	}
	cmd := exec.Command(pandoc, "--version")
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return "", fmtPandocError(err)
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", pandoc, eOut.String())
	}
	return out.String(), fmtPandocError(err)
}

// pandocProcessor accepts an array of bytes as input and returns
// a `pandoc -f {From} -t html` output of an array if
// bytes and error.
func pandocProcessor(input []byte, from string, to string) ([]byte, error) {
	var (
		out, eOut bytes.Buffer
	)

	if from == "" {
		from = PandocFrom
	}
	if to == "" {
		to = PandocTo
	}
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, fmtPandocError(err)
	}
	options := []string{}
	if from != "" {
		options = append(options, "-f", from)
	}
	if to != "" {
		options = append(options, "-t", to)
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return nil, fmtPandocError(err)
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", pandoc, eOut.String())
	}
	return out.Bytes(), fmtPandocError(err)
}

// MakePandoc resolves key/value map rendering metadata suitable for processing with pandoc along with template information
// rendering and returns an error if something goes wrong
func MakePandoc(wr io.Writer, templateName string, keyValues map[string]string) error {
	var (
		out, eOut bytes.Buffer
		options   []string
	)

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return fmtPandocError(err)
	}
	data, err := ResolveData(keyValues)
	if err != nil {
		return fmt.Errorf("Data resolution error: %s", err)
	}
	// NOTE: If a template is not provided (empty string) then
	// see is one is specified in the metadata
	if templateName == "" {
		if val, ok := data["template"]; ok == true {
			templateName = val.(string)
		}
	}
	// NOTE: Pandocs default template expects content to be called $body$.
	// we need to remap from data["content"] to data["body"] otherwise
	// we need to look in data to see if a template was specified.
	if templateName == "" {
		if val, ok := data["content"]; ok == true {
			delete(data, "content")
			data["body"] = val
		}
	}
	// NOTE: when using a template, title metadata is required.
	if _, ok := data["title"]; !ok {
		// Insert a title to prevent warning.
		data["title"] = "..."
	}

	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Marshal error, %q", err)
	}
	metadata, err := ioutil.TempFile(".", "pandoc.*.json")
	if err != nil {
		return fmt.Errorf("Cannot create temp metadata file, %s", err)
	}
	if _, err := metadata.Write(src); err != nil {
		return fmt.Errorf("Write error, %q", err)
	}
	defer os.Remove(metadata.Name())
	// Check if document has front matter, split and write to temp files.
	options = []string{}
	if PandocFrom != "" {
		options = append(options, "-f", PandocFrom)
	}
	if PandocTo != "" {
		options = append(options, "-t", PandocTo)
	}
	options = append(options, "--metadata-file", metadata.Name())
	if templateName != "" {
		options = append(options, []string{"--template", templateName}...)
	} else {
		options = append(options, "--standalone")
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return fmtPandocError(err)
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", pandoc, eOut.String())
	}
	wr.Write(out.Bytes())
	return fmtPandocError(err)
}

// MakePandocString resolves key/value map rendering metadata suitable for processing with pandoc along with template information
// rendering and returns an error if something goes wrong
func MakePandocString(tmplSrc string, keyValues map[string]string) (string, error) {
	var (
		out, eOut bytes.Buffer
		options   []string
	)

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return "", fmtPandocError(err)
	}
	data, err := ResolveData(keyValues)
	if err != nil {
		return "", fmt.Errorf("Data resolution error: %s", err)
	}

	src, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Marshal error, %q", err)
	}
	metadata, err := ioutil.TempFile(".", "pandoc.*.json")
	if err != nil {
		return "", fmt.Errorf("Cannot create temp metadata file, %s", err)
	}
	if _, err := metadata.Write(src); err != nil {
		return "", fmt.Errorf("Write error, %q", err)
	}
	defer os.Remove(metadata.Name())

	options = []string{}
	if PandocFrom != "" {
		options = append(options, "-f", PandocFrom)
	}
	if PandocTo != "" {
		options = append(options, "-t", PandocTo)
	}
	options = append(options, "--metadata-file", metadata.Name())
	if tmplSrc != "" {
		// Pandoc expects to read the template from disc so write
		// out to a temp file.
		// Check if document has front matter, split and write to temp files.
		template, err := ioutil.TempFile(".", "pandoc.*.tmpl")
		if err != nil {
			return "", fmt.Errorf("Cannot create temp template file, %s", err)
		}
		if _, err := template.Write([]byte(tmplSrc)); err != nil {
			return "", fmt.Errorf("Write error, %q", err)
		}
		defer os.Remove(template.Name())
		options = append(options, []string{"--template", template.Name()}...)
	} else {
		options = append(options, "--standalone")
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return "", fmtPandocError(err)
	}
	if eOut.Len() > 0 {
		return "", fmt.Errorf("%q warns, %s", pandoc, eOut.String())
	}
	return fmt.Sprintf("%s", out.Bytes()), nil
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
//     src, _ := ioutil.ReadFile("post.md")
//     opening := mkpage.OpenParagraphs(fmt.Sprintf("%s", src, "\n\n"), 2)
//
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