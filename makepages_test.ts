import { assert } from "@std/assert";
import { executeMakepages } from "./makepages.ts";
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

Deno.test("makepages test", async () => {
  const basePath = "./test_makepages";
  await setupTestDir(basePath);

  const markdownContent = `
---
title: Test Page
author: Test Author
dateCreated: 2023-10-05
dateModified: 2023-10-05
---
# Test Page
This is a test page.
`;

  const markdownPath = join(basePath, "test_page.md");
  await Deno.writeTextFile(markdownPath, markdownContent);

  await executeMakepages(join(basePath, "site.yaml"));

  const htmlPath = join(basePath, "test_page.html");
  const htmlContent = await Deno.readTextFile(htmlPath);
  assert(htmlContent.includes("Test Page"));

  await cleanupTestDir(basePath);
});
