import { parse as parseYaml, stringify as stringifyYaml } from "@std/yaml";

export interface SeriesConfig {
  short_name: string;
  name: string;
  description: string;
}

export interface Config {
  basePath: string;
  baseURL: string;
  views: string;
  viewsPartial: string;
  maxUrlsPerSitemap: number;
  changeFrequency: string;
  priority: string;
  author: string;
  rssTitle: string;
  rssDescription: string;
  rssLink: string;
  series: SeriesConfig[];  // An array of SeriesConfig objects.
}

export async function readConfig(configPath: string): Promise<Config> {
  const configContent = await Deno.readTextFile(configPath);
  const config = parseYaml(configContent) as Config;

  // Set default author if not provided
  if (!config.author) {
    config.author = Deno.env.get("USER") || "Unknown Author";
  }

  // Set default RSS channel configuration if not provided
  if (!config.rssTitle) {
    config.rssTitle = "My Blog";
  }

  if (!config.rssDescription) {
    config.rssDescription = "This is my blog feed.";
  }

  if (!config.rssLink) {
    config.rssLink = config.baseURL;
  }

  // Set default series configuration if not provided
  if (!config.series) {
    config.series = [];
  }

  return config;
}

export async function writeConfig(configPath: string, config: Config): Promise<void> {
  const configContent = stringifyYaml(config);
  await Deno.writeTextFile(configPath, configContent);
}

export async function initConfig(configPath: string): Promise<void> {
  const initialConfig: Config = {
    basePath: ".",
    baseURL: "http://localhost:8000",
    views: "./views",
    viewsPartial: "./views/partials",
    maxUrlsPerSitemap: 50000,
    changeFrequency: "monthly",
    priority: "0.5",
    author: Deno.env.get("USER") || "Unknown Author",
    rssTitle: "My Blog",
    rssDescription: "This is my blog feed.",
    rssLink: "http://localhost:8000",
    series: [],
  };
  await writeConfig(configPath, initialConfig);
}

export async function interactiveConfig(configPath: string): Promise<void> {
  const config: Config = {
    basePath: prompt("Enter base path (default: .): ") || ".",
    baseURL: prompt("Enter base URL (default: http://localhost:8000): ") || "http://localhost:8000",
    views: prompt("Enter views directory (default: ./views): ") || "./views",
    viewsPartial: prompt("Enter views partial directory (default: ./views/partials): ") || "./views/partials",
    maxUrlsPerSitemap: parseInt(prompt("Enter max URLs per sitemap (default: 50000): ") || "50000"),
    changeFrequency: prompt("Enter change frequency (default: monthly): ") || "monthly",
    priority: prompt("Enter priority (default: 0.5): ") || "0.5",
    author: prompt("Enter author (default: current user): ") || Deno.env.get("USER") || "Unknown Author",
    rssTitle: prompt("Enter RSS title (default: My Blog): ") || "My Blog",
    rssDescription: prompt("Enter RSS description (default: This is my blog feed.): ") || "This is my blog feed.",
    rssLink: prompt("Enter RSS link (default: base URL): ") || "http://localhost:8000",
    series: [],
  };

  const addSeries = prompt("Do you want to add a series configuration? (yes/no): ") === "yes";
  if (addSeries) {
    while (true) {
      const short_name = prompt("Enter series short name: ") || '';
      const name = prompt("Enter series name (default: short name): ") || short_name;
      const description = prompt("Enter series description: ") || '';
      config.series.push({ short_name, name, description });
      const addAnother = prompt("Do you want to add another series configuration? (yes/no): ");
      if (addAnother !== "yes") break;
    }
  }

  await writeConfig(configPath, config);
  console.log(`Configuration file updated at ${configPath}`);
}
