import { describe, expect, test } from "bun:test";
import { App } from "../../src/app/app";
import { add } from "../../src/tool/add";
import * as path from "path";
import * as fs from "fs/promises";

describe("tool.add", () => {
  test("creates new file", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const tmpDir = await fs.mkdtemp(path.join(process.cwd(), "add-test-"));
      const testFile = path.join(tmpDir, "new-file.txt");

      let result = await add.execute(
        {
          filePath: testFile,
          content: "Hello, World!\n",
        },
        {
          toolCallId: "test",
          messages: [],
        },
      );

      expect(result.metadata.created).toBe(true);
      expect(result.metadata.appended).toBe(false);
      expect(result.output).toContain("Created file");

      const content = await Bun.file(testFile).text();
      expect(content).toBe("Hello, World!\n");

      await fs.rm(tmpDir, { recursive: true });
    });
  });

  test("appends to existing file", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const tmpDir = await fs.mkdtemp(path.join(process.cwd(), "add-test-"));
      const testFile = path.join(tmpDir, "existing.txt");

      await Bun.write(testFile, "Original content.\n");

      let result = await add.execute(
        {
          filePath: testFile,
          content: "Appended content.\n",
        },
        {
          toolCallId: "test",
          messages: [],
        },
      );

      expect(result.metadata.created).toBe(false);
      expect(result.metadata.appended).toBe(true);
      expect(result.output).toContain("Added content");

      const content = await Bun.file(testFile).text();
      expect(content).toBe("Original content.\nAppended content.\n");

      await fs.rm(tmpDir, { recursive: true });
    });
  });

  test("appends multiple times", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const tmpDir = await fs.mkdtemp(path.join(process.cwd(), "add-test-"));
      const testFile = path.join(tmpDir, "multi.txt");

      await Bun.write(testFile, "Line 1\n");

      await add.execute(
        {
          filePath: testFile,
          content: "Line 2\n",
        },
        {
          toolCallId: "test-1",
          messages: [],
        },
      );

      await add.execute(
        {
          filePath: testFile,
          content: "Line 3\n",
        },
        {
          toolCallId: "test-2",
          messages: [],
        },
      );

      const content = await Bun.file(testFile).text();
      expect(content).toBe("Line 1\nLine 2\nLine 3\n");

      await fs.rm(tmpDir, { recursive: true });
    });
  });

  test("creates parent directories", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const tmpDir = await fs.mkdtemp(path.join(process.cwd(), "add-test-"));
      const testFile = path.join(tmpDir, "nested", "deep", "file.txt");

      let result = await add.execute(
        {
          filePath: testFile,
          content: "Deep content.\n",
        },
        {
          toolCallId: "test",
          messages: [],
        },
      );

      expect(result.metadata.created).toBe(true);
      expect(await Bun.file(testFile).exists()).toBe(true);

      await fs.rm(tmpDir, { recursive: true });
    });
  });

  test("rejects when filePath is missing", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const result = await add.execute(
        {
          content: "test",
        } as any,
        {
          toolCallId: "test",
          messages: [],
        },
      );
      expect(result.metadata.error).toBe(true);
      expect(result.output).toContain("filePath is required");
    });
  });

  test("rejects when content is missing", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const result = await add.execute(
        {
          filePath: "test.txt",
        } as any,
        {
          toolCallId: "test",
          messages: [],
        },
      );
      expect(result.metadata.error).toBe(true);
      expect(result.output).toContain("content is required");
    });
  });

  test("rejects for directory path", async () => {
    await App.provide({ directory: process.cwd() }, async () => {
      const tmpDir = await fs.mkdtemp(path.join(process.cwd(), "add-test-"));

      const result = await add.execute(
        {
          filePath: tmpDir,
          content: "test",
        },
        {
          toolCallId: "test",
          messages: [],
        },
      );
      expect(result.metadata.error).toBe(true);
      expect(result.output).toContain("Path is a directory");

      await fs.rm(tmpDir, { recursive: true });
    });
  });
});
