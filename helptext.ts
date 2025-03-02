import { version, releaseDate, releaseHash } from "./version.ts";

const helpText = `
PT Toolkit CLI

NAME
    pttk - A comprehensive toolkit for managing blog content.

SYNOPSIS
    pttk <action> [options]

DESCRIPTION
    The PT Toolkit CLI provides a set of tools for managing blog content, including generating sitemaps, converting Markdown to HTML, organizing blog posts, generating RSS feeds, and managing series.

ACTIONS
    help
        Display this help message.

    init
        Create a default configuration file.

    config
        Interactively set configuration values.

    sitemapper
        Generate sitemap files based on the configuration.

    makepages
        Generate HTML files from Markdown based on the configuration.

    blogit
        Organize blog posts and additional files based on the configuration.

    rss
        Generate an RSS feed from blog posts.

    series
        Generate series TOC and RSS feeds based on the configuration.

    version
        Display the application version information.

    license
        Display the license text.

OPTIONS
    --config <path>
        Specify the path to the configuration file (default: site.yaml).

    --output <path>
        Specify the output path for the RSS feed (default: index.xml).

    --filepath <path>
        Specify the file path for the blogit action (can be used multiple times).

    --date <date>
        Specify the date for the blogit action (format: YYYY-MM-DD).

VERSION
    ${version}

RELEASE DATE
    ${releaseDate}

RELEASE HASH
    ${releaseHash}
`;

export function fmtHelp(actionName: string): string {
  const sections = helpText.split("\n\n");
  const actionSection = sections.find(section => section.startsWith(`    ${actionName}`));

  if (actionSection) {
    return `
${sections[0]}

${sections[1]}

${actionSection}

${sections[sections.length - 3]}

${sections[sections.length - 2]}

${sections[sections.length - 1]}
    `.trim();
  }

  return helpText.trim();
}
