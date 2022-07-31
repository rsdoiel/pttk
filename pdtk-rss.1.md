% pdtk-rss(1) pdtk user manual
% R. S. Doiel
% July, 31, 2022

# NAME

pdtk rss PATH_TO_SITE

# SYNOPSIS

pdtk rss [OPTIONS] PATH_TO_SITE 

# DESCRIPTION

The rss renders an RSS file based on the content found in the
directory tree provided. If it encounters a "blog.json" file then
it'll use that file to generate feed content for that directory
and it's content otherwise it'll generate a feed backed on Markdown
front matter encountered in Markdown documents with corresponding
html file.

pdtk rss walks the file system to generate a RSS2 file. It assumes 
that the directory for HTDOCS is is the base directory containing 
subdirectories in the form of /YYYY/MM/DD/ARTICLE_HTML where 
YYYY/MM/DD (Year, Month, Day) corresponds to the publication date 
of ARTICLE_HTML.

If our htdocs folder is our document root and out blog is
htdocs/myblog.

    pdtk rss \
        -atom-link="http://blog.example.org/rss.xml" \
        -base-url="http://blog.example.org" \
        -channel-title="This Great Beyond" \
        -channel-description="Blog to save the world" \
        -channel-link="http://blog.example.org" \
        htdocs >htdocs/rss.xml

This would build an RSS 2 file in htdocs/rss.xml from the
articles in htdocs/myblog/YYYY/MM/DD.

# OPTIONS

What follows is are the options supported by the rss verb.

-atom-link string
: set atom:link href

-base-url string
: set site base url for links

-byline string
: set byline regexp (default "^[B|b]y\\s+(\\w|\\s|.)+[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]$")

-channel-builddate string
: Build Date for channel (e.g. 2006-01-02 15:04:05 -0700)

-channel-category string
: category for channel

-channel-copyright string
: Copyright for channel

-channel-description string
: Description of channel

-channel-generator string
: Name of RSS generator

-channel-language string
: Language, e.g. en-ca

-channel-link string
: link to channel

-channel-pubdate string
: Pub Date for channel (e.g. 2006-01-02 15:04:05 -0700)

-channel-title string
: Title of channel

-date-format string
: set date regexp (default "[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]")

-e string
: A colon delimited list of path exclusions

-help
: display rss help

-title string
: set title regexp (default "^#\\s+(\\w|\\s|.)+$")


# EXAMPLE

# SEE ALSO

pdtk website at https://rsdoiel.github.io/pdtk

The source code is avialable from https://github.com/rsdoiel/pdtk


