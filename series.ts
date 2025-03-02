import { walk, ensureDir } from "@std/fs";
import { join } from "@std/path";
import { parse as parseYaml } from "@std/yaml";
import { readConfig, Config, SeriesConfig } from "./config.ts";
import { FrontMatter } from './frontmatter.ts';


interface RSSItem {
  title: string;
  link: string;
  description?: string;
  pubDate: string;
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

function generateSeriesTOC(items: { title: string, link: string }[], series: SeriesConfig): string {
  const tocContent = `
    # ${series.name}

    ${series.description}

    ${items.map(item => `- [${item.title}](${item.link})`).join('\n')}
  `;
  return tocContent;
}

function generateRSS(items: RSSItem[], config: Config, series: SeriesConfig): string {
  const rssContent = `
    <rss version="2.0">
      <channel>
        <title>${series.name}</title>
        <link>${config.rssLink}/${series.short_name}.xml</link>
        <description>${series.description}</description>
        ${items.map(item => `
          <item>
            <title>${item.title}</title>
            <link>${item.link}</link>
            ${item.description ? `<description>${item.description}</description>` : ''}
            <pubDate>${new Date(item.pubDate).toUTCString()}</pubDate>
          </item>
        `).join('')}
      </channel>
    </rss>
  `;
  return rssContent;
}

export async function executeSeries(configPath: string) {
  const config = await readConfig(configPath);
  const markdownFiles = await findMarkdownFiles(join(config.basePath, "blog"));

  for (const series of config.series) {
    const seriesItems: { title: string, link: string }[] = [];
    const rssItems: RSSItem[] = [];

    for (const file of markdownFiles) {
      const frontMatter = await parseFrontMatter(file);
      if (frontMatter.series === series.short_name) {
        const relativePath = file.replace(join(config.basePath, "blog"), "").replace(/\\/g, "/");
        const link = join(config.baseURL, relativePath.replace(".md", ".html"));
        const title: string = frontMatter.title || "Untitled";
	const abstract: string = (frontMatter.abstract === undefined) ? '' : frontMatter.abstract;
	const dateModified: string = (frontMatter.dateModified === undefined) ? (new Date()).toUTCString() : frontMatter.dateModified;

        seriesItems.push({ title, link });
        rssItems.push({ title, link, description: abstract || '', pubDate: dateModified });
      }
    }

    // Sort items by series_no and title
    seriesItems.sort((a, b) => {
      if (a.title < b.title) return -1;
      if (a.title > b.title) return 1;
      return 0;
    });

    const seriesDir = join(config.basePath, "series");
    await ensureDir(seriesDir);

    const tocContent = generateSeriesTOC(seriesItems, series);
    const tocPath = join(seriesDir, `${series.short_name}.md`);
    await Deno.writeTextFile(tocPath, tocContent);
    console.log(`Series TOC written to ${tocPath}`);

    const rssContent = generateRSS(rssItems, config, series);
    const rssPath = join(config.basePath, `${series.short_name}.xml`);
    await Deno.writeTextFile(rssPath, rssContent);
    console.log(`Series RSS feed written to ${rssPath}`);
  }
}
