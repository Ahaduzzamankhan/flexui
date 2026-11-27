/**
 * Publishes the npm packages in dependency order.
 *
 *   bun scripts/publish.ts --dry-run   # validate tarballs without publishing
 *   bun scripts/publish.ts             # publish (needs an authenticated npm CLI)
 *
 * Order matters: the CLI depends on the platform packages through
 * optionalDependencies, and the framework helpers depend on flexui-runtime, so
 * the CLI is published last. Every package is checked before anything is
 * published, so a missing build output fails the run up front.
 */

import { spawnSync } from "node:child_process";
import { existsSync, readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";

import { readVersion } from "./version.js";

const DRY_RUN = process.argv.includes("--dry-run");
// npm requires a one-time password when publishing from an account with 2FA.
// Pass it with --otp <code>, or set NPM_OTP in the environment.
const OTP = (() => {
  const index = process.argv.indexOf("--otp");
  if (index >= 0) {
    return process.argv[index + 1];
  }
  return process.env["NPM_OTP"];
})();
// --only a,b restricts the run to those package names, which keeps an
// OTP-authenticated publish inside the code's validity window.
const ONLY = (() => {
  const index = process.argv.indexOf("--only");
  if (index < 0) {
    return null;
  }
  return new Set((process.argv[index + 1] ?? "").split(",").map((name) => name.trim()).filter(Boolean));
})();
const ROOT = process.cwd();
const PACKAGES = join(ROOT, "packages");

interface Manifest {
  readonly name?: string;
  readonly version?: string;
  readonly private?: boolean;
}

function fail(message: string): never {
  console.error(`publish: ${message}`);
  process.exit(1);
}

/** Platform packages first, then integrations, then the CLI. */
function publishOrder(): string[] {
  const platform = readdirSync(PACKAGES)
    .filter((entry) => entry.startsWith("cli-"))
    .sort()
    .map((entry) => join(PACKAGES, entry));
  const rest = ["runtime", "react", "vue", "svelte", "vite-plugin", "cli", "installer"].map((entry) =>
    join(PACKAGES, entry),
  );
  return [...platform, ...rest];
}

function verify(path: string, version: string): Manifest {
  const manifest = JSON.parse(readFileSync(join(path, "package.json"), "utf8")) as Manifest;
  if (manifest.name === undefined || manifest.name === "") {
    fail(`${path}/package.json has no name`);
  }
  if (manifest.private === true) {
    fail(`${manifest.name} is marked private`);
  }
  if (manifest.version !== version) {
    fail(`${manifest.name} is at ${manifest.version ?? "(none)"}, expected ${version}; run version:sync`);
  }
  const needed = [join(path, "dist", "index.js"), join(path, "dist", "index.d.ts")];
  const isPlatform = manifest.name?.startsWith("flexui-cli-") ?? false;
  // The installer only carries manifests and documentation.
  const isInstaller = manifest.name === "flexui-v2";
  if (!isPlatform && !isInstaller && !needed.every((file) => existsSync(file))) {
    fail(`${manifest.name} is missing build output; run build:js, build:runtime and build:integrations`);
  }
  if (manifest.name === "flexui-compiler" && !existsSync(join(path, "checksums.json"))) {
    fail("flexui-compiler is missing checksums.json; run scripts/build-binaries.sh");
  }
  return manifest;
}

/** True when this exact version already exists, so a retried batch resumes. */
function alreadyPublished(name: string, version: string): boolean {
  const result = spawnSync("npm", ["view", name, "versions", "--json"], {
    cwd: ROOT,
    encoding: "utf8",
    shell: false,
  });
  if ((result.status ?? 1) !== 0) {
    return false;
  }
  try {
    const versions = JSON.parse(result.stdout) as string | string[];
    return Array.isArray(versions) ? versions.includes(version) : versions === version;
  } catch {
    return false;
  }
}

interface PublishOutcome {
  readonly code: number;
  readonly stderr: string;
}

function publish(path: string, args: readonly string[]): PublishOutcome {
  const result = spawnSync("npm", ["publish", path, "--access", "public", "--ignore-scripts", ...args], {
    cwd: ROOT,
    stdio: ["ignore", "inherit", "pipe"],
    encoding: "utf8",
    shell: false,
  });
  if (result.error !== undefined) {
    fail(`could not run npm: ${result.error.message}`);
  }
  return { code: result.status ?? 1, stderr: result.stderr ?? "" };
}

function main(): void {
  const info = readVersion(ROOT);
  const all = publishOrder();
  const manifestsByPath = new Map(all.map((path) => [path, verify(path, info.version)]));
  const order =
    ONLY === null
      ? all
      : all.filter((path) => {
          const name = manifestsByPath.get(path)?.name ?? "";
          return ONLY.has(name) || ONLY.has(name.replace(/^@flexui\//, ""));
        });
  if (order.length === 0) {
    fail(`--only matched no packages; available: ${all.length}`);
  }
  const manifests = order.map((path) => manifestsByPath.get(path) as Manifest);

  // npm refuses to publish a prerelease without an explicit dist-tag.
  const prerelease = info.version.includes("-");
  const distTag = prerelease ? (info.channel !== "" ? info.channel : "next") : "latest";

  console.log(`${DRY_RUN ? "validating" : "publishing"} ${manifests.length} packages at ${info.version} (tag ${distTag})`);

  const auth = OTP === undefined || OTP === "" ? [] : ["--otp", OTP];
  for (const [index, path] of order.entries()) {
    const manifest = manifests[index];
    const name = manifest?.name ?? path;
    if (!DRY_RUN && alreadyPublished(name, info.version)) {
      console.log(`  skip    ${name} (${info.version} already published)`);
      continue;
    }
    console.log(`  ${DRY_RUN ? "pack" : "publish"} ${name}`);
    const outcome = publish(path, [...(DRY_RUN ? ["--dry-run"] : []), "--tag", distTag, ...auth]);
    if (outcome.code === 0) {
      continue;
    }
    // npm reports a version that already exists (including a staged one) as a
    // conflict rather than an error worth retrying, so a resumed batch skips it.
    if (/E409|previously staged|cannot publish over/i.test(outcome.stderr)) {
      console.log(`  skip    ${name} (${info.version} already published)`);
      continue;
    }
    process.stderr.write(outcome.stderr);
    fail(`npm publish failed for ${path}`);
  }

  console.log(DRY_RUN ? "dry run complete" : `published ${info.version}`);
}

try {
  main();
} catch (error) {
  fail(error instanceof Error ? error.message : String(error));
}
