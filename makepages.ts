import { walk, ensureDir } from "@std/fs";
import { join, dirname } from "@std/path";
import { parse as parseYaml } from "@std/yaml";
import { readConfig } from "./config.ts";
import MarkdownIt from "https://esm.sh/markdown-it@13.0.1";
import Handlebars from "https://esm.sh/handlebars@4.7.7";

export async function findMarkdownFiles(dirPath: string): Promise<string[]> {
  const files: string[] = [];
  for await (const entry of walk(dirPath)) {
    if (entry.path.endsWith(".md")) {
      files.push(entry.path);
    }
  }
  return files;
}

export async function parseFrontMatter(filePath: string): Promise<{ frontMatter: any, content: string }> {
  const fileContent = await Deno.readTextFile(filePath);
  const [frontMatterSection, ...markdownContent] = fileContent.split("---\n");

  const frontMatter = parseYaml(frontMatterSection.replace("---", "").trim());
  const content = markdownContent.join("---\n").trim();

  return { frontMatter, content };
}

export function convertMarkdownToHtml(markdown: string): string {
  const md = new MarkdownIt();
  return md.render(markdown);
}

export function renderHtml(templatePath: string, data: { content: string, frontMatter: any }): string {
  const templateContent = Deno.readTextFileSync(templatePath);
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
    const renderedHtml = renderHtml(templatePath, { content: htmlContent, frontMatter });
    await Deno.writeTextFile(outputFilePath, renderedHtml);

    console.log(`HTML file written to ${outputFilePath}`);
  }
}
