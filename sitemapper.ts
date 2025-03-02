import { walk, ensureDir } from "@std/fs";
import { join } from "@std/path";
import { readConfig, Config } from "./config.ts";

export async function findMarkdownFiles(dirPath: string): Promise<string[]> {
  const files: string[] = [];
  for await (const entry of walk(dirPath)) {
    if (entry.path.endsWith(".md")) {
      files.push(entry.path);
    }
  }
  return files;
}

export async function generateSitemap(config: Config): Promise<{ index: number, content: string }[]> {
  const markdownFiles = await findMarkdownFiles(config.basePath);
  const markdownLinks = markdownFiles.map((file) => {
    const relativePath = file.replace(config.basePath, "").replace(/\\/g, "/");
    return join(config.baseURL, relativePath.replace(".md", ".html"));
  });

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

  for (let i = 0; i < markdownLinks.length; i++) {
    const link = markdownLinks[i];
    sitemapContent += `
      <url>
        <loc>${link}</loc>
        <changefreq>${config.changeFrequency}</changefreq>
        <priority>${config.priority}</priority>
      </url>
    `;

    if ((i + 1) % config.maxUrlsPerSitemap === 0 || i === markdownLinks.length - 1) {
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

export function generateSitemapIndex(config: Config, sitemapFiles: string[]): string {
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
  const sitemapDir = join(config.basePath, "sitemaps");

  await ensureDir(sitemapDir);

  const sitemaps = await generateSitemap(config);
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
