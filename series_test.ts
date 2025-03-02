import { assert } from "@std/assert";
import { executeSeries } from "./series.ts";
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
    series: [
      { short_name: "test_series", name: "Test Series", description: "This is a test series." },
    ],
  } as Config);
}

async function cleanupTestDir(basePath: string) {
  await Deno.remove(basePath, { recursive: true });
}

Deno.test("series test", async () => {
  const basePath = "./test_series";
  await setupTestDir(basePath);

  const markdownContent = `
---
title: Test Series Post
author: Test Author
dateCreated: 2023-10-05
dateModified: 2023-10-05
series: test_series
series_no: 1
---
# Test Series Post
This is a test series post.
`;

  const markdownPath = join(basePath, "blog", "test_series_post.md");
  await ensureDir(join(basePath, "blog"));
  await Deno.writeTextFile(markdownPath, markdownContent);

  await executeSeries(join(basePath, "site.yaml"));

  const seriesTOCPath = join(basePath, "series", "test_series.md");
  const seriesTOCContent = await Deno.readTextFile(seriesTOCPath);
  assert(seriesTOCContent.includes("Test Series Post"));

  const seriesRSSPath = join(basePath, "test_series.xml");
  const seriesRSSContent = await Deno.readTextFile(seriesRSSPath);
  assert(seriesRSSContent.includes("<rss"));

  await cleanupTestDir(basePath);
});
