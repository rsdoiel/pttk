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
