#!/usr/bin/env node
/** Bin entry point for the `flexui` command (package `flexui-compiler`). */

import { runInherit } from "./index";

try {
  process.exitCode = runInherit(process.argv.slice(2));
} catch (error) {
  console.error(error instanceof Error ? error.message : String(error));
  process.exitCode = 1;
}
