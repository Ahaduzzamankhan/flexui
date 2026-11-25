/**
 * Resolves the precompiled FlexUI compiler for the current platform.
 *
 * Resolution order:
 *   1. `FLEXUI_BINARY` (explicit override used by CI and by release checks).
 *   2. The optional platform package (`@flexui/cli-<platform>-<arch>`), which
 *      npm installs as a normal dependency. No install script is involved, so
 *      it also works when lifecycle scripts are disabled.
 *   3. A vendored copy inside this package (`binaries/<platform>-<arch>/`),
 *      produced by `scripts/build-binaries.sh` for offline use.
 */

import { createHash } from "node:crypto";
import { existsSync, readFileSync } from "node:fs";
import { dirname, join } from "node:path";

export interface PlatformInfo {
  readonly pkg: string;
  readonly binary: string;
}

export const PLATFORMS: Readonly<Record<string, PlatformInfo>> = {
  "win32-x64": { pkg: "@flexui/cli-win32-x64", binary: "flexui.exe" },
  "win32-arm64": { pkg: "@flexui/cli-win32-arm64", binary: "flexui.exe" },
  "linux-x64": { pkg: "@flexui/cli-linux-x64", binary: "flexui" },
  "linux-arm64": { pkg: "@flexui/cli-linux-arm64", binary: "flexui" },
  "darwin-x64": { pkg: "@flexui/cli-darwin-x64", binary: "flexui" },
  "darwin-arm64": { pkg: "@flexui/cli-darwin-arm64", binary: "flexui" },
};

export interface ChecksumEntry {
  readonly file: string;
  readonly sha256: string;
  readonly version: string;
}

export type Checksums = Record<string, ChecksumEntry>;

export interface ResolvedBinary {
  readonly path: string;
  readonly verified: boolean;
}

/** Returns the current platform key, for example `win32-x64`. */
export function platformKey(): string {
  return `${process.platform}-${process.arch}`;
}

function checksumPath(): string {
  return join(__dirname, "..", "checksums.json");
}

/** Reads the published checksum manifest, or an empty map when absent. */
export function readChecksums(): Checksums {
  const file = checksumPath();
  if (!existsSync(file)) {
    return {};
  }
  try {
    return JSON.parse(readFileSync(file, "utf8")) as Checksums;
  } catch {
    return {};
  }
}

function verify(binaryPath: string, key: string): boolean {
  const expected = readChecksums()[key]?.sha256;
  if (expected === undefined) {
    // A locally built binary has no published checksum yet. Report it rather
    // than failing, so development builds keep working.
    return false;
  }
  const actual = createHash("sha256").update(readFileSync(binaryPath)).digest("hex");
  if (actual !== expected) {
    throw new Error(
      [
        `flexui: binary integrity check failed for ${binaryPath}`,
        `  expected sha256 ${expected}`,
        `  actual   sha256 ${actual}`,
        "Reinstall flexui to get an intact binary.",
      ].join("\n"),
    );
  }
  return true;
}

function candidatePaths(key: string): string[] {
  const info = PLATFORMS[key];
  if (info === undefined) {
    return [];
  }
  const paths: string[] = [];
  try {
    const manifest = require.resolve(`${info.pkg}/package.json`);
    paths.push(join(dirname(manifest), info.binary));
  } catch {
    // The optional platform package is not installed for this platform.
  }
  paths.push(join(__dirname, "..", "binaries", key, info.binary));
  return paths;
}

/** Locates the compiler binary or throws a message that says what to do. */
export function resolveBinary(): ResolvedBinary {
  const override = process.env["FLEXUI_BINARY"];
  if (override !== undefined && override !== "") {
    if (!existsSync(override)) {
      throw new Error(`flexui: FLEXUI_BINARY points at a missing file: ${override}`);
    }
    return { path: override, verified: false };
  }

  const key = platformKey();
  const info = PLATFORMS[key];
  if (info === undefined) {
    throw new Error(
      [
        `flexui: no precompiled compiler for ${key}.`,
        `Supported platforms: ${Object.keys(PLATFORMS).join(", ")}.`,
      ].join("\n"),
    );
  }

  const tried: string[] = [];
  for (const candidate of candidatePaths(key)) {
    tried.push(candidate);
    if (existsSync(candidate)) {
      return { path: candidate, verified: verify(candidate, key) };
    }
  }

  throw new Error(
    [
      `flexui: could not find the compiler binary for ${key}.`,
      "Looked in:",
      ...tried.map((candidate) => `  ${candidate}`),
      `Install the platform package with: npm install ${info.pkg}`,
      "or set FLEXUI_BINARY to a flexui binary.",
    ].join("\n"),
  );
}
