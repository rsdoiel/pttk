import { walk } from "@std/fs";
import { join } from "@std/path";
import { parse as parseYaml } from "@std/yaml";
import { readConfig, Config } from "./config.ts";

interface FrontMatter {
  title?: string;
  author: string;
  dateCreated: string;
  dateModified: string;
  abstract?: string;
  series?: string;
  series_no?: number;
  enclosureUrl?: string;
  enclosureLength?: string;
  enclosureType?: string;
}

interface RSSItem {
  title: string;
  link: string;
  description?: string;
  pubDate: string;
  enclosure?: {
    url: string;
    length?: string;
    type: string;
  };
}

async function findMarkdownFiles(dirPath: string): Promise<string[]> {
  const files: string[] = [];
  for await (const entry of walk(dirPath)) {
    if (entry.path.endsWith(".md")) {
      files.push(entry.path);
    }
  }
  return files;
}

async function parseFrontMatter(filePath: string): Promise<FrontMatter> {
  const fileContent = await Deno.readTextFile(filePath);
  const [frontMatterSection] = fileContent.split("---\n");
  return parseYaml(frontMatterSection.replace("---", "").trim()) as FrontMatter;
}

function generateRSS(items: RSSItem[], config: Config): string {
  const rssContent = `
    <rss version="2.0">
      <channel>
        <title>${config.rssTitle}</title>
        <link>${config.rssLink}</link>
        <description>${config.rssDescription}</description>
        ${items.map(item => `
          <item>
            <title>${item.title}</title>
            <link>${item.link}</link>
            ${item.description ? `<description>${item.description}</description>` : ''}
            <pubDate>${new Date(item.pubDate).toUTCString()}</pubDate>
            ${item.enclosure ? `
            <enclosure url="${item.enclosure.url}" length="${item.enclosure.length || ''}" type="${item.enclosure.type}" />
            ` : ''}
          </item>
        `).join('')}
      </channel>
    </rss>
  `;
  return rssContent;
}

export async function executeRSSGenerator(configPath: string, outputPath: string) {
  const config = await readConfig(configPath);
  const markdownFiles = await findMarkdownFiles(join(config.basePath, "blog"));

  const rssItems: RSSItem[] = [];

  for (const file of markdownFiles) {
    const frontMatter = await parseFrontMatter(file);
    const relativePath = file.replace(join(config.basePath, "blog"), "").replace(/\\/g, "/");
    const link = join(config.baseURL, relativePath.replace(".md", ".html"));
    const title = frontMatter.title || "Untitled";

    rssItems.push({
      title,
      link,
      description: frontMatter.abstract,
      pubDate: frontMatter.dateModified,
      enclosure: frontMatter.enclosureUrl ? {
        url: frontMatter.enclosureUrl,
        type: frontMatter.enclosureType || "audio/mpeg",
        length: frontMatter.enclosureLength,
      } : undefined,
    });
  }

  const rssContent = generateRSS(rssItems, config);
  await Deno.writeTextFile(outputPath, rssContent);
  console.log(`RSS feed written to ${outputPath}`);
}
