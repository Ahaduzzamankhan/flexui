/**
 * Programmatic entry point for the FlexUI npm package.
 *
 * Every function here runs the native compiler; nothing duplicates compiler
 * logic in TypeScript, so the CLI, the npm wrapper and any future bundler
 * integration share one implementation.
 */

import { spawn, spawnSync } from "node:child_process";

import { resolveBinary } from "./binary";

export * from "./binary";

export interface RunOptions {
  /** Working directory for the compiler process. */
  cwd?: string;
  /** Project directory passed as a positional argument. */
  root?: string;
  config?: string;
  input?: string;
  output?: string;
  content?: string | readonly string[];
  safelist?: string | readonly string[];
  minify?: boolean;
  verbose?: boolean;
  noCache?: boolean;
  env?: NodeJS.ProcessEnv;
}

export interface RunResult {
  readonly code: number;
  readonly stdout: string;
  readonly stderr: string;
}

/** Raised when the compiler exits with a non-zero status. */
export class FlexUIError extends Error {
  readonly exitCode: number;
  readonly stdout: string;
  readonly stderr: string;

  constructor(message: string, exitCode: number, stdout: string, stderr: string) {
    super(message);
    this.name = "FlexUIError";
    this.exitCode = exitCode;
    this.stdout = stdout;
    this.stderr = stderr;
  }
}

function asList(value: string | readonly string[]): string {
  return Array.isArray(value) ? value.join(",") : String(value);
}

/** Translates options into compiler flags. */
export function optArgs(options: RunOptions = {}): string[] {
  const args: string[] = [];
  if (options.config !== undefined) args.push("--config", options.config);
  if (options.input !== undefined) args.push("--input", options.input);
  if (options.output !== undefined) args.push("--output", options.output);
  if (options.content !== undefined) args.push("--content", asList(options.content));
  if (options.safelist !== undefined) args.push("--safelist", asList(options.safelist));
  if (options.minify === true) args.push("--minify");
  if (options.verbose === true) args.push("--verbose");
  if (options.noCache === true) args.push("--no-cache");
  return args;
}

/**
 * Runs the compiler with inherited stdio, so exit codes and signals behave
 * exactly like a native binary. Used by the bin entry point.
 */
export function runInherit(args: readonly string[], options: RunOptions = {}): number {
  const binary = resolveBinary();
  const result = spawnSync(binary.path, [...args], {
    stdio: "inherit",
    cwd: options.cwd ?? process.cwd(),
    env: options.env ?? process.env,
    windowsHide: true,
  });
  if (result.error !== undefined) {
    throw result.error;
  }
  return result.status ?? 1;
}

/** Runs the compiler and captures its output. */
export function run(args: readonly string[], options: RunOptions = {}): Promise<RunResult> {
  const binary = resolveBinary();
  return new Promise<RunResult>((resolve, reject) => {
    const child = spawn(binary.path, [...args], {
      cwd: options.cwd ?? process.cwd(),
      env: options.env ?? process.env,
      windowsHide: true,
    });

    let stdout = "";
    let stderr = "";
    child.stdout?.setEncoding("utf8");
    child.stderr?.setEncoding("utf8");
    child.stdout?.on("data", (chunk: string) => {
      stdout += chunk;
    });
    child.stderr?.on("data", (chunk: string) => {
      stderr += chunk;
    });
    child.on("error", reject);
    child.on("close", (code) => {
      const exitCode = code ?? 1;
      if (exitCode === 0) {
        resolve({ code: exitCode, stdout, stderr });
        return;
      }
      reject(new FlexUIError(stderr.trim() || `flexui exited with code ${exitCode}`, exitCode, stdout, stderr));
    });
  });
}

function command(name: string, options: RunOptions): string[] {
  const args = [name, ...optArgs(options)];
  if (options.root !== undefined) {
    args.push(options.root);
  }
  return args;
}

/** Compiles the project and writes the stylesheet. */
export function build(options: RunOptions = {}): Promise<RunResult> {
  return run(command("build", options), options);
}

/** Starts watch mode; resolves when the process exits. */
export function watch(options: RunOptions = {}): Promise<RunResult> {
  return run(command("watch", options), options);
}

/** Prints resolved rules without writing a file. */
export function inspect(options: RunOptions = {}): Promise<RunResult> {
  return run(command("inspect", options), options);
}

/** Reads the compiler version. */
export function version(options: RunOptions = {}): Promise<RunResult> {
  return run(["version"], options);
}
