import { parse as parseYaml } from '@std/yaml';
import { walk, ensureDir } from "@std/fs";
import { join, dirname } from "@std/path";
import { readConfig } from "./config.ts";
import { micromark } from 'https://esm.sh/micromark@3';
import Handlebars from "npm:handlebars";

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

function convertMarkdownToHtml(src: string): string {
  return micromark(src);
}

async function renderHtml(templatePath: string, data: { content: string, frontMatter: any }): Promise<string> {
  const templateContent = await Deno.readTextFile(templatePath);
  const template = Handlebars.compile(templateContent);
  return template(data);
}

export async function executeMakepages(configPath: string) {
  const config = await readConfig(configPath);
  const markdownFiles = await findMarkdownFiles(config.basePath);

  for (const file of markdownFiles) {
    const { frontMatter, content } = await parseFrontMatter(file);
    const htmlContent = convertMarkdownToHtml(content);

    const relativePath = file.replace(config.basePath, "").replace(/\\/g, "/");
    const outputFilePath = join(config.basePath, relativePath.replace(".md", ".html"));

    await ensureDir(dirname(outputFilePath));

    const templatePath = join(config.views, "template.hbs");
    const renderedHtml = await renderHtml(templatePath, { content: htmlContent, frontMatter });
    await Deno.writeTextFile(outputFilePath, renderedHtml);

    console.log(`HTML file written to ${outputFilePath}`);
  }
}
