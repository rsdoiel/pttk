
# Deno+TypeScript Port

## 1. Project Structure

Create a directory structure for your project:

```
pttk/
│
├── config.ts
├── frontmatter.ts
├── helptext.ts
├── makepages.ts
├── rss_generator.ts
├── series.ts
├── sitemapper.ts
├── blogit.ts
├── pttk.ts
├── version.ts
├── deno.json
├── sitemapper_test.ts
├── rss_generator_test.ts
├── blogit_test.ts
├── series_test.ts
├── makepages_test.ts
└── site.yaml
```

## 2. Configuration Module (`config.ts`)

- Manages reading, writing, and initializing the configuration file (`site.yaml`).
- Includes settings for base path, base URL, views, RSS, series, and author.

## 3. Front Matter Module (`frontmatter.ts`)

- Handles parsing and updating the front matter of Markdown files.
- Includes fields like title, author, dateCreated, dateModified, draft, byline, abstract, keywords, series, series_no, enclosureUrl, enclosureLength, and enclosureType.

## 4. Help Text Module (`helptext.ts`)

- Formats and displays help text for each action.
- Uses Markdown notation to structure the help text in a man page style.

## 5. Version Module (`version.ts`)

- Contains version information, release date, release hash, and license text.

## 6. Sitemapper Module (`sitemapper.ts`)

- Generates sitemap files based on the configuration.
- Walks the directory structure to find Markdown files and generates sitemap XML files.

## 7. RSS Generator Module (`rss_generator.ts`)

- Generates an RSS feed from blog posts.
- Supports podcast enclosures and other metadata.

## 8. Blogit Module (`blogit.ts`)

- Organizes blog posts and additional files based on the configuration.
- Handles front matter, including optional fields and enclosure metadata.

## 9. Series Module (`series.ts`)

- Generates a table of contents and RSS feed for each series.
- Walks the directory structure to find Markdown files and generates series-specific content.

## 10. Makepages Module (`makepages.ts`)

- Converts Markdown files to HTML using HandlebarsJS templates.
- Walks the directory structure to find Markdown files and generates HTML files.

## 11. CLI Entry Point (`pttk.ts`)

- Provides a unified CLI interface for managing blog content.
- Supports actions like `init`, `config`, `sitemapper`, `makepages`, `blogit`, `rss`, `series`, `version`, and `license`.
- Uses the help text module to display help information.

## 12. Test Modules

- Individual test files for each module (`sitemapper_test.ts`, `rss_generator_test.ts`, `blogit_test.ts`, `series_test.ts`, `makepages_test.ts`).
- Uses Deno's testing framework to verify the functionality of each module.

## 13. Deno Configuration (`deno.json`)

- Includes tasks for building, cross-compiling, testing, and running the CLI tool.
- Allows running tests individually and all together.

## 14. Configuration File (`site.yaml`)

- Contains the shared configuration settings for the CLI tool.
- Includes settings for base path, base URL, views, RSS, series, and author.

## Running the CLI

- Use Deno tasks to build, test, and run the CLI tool.
- Example commands:
  - Build: `deno task build`
  - Test: `deno task test:all`
  - Run actions: `deno run --allow-read --allow-write pttk.ts <action>`

This structure provides a comprehensive and modular solution for managing blog content using Deno and TypeScript.
