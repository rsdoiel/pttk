import { assert } from "@std/assert";
import { readConfig, writeConfig, Config } from "./config.ts";
import { ensureDir } from "@std/fs";
import { join } from "@std/path";

async function setupTestDir(basePath: string) {
  await ensureDir(basePath);
}

async function cleanupTestDir(basePath: string) {
  await Deno.remove(basePath, { recursive: true });
}

Deno.test("config test", async () => {
  const basePath = "./test_config";
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
    series: [
      { short_name: "test_series", name: "Test Series", description: "This is a test series." },
    ],
  };

  await writeConfig(join(basePath, "site.yaml"), config);

  const readConfigResult = await readConfig(join(basePath, "site.yaml"));
  assert(readConfigResult.basePath === basePath);
  assert(readConfigResult.author === "Test Author");
  assert(readConfigResult.series.length === 1);
  assert(readConfigResult.series[0].short_name === "test_series");

  await cleanupTestDir(basePath);
});

Deno.test("default config test", async () => {
  const basePath = "./test_default_config";
  await setupTestDir(basePath);

  await writeConfig(join(basePath, "site.yaml"), {
    author: '',
    basePath,
    baseURL: "http://localhost:8000",
    views: "./views",
    viewsPartial: "./views/partials",
    maxUrlsPerSitemap: 50000,
    changeFrequency: "monthly",
    priority: "0.5",
    rssTitle: "Test Blog",
    rssDescription: "This is a test blog feed.",
    rssLink: "http://localhost:8000",
    series: [],
  });

  const readConfigResult = await readConfig(join(basePath, "site.yaml"));
  assert(readConfigResult.basePath === basePath);
  assert(readConfigResult.author === "");
  assert(readConfigResult.series.length === 0);

  await cleanupTestDir(basePath);
});
