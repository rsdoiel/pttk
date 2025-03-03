import { assert } from "@std/assert";
import { parseFrontMatter, updateFrontMatter } from "./frontmatter.ts";
import { ensureDir } from "@std/fs";
import { join } from "@std/path";
import { Config } from "./config.ts";

async function setupTestDir(basePath: string) {
  await ensureDir(basePath);
}

async function cleanupTestDir(basePath: string) {
  await Deno.remove(basePath, { recursive: true });
}

Deno.test("parseFrontMatter test", async () => {
  const basePath = "./test_frontmatter";
  await setupTestDir(basePath);

  const config: Config = {
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
  };

  const markdownContent = `
---
title: Test Post
dateCreated: 2023-10-05
dateModified: 2023-10-05
---
# Test Post
This is a test post.
`;

  const markdownPath = join(basePath, "test_post.md");
  await Deno.writeTextFile(markdownPath, markdownContent);

  const { frontMatter, content } = await parseFrontMatter(markdownPath, config);

  assert(frontMatter.title === "Test Post");
  assert(frontMatter.author === "Test Author");
  assert(frontMatter.dateCreated === "2023-10-05");
  assert(frontMatter.dateModified === "2023-10-05");
  assert(content.includes("Test Post"));

  await cleanupTestDir(basePath);
});

Deno.test("updateFrontMatter test", async () => {
  const basePath = "./test_frontmatter";
  await setupTestDir(basePath);

  const config: Config = {
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
  };

  const markdownContent = `
---
title: Test Post
dateCreated: 2023-10-05
dateModified: 2023-10-05
---
# Test Post
This is a test post.
`;

  const markdownPath = join(basePath, "test_post.md");
  await Deno.writeTextFile(markdownPath, markdownContent);

  const { frontMatter, content } = await parseFrontMatter(markdownPath, config);
  frontMatter.title = "Updated Test Post";
  await updateFrontMatter(markdownPath, frontMatter, content);

  const updatedContent = await Deno.readTextFile(markdownPath);
  assert(updatedContent.includes("Updated Test Post"));

  await cleanupTestDir(basePath);
});

Deno.test("defaultFrontMatter test", async () => {
  const basePath = "./test_frontmatter";
  await setupTestDir(basePath);

  const config: Config = {
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
  };

  const markdownContent = `
# Test Post
This is a test post without front matter.
`;

  const markdownPath = join(basePath, "test_post.md");
  await Deno.writeTextFile(markdownPath, markdownContent);

  const { frontMatter, content } = await parseFrontMatter(markdownPath, config);

  assert(frontMatter.title === "Untitled");
  assert(frontMatter.author === "Test Author");
  assert(frontMatter.dateCreated !== "");
  assert(frontMatter.dateModified !== "");
  assert(content.includes("Test Post"));

  await cleanupTestDir(basePath);
});
