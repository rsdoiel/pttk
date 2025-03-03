
import { parse as parseYaml, stringify as stringifyYaml } from "@std/yaml";
import { Config } from "./config.ts";

export interface FrontMatter {
  title?: string;
  author: string;
  dateCreated: string;
  dateModified: string;
  draft?: boolean;
  byline?: string;
  abstract?: string;
  keywords?: string[];
  series?: string;
  series_no?: number;
  enclosureUrl?: string;
  enclosureLength?: string;
  enclosureType?: string;
}

export async function parseFrontMatter(filePath: string, config: Config): Promise<{ frontMatter: FrontMatter, content: string }> {
  const fileContent = await Deno.readTextFile(filePath);
  const [frontMatterSection, ...markdownContent] = fileContent.split("---\n");

  let frontMatter: FrontMatter = parseYaml(frontMatterSection.replace("---", "").trim()) as FrontMatter;
  const content = markdownContent.join("---\n").trim();

  if (!frontMatter.title) {
    frontMatter.title = "Untitled";
  }

  if (!frontMatter.author) {
    frontMatter.author = config.author;
  }

  const now = new Date().toISOString().split('T')[0];
  if (!frontMatter.dateCreated) {
    frontMatter.dateCreated = now;
  }

  if (!frontMatter.dateModified) {
    frontMatter.dateModified = now;
  }

  if (frontMatter.draft === undefined) {
    frontMatter.draft = false;
  }

  if (!frontMatter.byline) {
    frontMatter.byline = `${frontMatter.author}, ${now}`;
  }

  if (!frontMatter.abstract) {
    frontMatter.abstract = "";
  }

  if (!frontMatter.keywords) {
    frontMatter.keywords = [];
  }

  if (!frontMatter.series) {
    frontMatter.series = "";
  }

  if (frontMatter.series_no === undefined) {
    frontMatter.series_no = 1;
  }

  return { frontMatter, content };
}

export async function updateFrontMatter(filePath: string, frontMatter: FrontMatter, content: string) {
  const updatedContent = `---\n${stringifyYaml(frontMatter)}---\n${content}`;
  await Deno.writeTextFile(filePath, updatedContent);
}
