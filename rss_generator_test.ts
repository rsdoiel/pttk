import { assert } from "@std/assert";
import { executeRSSGenerator } from "./rss_generator.ts";
import { writeConfig, Config } from "./config.ts";
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

Deno.test("rss_generator test", async () => {
  const basePath = "./test_rss_generator";
  await setupTestDir(basePath);

  await executeRSSGenerator(join(basePath, "site.yaml"), "index.xml");

  const rssPath = join(basePath, "index.xml");
  const rssContent = await Deno.readTextFile(rssPath);
  assert(rssContent.includes("<rss"));

  await cleanupTestDir(basePath);
});
