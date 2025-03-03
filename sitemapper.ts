
import { walk, ensureDir } from "@std/fs";
import { join } from "@std/path";
import { parse as parseYaml } from '@std/yaml';
import { readConfig, Config } from "./config.ts";

interface RSSItem {
  title: string;
  link: string;
  description: string;
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

async function parseFrontMatter(filePath: string): Promise<{ frontMatter: any, content: string }> {
  const fileContent = await Deno.readTextFile(filePath);
  const [frontMatterSection, ...markdownContent] = fileContent.split("---\n");

  let frontMatter: any = parseYaml(frontMatterSection.replace("---", "").trim());
  const content = markdownContent.join("---\n").trim();

  return { frontMatter, content };
}

function generateSitemap(items: string[], config: Config): { index: number, content: string }[] {
  const sitemaps: { index: number, content: string }[] = [];
  let sitemapIndex = 1;
  let sitemapContent = `
    <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
      <url>
        <loc>${config.baseURL}</loc>
        <changefreq>${config.changeFrequency}</changefreq>
        <priority>${config.priority}</priority>
      </url>
  `;

  for (let i = 0; i < items.length; i++) {
    const link = items[i];
    sitemapContent += `
      <url>
        <loc>${link}</loc>
        <changefreq>${config.changeFrequency}</changefreq>
        <priority>${config.priority}</priority>
      </url>
    `;

    if ((i + 1) % config.maxUrlsPerSitemap === 0 || i === items.length - 1) {
      sitemapContent += "</urlset>";
      sitemaps.push({ index: sitemapIndex, content: sitemapContent });
      sitemapIndex++;
      sitemapContent = `
        <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
          <url>
            <loc>${config.baseURL}</loc>
            <changefreq>${config.changeFrequency}</changefreq>
            <priority>${config.priority}</priority>
          </url>
      `;
    }
  }

  return sitemaps;
}

function generateSitemapIndex(config: Config, sitemapFiles: string[]): string {
  let sitemapIndexContent = `
    <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  `;

  sitemapFiles.forEach((file) => {
    sitemapIndexContent += `
      <sitemap>
        <loc>${join(config.baseURL, file)}</loc>
      </sitemap>
    `;
  });

  sitemapIndexContent += "</sitemapindex>";
  return sitemapIndexContent;
}

export async function executeSitemapper(configPath: string) {
  const config = await readConfig(configPath);
  const markdownFiles = await findMarkdownFiles(config.basePath);

  const sitemapDir = join(config.basePath, "sitemaps");
  await ensureDir(sitemapDir);

  const sitemaps = await generateSitemap(
    markdownFiles.map(
      file => join(config.baseURL, 
        file.replace(config.basePath, "").replace(/\\/g, "/"))), config);
  const sitemapFiles: string[] = [];

  for (const sitemap of sitemaps) {
    const sitemapFilePath = join(sitemapDir, `sitemap-${sitemap.index}.xml`);
    await Deno.writeTextFile(sitemapFilePath, sitemap.content);
    sitemapFiles.push(sitemapFilePath.replace(config.basePath, "").replace(/\\/g, "/"));
    console.log(`Sitemap ${sitemap.index} written to ${sitemapFilePath}`);
  }

  const sitemapIndexContent = generateSitemapIndex(config, sitemapFiles);
  const sitemapIndexPath = join(config.basePath, "sitemap.xml");
  await Deno.writeTextFile(sitemapIndexPath, sitemapIndexContent);
  console.log(`Sitemap index written to ${sitemapIndexPath}`);
}
