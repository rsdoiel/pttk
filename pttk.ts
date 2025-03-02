import { parse } from "@std/flags";
import { initConfig, interactiveConfig } from "./config.ts";
import { executeSitemapper } from "./sitemapper.ts";
import { executeMakepages } from "./makepages.ts";
import { executeBlogit } from "./blogit.ts";
import { executeRSSGenerator } from "./rss_generator.ts";
import { executeSeries } from "./series.ts";

function displayHelp() {
  console.log(`
    PT Toolkit CLI

    Usage:
      pttk.ts <action> [options]

    Actions:
      help       Display this help message
      init       Create a default configuration file
      config     Interactively set configuration values
      sitemapper Generate sitemap files based on the configuration
      makepages   Generate HTML files from Markdown based on the configuration
      blogit      Organize blog posts and additional files based on the configuration
      rss         Generate an RSS feed from blog posts
      series      Generate series TOC and RSS feeds based on the configuration

    Options:
      --config <path>  Specify the path to the configuration file (default: site.yaml)
      --output <path>  Specify the output path for the RSS feed (default: index.xml)
      --filepath <path> Specify the file path for blogit action (can be used multiple times)
      --date <date>    Specify the date for blogit action (format: YYYY-MM-DD)
  `);
}

async function main() {
  const args = parse(Deno.args, {
    string: ["config", "filepath", "date", "output"],
    boolean: ["filepath"],
    default: { config: "site.yaml", output: "index.xml" },
  });

  const action = args._[0];

  if (!action) {
    displayHelp();
    return;
  }

  switch (action) {
    case "help":
      displayHelp();
      break;

    case "init":
      await initConfig(args.config);
      console.log(`Configuration file generated at ${args.config}`);
      break;

    case "config":
      await interactiveConfig(args.config);
      break;

    case "sitemapper":
      await executeSitemapper(args.config);
      break;

    case "makepages":
      await executeMakepages(args.config);
      break;

    case "blogit":
      const date = new Date(args.date);
      if (isNaN(date.getTime())) {
        console.error("Invalid date format. Please use YYYY-MM-DD.");
        return;
      }
      const filepaths = args.filepath ? Array.isArray(args.filepath) ? args.filepath : [args.filepath] : [];
      await executeBlogit(args.config, filepaths, date);
      break;

    case "rss":
      await executeRSSGenerator(args.config, args.output);
      break;

    case "series":
      await executeSeries(args.config);
      break;

    default:
      console.error(`Unknown action: ${action}`);
      displayHelp();
      break;
  }
}

if (import.meta.main) {
  main().catch(console.error);
}
