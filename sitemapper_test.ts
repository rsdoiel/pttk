import { assert } from "@std/assert";
import { executeSitemapper } from "./sitemapper.ts";
import { readConfig, writeConfig, Config } from "./config.ts";
import { ensureDir } from "@std/fs";
import { join } from "@std/path";

async function setupTestDir(basePath: string) {
  await ensureDir(basePath);
  await writeConfig(join(basePath, "site.yaml"), {
    basePath,
    baseURL: "http://localhost:8000",
    views: "./views",
    viewsPartial: "./views/partials",
    maxUrlsPerSitemap: 50000,
    changeFrequency: "monthly",
    priority: "0.5",
    author: "Test Author",
    rssTitle: "Test Blog",
    rssDescription: "This is a test blog feed.",
    rssLink: "http://localhost:8000",
    series: [],
  } as Config);
}

async function cleanupTestDir(basePath: string) {
  await Deno.remove(basePath, { recursive: true });
}

Deno.test("sitemapper test", async () => {
  const basePath = "./test_sitemapper";
  await setupTestDir(basePath);

  await executeSitemapper(join(basePath, "site.yaml"));

  const sitemapIndexPath = join(basePath, "sitemap.xml");
  const sitemapIndexContent = await Deno.readTextFile(sitemapIndexPath);
  assert(sitemapIndexContent.includes("<sitemapindex"));

  await cleanupTestDir(basePath);
});
