package rss

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	// My packages
	"github.com/rsdoiel/pdtk"
	"github.com/rsdoiel/pdtk/blogit"
	"github.com/rsdoiel/pdtk/help"
)

var (
	// Standard options
	showHelp bool

	// App specific options
	excludeList        string
	atomLink           string
	baseURL            string
	channelLanguage    string
	channelTitle       string
	channelDescription string
	channelLink        string
	channelGenerator   string
	channelPubDate     string
	channelBuildDate   string
	channelCopyright   string
	channelCategory    string
	bylineExp          string
	titleExp           string
	dateExp            string
)

func usage(appName string, verb string) string {
	return help.Render(appName, verb, helpText)
}

func RunRSS(appName string, verb string, options []string) ([]byte, error) {

	flagSet := flag.NewFlagSet(appName+":"+verb, flag.ExitOnError)

	// Standard options
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")

	// App specific options
	flagSet.StringVar(&atomLink, "atom-link", "", "set atom:link href")
	flagSet.StringVar(&baseURL, "base-url", "", "set site base url for links")
	flagSet.StringVar(&excludeList, "e", "", "A colon delimited list of path exclusions")
	flagSet.StringVar(&channelLanguage, "channel-language", "", "Language, e.g. en-ca")
	flagSet.StringVar(&channelTitle, "channel-title", "", "Title of channel")
	flagSet.StringVar(&channelDescription, "channel-description", "", "Description of channel")
	flagSet.StringVar(&channelLink, "channel-link", "", "link to channel")
	flagSet.StringVar(&channelGenerator, "channel-generator", "", "Name of RSS generator")
	flagSet.StringVar(&channelPubDate, "channel-pubdate", "", "Pub Date for channel (e.g. 2006-01-02 15:04:05 -0700)")
	flagSet.StringVar(&channelBuildDate, "channel-builddate", "", "Build Date for channel (e.g. 2006-01-02 15:04:05 -0700)")
	flagSet.StringVar(&channelCopyright, "channel-copyright", "", "Copyright for channel")
	flagSet.StringVar(&channelCategory, "channel-category", "", "category for channel")
	flagSet.StringVar(&dateExp, "date-format", DateExp, "set date regexp")
	flagSet.StringVar(&titleExp, "title", TitleExp, "set title regexp")
	flagSet.StringVar(&bylineExp, "byline", BylineExp, "set byline regexp")

	flagSet.Parse(options)
	args := flagSet.Args()

	// Setup IO
	var err error

	// Process options
	if showHelp {
		usage(appName, verb)
		os.Exit(0)
	}

	if len(channelTitle) == 0 {
		channelTitle = `A website`
	}
	if len(channelDescription) == 0 {
		channelDescription = `A collection of web pages`
	}
	if len(channelLink) == 0 {
		channelLink = `http://localhost:8000`
	}

	// Setup the Channel metadata for feed.
	feed := new(RSS2)
	feed.Version = "2.0"
	feed.Title = channelTitle
	feed.Description = channelDescription
	feed.Link = channelLink
	feed.AtomNameSpace = "http://www.w3.org/2005/Atom"
	if len(channelLanguage) > 0 {
		feed.Language = channelLanguage
	}
	if len(channelCopyright) > 0 {
		feed.Copyright = channelCopyright
	}
	if len(channelCategory) > 0 {
		feed.Category = channelCategory
	}
	if len(channelGenerator) == 0 {
		feed.Generator = fmt.Sprintf("%s %s %s", appName, verb, pdtk.Version)
	} else {
		feed.Generator = channelGenerator
	}
	now := time.Now()
	if len(channelPubDate) == 0 {
		// RSS spec shows RFC 1123 dates
		// Validators indicate the RFC-822, note UTC isn't list on RFC-822
		//feed.PubDate = now.Format(time.RFC822Z)
		feed.PubDate = now.Format(time.RFC1123Z)
	} else {
		dt, err := NormalizeDate(channelPubDate)
		if err != nil {
			return nil, fmt.Errorf("Can't parse %q, %s\n", channelPubDate, err)
		}
		// RSS spec shows RFC 1123 dates
		// Validators indicate the RFC-822 and "UTC" isn't that.
		//feed.PubDate = dt.Format(time.RFC822Z)
		feed.PubDate = dt.Format(time.RFC1123Z)
	}
	if len(channelBuildDate) == 0 {
		// RSS spec shows RFC 1123 dates
		// Validators indicate the RFC-822
		//feed.LastBuildDate = now.Format(time.RFC822Z)
		feed.LastBuildDate = now.Format(time.RFC1123Z)
	} else {
		dt, err := NormalizeDate(channelBuildDate)
		if err != nil {
			return nil, fmt.Errorf("Can't parse %q, %s\n", channelBuildDate, err)
		}
		// RSS spec shows RFC 1123 dates
		// Validators indicate the RFC-822
		//feed.LastBuildDate = dt.Format(time.RFC822Z)
		feed.LastBuildDate = dt.Format(time.RFC1123Z)
	}

	// Process command line parameters
	htdocs := "."
	if len(args) > 0 {
		htdocs = args[0]
	}
	blogJSON := path.Join(htdocs, "blog.json")
	if _, err := os.Stat(blogJSON); os.IsNotExist(err) {

		err = WalkRSS(feed, htdocs, baseURL, excludeList, titleExp, bylineExp, dateExp)
	} else {
		blog := new(blogit.BlogMeta)
		src, err := ioutil.ReadFile(blogJSON)
		if err != nil {
			return nil, fmt.Errorf("Reading %q, %s\n", blogJSON, err)
		}
		if err := json.Unmarshal(src, &blog); err != nil {
			return nil, fmt.Errorf("Unmashal, %s\n", err)
		}
		if blog.BaseURL == "" {
			blog.BaseURL = baseURL
		}
		err = BlogMetaToRSS(blog, feed)
	}
	if err != nil {
		return nil, err
	}
	if atomLink != "" {
		feed.AtomLink = new(AtomLink)
		feed.AtomLink.HRef = atomLink
		feed.AtomLink.Rel = "self"
		feed.AtomLink.Type = "application/rss+xml"
	}

	// Marshal RSS2 and render output
	src, err := xml.MarshalIndent(feed, "", "    ")
	if err != nil {
		return nil, err
	}

	txt := strings.ReplaceAll(fmt.Sprintf(`%s%s`, xml.Header, src), "></atom:link>", "/>")
	return []byte(txt), nil
}
