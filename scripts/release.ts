/**
 * Release automation: version check, binaries, changelog notes, git tag.
 *
 *   bun scripts/release.ts                 # build + prepare dist-release, no tag
 *   bun scripts/release.ts --publish       # also tag, push and create the GitHub release
 *   bun scripts/release.ts --skip-build    # reuse existing packages/cli-<platform>/ binaries
 *
 * The version comes from version.json (see scripts/version.ts). Release notes
 * are the CHANGELOG.md section for that version, so a release cannot be
 * published without notes.
 */

import { spawnSync } from "node:child_process";
import { mkdirSync, readFileSync, readdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";

import { checkVersions, notesFor, readVersion, tagName, type VersionInfo } from "./version.js";

const PUBLISH = process.argv.includes("--publish");
const SKIP_BUILD = process.argv.includes("--skip-build");
const ROOT = process.cwd();
const OUT = join(ROOT, "dist-release");

interface Checksums {
  readonly [platform: string]: { readonly file: string; readonly sha256: string; readonly version: string };
}

function fail(message: string): never {
  console.error(`release: ${message}`);
  process.exit(1);
}

function run(command: string, args: readonly string[], quiet = false): void {
  const result = spawnSync(command, [...args], {
    cwd: ROOT,
    stdio: quiet ? "pipe" : "inherit",
    encoding: "utf8",
    shell: false,
  });
  if (result.error !== undefined) {
    fail(`could not run ${command}: ${result.error.message}`);
  }
  if ((result.status ?? 1) !== 0) {
    const detail = (result.stderr ?? "").trim();
    fail(`${command} ${args.join(" ")} exited with ${result.status ?? "no status"}${detail === "" ? "" : `\n${detail}`}`);
  }
}

function capture(command: string, args: readonly string[]): { status: number; stdout: string } {
  const result = spawnSync(command, [...args], { cwd: ROOT, encoding: "utf8", shell: false });
  return { status: result.status ?? 1, stdout: (result.stdout ?? "").trim() };
}

/** Every cross-compiled binary produced by scripts/build-binaries.sh. */
function binaryAssets(): string[] {
  const assets: string[] = [];
  for (const entry of readdirSync(ROOT)) {
    if (!entry.startsWith("cli-") || !entry.includes("-")) {
      continue;
    }
    const dir = join(ROOT, "packages", entry);
    for (const name of readdirSync(dir)) {
      if (name === "flexui" || name === "flexui.exe") {
        assets.push(join("packages", entry, name));
      }
    }
  }
  if (assets.length === 0) {
    fail("no binaries found; run `bash scripts/build-binaries.sh` first");
  }
  return assets.sort();
}

function writeChecksumFile(): string {
  const path = join(ROOT, "packages", "cli", "checksums.json");
  const sums = JSON.parse(readFileSync(path, "utf8")) as Checksums;
  const lines = Object.entries(sums)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([platform, entry]) => `${entry.sha256}  packages/cli-${platform}/${entry.file}`);
  const out = join(OUT, "SHA256SUMS");
  writeFileSync(out, `${lines.join("\n")}\n`);
  return out;
}

function writeNotes(info: VersionInfo, notes: string, assets: readonly string[]): string {
  const body = [
    `# FlexUI ${info.version}`,
    "",
    notes,
    "",
    "## Assets",
    "",
    "- Precompiled compilers for Windows, Linux and macOS on x64 and ARM64",
    "- `checksums.json` and `SHA256SUMS` with the SHA-256 of every binary",
    "",
    "```",
    ...assets.map((asset) => asset.replace(/\\\\/g, "/")),
    "```",
    "",
  ].join("\n");
  const path = join(OUT, "RELEASE_NOTES.md");
  writeFileSync(path, body);
  return path;
}

function main(): void {
  const info = readVersion(ROOT);
  const tag = tagName(info);

  const issues = checkVersions(ROOT);
  if (issues.length > 0) {
    for (const issue of issues) {
      console.error(`  ${issue.file}: ${issue.found} (expected ${issue.expected})`);
    }
    fail("package versions do not match version.json; run `bun run version:sync`");
  }

  const notes = notesFor(ROOT, info);
  if (notes === null) {
    fail(`CHANGELOG.md has no "## ${info.version}" section; add release notes first`);
  }

  mkdirSync(OUT, { recursive: true });

  if (!SKIP_BUILD) {
    console.log(`building binaries for ${info.version}`);
    run("bash", ["scripts/build-binaries.sh"]);
  }

  const assets = binaryAssets();
  const sums = writeChecksumFile();
  const notesPath = writeNotes(info, notes, assets);

  console.log(`\nprepared ${tag}`);
  console.log(`  notes    ${notesPath}`);
  console.log(`  sums     ${sums}`);
  console.log(`  assets   ${assets.length} binaries`);

  if (!PUBLISH) {
    console.log("\nnothing published. Re-run with --publish to tag and create the release:");
    console.log(`  git tag -a ${tag} -m "FlexUI ${info.version}"`);
    console.log(`  git push origin ${tag}`);
    console.log(`  gh release create ${tag} --notes-file dist-release/RELEASE_NOTES.md`);
    return;
  }

  const status = capture("git", ["status", "--porcelain"]);
  if (status.status !== 0) {
    fail("not a git repository");
  }
  if (status.stdout !== "") {
    fail("working tree is dirty; commit the version bump and changelog first");
  }

  if (capture("git", ["rev-parse", "--verify", `refs/tags/${tag}`]).status === 0) {
    fail(`tag ${tag} already exists`);
  }

  run("git", ["tag", "-a", tag, "-m", `FlexUI ${info.version}`]);
  run("git", ["push", "origin", tag]);

  if (capture("gh", ["--version"]).status !== 0) {
    fail(`tag ${tag} was pushed but the GitHub CLI (gh) is not installed, so no release was created`);
  }

  const releaseAssets = [...assets, "packages/cli/checksums.json", sums];
  run("gh", ["release", "create", tag, "--title", `FlexUI ${info.version}`, "--notes-file", notesPath, ...releaseAssets]);

  console.log(`\nreleased ${tag}`);
}

try {
  main();
} catch (error) {
  fail(error instanceof Error ? error.message : String(error));
}
