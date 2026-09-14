import { z } from "zod";
import * as path from "path";
import fs from "fs/promises";
import { Tool } from "./tool";
import { App } from "../app/app";

const DESCRIPTION = `Adds content to a file. If the file exists, appends content to the end. If the file does not exist, creates a new file with the given content.

WHEN TO USE THIS TOOL:
- Use when you need to append text to an existing file
- Use when you need to create a new file with content
- Use when you want to add lines or sections to a file without overwriting existing content

PARAMETERS:
- file_path: The relative path to the file (must be relative, not absolute)
- content: The content to append to the file (required)

IMPORTANT:
- This tool appends content; it does NOT replace existing content
- To overwrite a file, use the Edit tool instead
- Always use relative file paths`;

export const add = Tool.define({
  name: "kotha.add",
  description: DESCRIPTION,
  parameters: z.object({
    filePath: z.string().describe("The relative path to the file to add content to"),
    content: z.string().describe("The content to append to the file"),
  }),
  async execute(params) {
    if (!params.filePath) {
      throw new Error("filePath is required");
    }
    if (!params.content) {
      throw new Error("content is required");
    }

    const app = await App.use();
    let filePath = params.filePath;
    if (!path.isAbsolute(filePath)) {
      filePath = path.join(app.root, filePath);
    }

    const stats = await fs.stat(filePath).catch(() => null);
    if (stats?.isDirectory()) {
      throw new Error(`Path is a directory, not a file: ${filePath}`);
    }

    const exists = stats !== null;

    if (exists) {
      await Bun.write(filePath, await Bun.file(filePath).text() + params.content);
    } else {
      const dir = path.dirname(filePath);
      await fs.mkdir(dir, { recursive: true });
      await Bun.write(filePath, params.content);
    }

    return {
      metadata: {
        created: !exists,
        appended: exists,
      },
      output: exists
        ? `Added content to ${filePath}`
        : `Created file ${filePath} with content`,
    };
  },
});
