import { assert } from "@std/assert";
import { executeBlogit } from "./blogit.ts";
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

Deno.test("blogit test", async () => {
  const basePath = "./test_blogit";
  await setupTestDir(basePath);

  const markdownContent = `
---
title: Test Post
author: Test Author
dateCreated: 2023-10-05
dateModified: 2023-10-05
---
# Test Post
This is a test post.
`;

  const markdownPath = join(basePath, "test_post.md");
  await Deno.writeTextFile(markdownPath, markdownContent);

  await executeBlogit(join(basePath, "site.yaml"), [markdownPath], new Date("2023-10-05"));

  const blogPath = join(basePath, "blog", "2023", "10", "05", "test_post.md");
  const blogContent = await Deno.readTextFile(blogPath);
  assert(blogContent.includes("Test Post"));

  await cleanupTestDir(basePath);
});
