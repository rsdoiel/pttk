package rss

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	// My packages
	"github.com/rsdoiel/pdtk"
	"github.com/rsdoiel/pdtk/blogit"
	"github.com/rsdoiel/pdtk/help"

	// Caltech package
	"github.com/caltechlibrary/rss2"
)

var (
	// Standard options
	showHelp    bool
	inputFName  string
	outputFName string

	// App specific options
	excludeList        string
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
	flagSet.StringVar(&inputFName, "i", "", "set input filename")
	flagSet.StringVar(&outputFName, "o", "", "set output filename")

	// App specific options
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
	flagSet.StringVar(&dateExp, "d,date-format", DateExp, "set date regexp")
	flagSet.StringVar(&titleExp, "t,title", TitleExp, "set title regexp")
	flagSet.StringVar(&bylineExp, "b,byline", BylineExp, "set byline regexp")

	flagSet.Parse(options)
	args := flagSet.Args()

	// Setup IO
	var err error

	in := os.Stdin

	if inputFName != "" {
		in, err = os.Open(inputFName)
		if err != nil {
			return nil, err
		}
		defer in.Close()
	}

	if outputFName != "" {
		out, err := os.Create(outputFName)
		if err != nil {
			return nil, err
		}
		defer out.Close()
	}

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
	feed := new(rss2.RSS2)
	feed.Version = "2.0"
	feed.Title = channelTitle
	feed.Description = channelDescription
	feed.Link = channelLink
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
		feed.Generator = pdtk.Version
	} else {
		feed.Generator = channelGenerator
	}
	now := time.Now()
	if len(channelPubDate) == 0 {
		// RSS spec shows RTF 1123 dates
		//feed.PubDate = now.Format(time.RFC822Z)
		feed.PubDate = now.Format(time.RFC1123)
	} else {
		dt, err := NormalizeDate(channelPubDate)
		if err != nil {
			return nil, fmt.Errorf("Can't parse %q, %s\n", channelPubDate, err)
		}
		feed.PubDate = dt.Format(time.RFC1123)
	}
	if len(channelBuildDate) == 0 {
		// RSS spec shows RTF 1123 dates
		//feed.LastBuildDate = now.Format(time.RFC822Z)
		feed.LastBuildDate = now.Format(time.RFC1123)
	} else {
		dt, err := NormalizeDate(channelBuildDate)
		if err != nil {
			return nil, fmt.Errorf("Can't parse %q, %s\n", channelBuildDate, err)
		}
		feed.LastBuildDate = dt.Format(time.RFC1123)
	}

	// Process command line parameters
	htdocs := "."
	rssPath := ""
	if len(args) > 0 {
		htdocs = args[0]
	}
	if len(args) > 1 {
		rssPath = args[1]
	}
	blogJSON := path.Join(htdocs, "blog.json")
	if _, err := os.Stat(blogJSON); os.IsNotExist(err) {
		err = WalkRSS(feed, htdocs, excludeList, titleExp, bylineExp, dateExp)
	} else {
		blog := new(blogit.BlogMeta)
		if src, err := ioutil.ReadFile(blogJSON); err != nil {
			return nil, fmt.Errorf("Reading %q, %s\n", blogJSON, err)
		} else {
			if err := json.Unmarshal(src, &blog); err != nil {
				return nil, fmt.Errorf("Unmashal, %s\n", err)
			}
		}
		err = BlogMetaToRSS(blog, feed)
	}
	if err != nil {
		return nil, err
	}

	// Marshal RSS2 and render output
	src, err := xml.MarshalIndent(feed, "", "    ")
	if err != nil {
		return nil, err
	}

	txt := fmt.Sprintf(`<?xml version="1.0"?>
%s`, src)
	if len(rssPath) > 0 {
		err = ioutil.WriteFile(rssPath, []byte(txt), 0664)
		if err != nil {
			return nil, err
		}
	}
	return []byte(txt), nil
}
