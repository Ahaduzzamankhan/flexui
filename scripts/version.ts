/**
 * version.json is the single source of truth for the release version.
 *
 *   bun scripts/version.ts show      # resolved version info as JSON
 *   bun scripts/version.ts check     # exit 1 if any package or changelog drifts
 *   bun scripts/version.ts sync      # rewrite package versions from version.json
 *   bun scripts/version.ts tag       # print the git tag (v2.0.0-alpha.2)
 *   bun scripts/version.ts notes     # print the CHANGELOG section for this version
 *   bun scripts/version.ts notes --out dist-release/RELEASE_NOTES.md
 *
 * Every command reads version.json, so there is exactly one value to bump.
 * Run from the repository root (npm and bun scripts already do).
 */

import { mkdirSync, readFileSync, readdirSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";

export interface VersionInfo {
  readonly name: string;
  readonly version: string;
  readonly channel: string;
}

export interface ManifestIssue {
  readonly file: string;
  readonly found: string;
  readonly expected: string;
}

export const VERSION_FILE = "version.json";
export const CHANGELOG_FILE = "CHANGELOG.md";

function readJSON(path: string): unknown {
  return JSON.parse(readFileSync(path, "utf8")) as unknown;
}

/** Reads version.json and validates its fields. */
export function readVersion(root: string = process.cwd()): VersionInfo {
  const data = readJSON(join(root, VERSION_FILE)) as Partial<VersionInfo>;
  if (typeof data.name !== "string" || data.name === "") {
    throw new Error(`${VERSION_FILE}: missing "name"`);
  }
  if (typeof data.version !== "string" || !/^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?(\+[0-9A-Za-z.]+)?$/.test(data.version)) {
    throw new Error(`${VERSION_FILE}: "version" must be a semantic version, got ${JSON.stringify(data.version)}`);
  }
  if (data.channel !== undefined && typeof data.channel !== "string") {
    throw new Error(`${VERSION_FILE}: "channel" must be a string`);
  }
  return { name: data.name, version: data.version, channel: data.channel ?? "" };
}

/** Returns a git tag name such as v2.0.0-alpha.2. */
export function tagName(info: VersionInfo): string {
  return `v${info.version}`;
}

/** Every package.json that must carry the release version. */
export function manifestPaths(root: string = process.cwd()): string[] {
  const paths = [join(root, "package.json")];
  const packages = join(root, "packages");
  let entries: string[] = [];
  try {
    entries = readdirSync(packages);
  } catch {
    return paths;
  }
  for (const entry of entries.sort()) {
    const manifest = join(packages, entry, "package.json");
    try {
      readJSON(manifest);
      paths.push(manifest);
    } catch {
      // Not a package directory.
    }
  }
  return paths;
}

interface Manifest {
  version?: string;
  dependencies?: Record<string, string>;
  devDependencies?: Record<string, string>;
  optionalDependencies?: Record<string, string>;
  peerDependencies?: Record<string, string>;
}

const DEP_FIELDS = ["dependencies", "devDependencies", "optionalDependencies", "peerDependencies"] as const;

/** Returns every manifest or dependency whose version differs from version.json. */
export function checkVersions(root: string = process.cwd()): ManifestIssue[] {
  const info = readVersion(root);
  const issues: ManifestIssue[] = [];

  for (const path of manifestPaths(root)) {
    const manifest = readJSON(path) as Manifest;
    if (manifest.version !== info.version) {
      issues.push({ file: path, found: manifest.version ?? "(none)", expected: info.version });
    }
    for (const field of DEP_FIELDS) {
      for (const [name, range] of Object.entries(manifest[field] ?? {})) {
        if (range === info.version) {
          continue;
        }
        // Workspace-internal FlexUI packages are pinned to the release version.
        if (name.startsWith("flexui-") || name.startsWith("@flexui/")) {
          if (/^[0-9]/.test(range)) {
            issues.push({ file: `${path}#${field}.${name}`, found: range, expected: info.version });
          }
        }
      }
    }
  }
  return issues;
}

/**
 * Rewrites every package manifest to the version.json version.
 * Returns the files that changed.
 */
export function syncVersions(root: string = process.cwd()): string[] {
  const info = readVersion(root);
  const changed: string[] = [];

  for (const path of manifestPaths(root)) {
    const before = readFileSync(path, "utf8");
    const manifest = JSON.parse(before) as Manifest;
    const previous = manifest.version;
    if (previous === undefined) {
      continue;
    }

    let after = before.replace(/"version"\s*:\s*"[^"]*"/, `"version": "${info.version}"`);
    if (previous !== info.version) {
      // Internal dependencies pin the same version, so rewrite those too.
      after = after.split(`"${previous}"`).join(`"${info.version}"`);
    }
    if (after !== before) {
      writeFileSync(path, after);
      changed.push(path);
    }
  }
  return changed;
}

/**
 * Extracts the CHANGELOG section for a version.
 * Matches "## [1.2.3]", "## 1.2.3" and "## [1.2.3] - 2026-01-01" headings.
 */
export function changelogSection(markdown: string, version: string): string | null {
  const lines = markdown.split(/\r?\n/);
  const escaped = version.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const heading = new RegExp(`^##\\s+\\[?v?${escaped}\\]?(\\s|$)`);
  const start = lines.findIndex((line) => heading.test(line));
  if (start < 0) {
    return null;
  }
  const body: string[] = [];
  for (let index = start + 1; index < lines.length; index += 1) {
    const line = lines[index] ?? "";
    if (/^##\s+/.test(line)) {
      break;
    }
    body.push(line);
  }
  return body.join("\n").trim();
}

/** Reads CHANGELOG.md and returns this version's notes, or null. */
export function notesFor(root: string, info: VersionInfo): string | null {
  try {
    return changelogSection(readFileSync(join(root, CHANGELOG_FILE), "utf8"), info.version);
  } catch {
    return null;
  }
}

function flagValue(args: readonly string[], name: string): string | undefined {
  const index = args.indexOf(name);
  if (index < 0) {
    return undefined;
  }
  return args[index + 1];
}

function main(args: readonly string[]): number {
  const command = args[0] ?? "show";
  const info = readVersion();

  switch (command) {
    case "show": {
      console.log(JSON.stringify({ ...info, tag: tagName(info) }, null, 2));
      return 0;
    }
    case "tag": {
      console.log(tagName(info));
      return 0;
    }
    case "check": {
      const issues = checkVersions();
      const notes = notesFor(process.cwd(), info);
      if (issues.length === 0 && notes !== null) {
        console.log(`versions consistent at ${info.version}`);
        return 0;
      }
      for (const issue of issues) {
        console.error(`version drift: ${issue.file} has ${issue.found}, expected ${issue.expected}`);
      }
      if (notes === null) {
        console.error(`CHANGELOG.md has no "## ${info.version}" section`);
      }
      console.error("run `bun run version:sync` after editing version.json");
      return 1;
    }
    case "sync": {
      const changed = syncVersions();
      if (changed.length === 0) {
        console.log(`all packages already at ${info.version}`);
      } else {
        for (const file of changed) {
          console.log(`  updated ${file}`);
        }
      }
      return 0;
    }
    case "notes": {
      const notes = notesFor(process.cwd(), info);
      if (notes === null) {
        console.error(`CHANGELOG.md has no "## ${info.version}" section`);
        return 1;
      }
      const out = flagValue(args, "--out");
      if (out === undefined) {
        console.log(notes);
      } else {
        mkdirSync(dirname(out), { recursive: true });
        writeFileSync(out, `${notes}\n`);
        console.log(`wrote ${out}`);
      }
      return 0;
    }
    default:
      console.error(`unknown command ${JSON.stringify(command)} (show, check, sync, tag, notes)`);
      return 2;
  }
}

/**
 * True when this file is the entry point. `import.meta.main` is Bun-only and
 * the scripts compile as CommonJS for `tsc`, so compare argv instead.
 */
function invokedDirectly(): boolean {
  const entry = process.argv[1];
  return entry !== undefined && /(^|[\\/])version\.(ts|js|mjs|cjs)$/.test(entry);
}

if (invokedDirectly()) {
  try {
    process.exitCode = main(process.argv.slice(2));
  } catch (error) {
    console.error(error instanceof Error ? error.message : String(error));
    process.exitCode = 1;
  }
}
