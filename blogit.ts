import { ensureDir } from "@std/fs";
import { join, basename } from "@std/path";
import { readConfig } from "./config.ts";
import { parseFrontMatter, updateFrontMatter, FrontMatter } from "./frontmatter.ts";

async function createBlogDirectory(basePath: string, date: Date) {
  const year = date.getFullYear().toString();
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const day = date.getDate().toString().padStart(2, '0');

  const blogPath = join(basePath, "blog", year, month, day);
  await ensureDir(blogPath);

  return blogPath;
}

export async function executeBlogit(configPath: string, filepaths: string[], date: Date) {
  const config = await readConfig(configPath);

  const markdownFiles = filepaths.filter(file => file.endsWith(".md") || file.endsWith(".markdown"));
  if (markdownFiles.length === 0) {
    throw new Error("No Markdown file provided.");
  }

  for (const markdownFile of markdownFiles) {
    const { frontMatter, content } = await parseFrontMatter(markdownFile, config);
    await updateFrontMatter(markdownFile, frontMatter, content);

    const blogPath = await createBlogDirectory(config.basePath, date);

    const markdownFilename = basename(markdownFile);
    const markdownDestinationPath = join(blogPath, markdownFilename);
    await Deno.copyFile(markdownFile, markdownDestinationPath);
    console.log(`Markdown file copied to ${markdownDestinationPath}`);

    const additionalFiles = filepaths.filter(file => !markdownFiles.includes(file));
    for (const file of additionalFiles) {
      const filename = basename(file);
      const destinationPath = join(blogPath, filename);
      await Deno.copyFile(file, destinationPath);
      console.log(`Additional file copied to ${destinationPath}`);
    }

    if (frontMatter.enclosureUrl) {
      const enclosureFilename = basename(frontMatter.enclosureUrl);
      const enclosureDestinationPath = join(blogPath, enclosureFilename);
      await Deno.copyFile(frontMatter.enclosureUrl, enclosureDestinationPath);
      console.log(`Enclosure file copied to ${enclosureDestinationPath}`);
    }
  }
}
