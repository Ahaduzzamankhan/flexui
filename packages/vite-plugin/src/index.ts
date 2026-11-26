/**
 * Vite plugin for FlexUI.
 *
 * The plugin never reimplements the compiler: it spawns the same precompiled
 * binary as the CLI, so a Vite project and a CLI project produce byte-identical
 * stylesheets. During `vite build` it compiles once before the bundle is
 * written; during `vite dev` it starts `flexui watch` and streams compiler
 * errors to the Vite logger.
 */

import { spawn, type ChildProcess } from "node:child_process";
import { existsSync } from "node:fs";
import { join } from "node:path";

export interface FlexUIOptions {
  /** Project directory. Defaults to Vite's `root`. */
  root?: string;
  /** Config file path, relative to the project root. */
  config?: string;
  /** Entry stylesheet with @flexui directives. */
  input?: string;
  /** Output stylesheet, relative to the project root. */
  output?: string;
  /** Globs to scan instead of the config file's `content`. */
  content?: string | readonly string[];
  /** Classes that must always be emitted. */
  safelist?: string | readonly string[];
  /** Force minification on or off instead of following Vite's build.minify. */
  minify?: boolean;
  /** Run `flexui watch` during dev. Defaults to true. */
  watch?: boolean;
  /** Compiler command. Defaults to FLEXUI_BINARY, the local .bin, or PATH. */
  command?: string;
}

/** The subset of Vite's resolved config this plugin reads. */
export interface ViteResolvedConfigLike {
  root?: string;
  command?: "build" | "serve" | string;
  build?: { minify?: boolean | string };
}

/** The subset of Vite's dev server this plugin uses. */
export interface ViteDevServerLike {
  config?: {
    logger?: {
      info?(message: string): void;
      warn?(message: string): void;
      error?(message: string): void;
    };
  };
  httpServer?: {
    once(event: "close", listener: () => void): unknown;
  } | null;
}

/** Minimal plugin context: only `error` is used, and it is optional in tests. */
export interface VitePluginContextLike {
  error?(message: string): never | void;
  warn?(message: string): void;
}

export interface FlexUIPlugin {
  name: string;
  enforce: "pre";
  configResolved?(config: ViteResolvedConfigLike): void;
  buildStart?(this: VitePluginContextLike): Promise<void>;
  configureServer?(server: ViteDevServerLike): void;
  closeBundle?(): void;
}

function binaryName(): string {
  return process.platform === "win32" ? "flexui.cmd" : "flexui";
}

/** Locates the compiler binary without failing when it is not installed yet. */
export function resolveCommand(root: string, command?: string): string {
  if (command !== undefined && command !== "") {
    return command;
  }
  const override = process.env["FLEXUI_BINARY"];
  if (override !== undefined && override !== "") {
    return override;
  }
  const local = join(root, "node_modules", ".bin", binaryName());
  return existsSync(local) ? local : "flexui";
}

function listArgs(flag: string, value: string | readonly string[] | undefined): string[] {
  if (value === undefined) {
    return [];
  }
  const text = Array.isArray(value) ? value.join(",") : String(value);
  return text === "" ? [] : [flag, text];
}

/**
 * Builds the argument vector for `flexui build` or `flexui watch`.
 * Flags precede the project directory, matching the CLI's own wrapper.
 */
export function compilerArgs(
  command: "build" | "watch",
  options: FlexUIOptions,
  root: string,
  minify: boolean,
): string[] {
  const args: string[] = [command];
  if (options.config !== undefined && options.config !== "") {
    args.push("--config", options.config);
  }
  args.push(
    ...listArgs("--input", options.input),
    ...listArgs("--output", options.output),
    ...listArgs("--content", options.content),
    ...listArgs("--safelist", options.safelist),
  );
  if (minify) {
    args.push("--minify");
  }
  args.push(root);
  return args;
}

/** The message shown when the compiler exits non-zero. */
export function formatFailure(args: readonly string[], stderr: string): string {
  const detail = stderr.trim();
  if (detail !== "") {
    return detail;
  }
  return `flexui ${args.join(" ")} exited without a diagnostic`;
}

interface RunOutcome {
  readonly code: number;
  readonly stderr: string;
}

function runOnce(command: string, args: readonly string[], cwd: string): Promise<RunOutcome> {
  return new Promise<RunOutcome>((resolve, reject) => {
    const child = spawn(command, [...args], {
      cwd,
      windowsHide: true,
      stdio: ["ignore", "ignore", "pipe"],
    });
    let stderr = "";
    child.stderr?.setEncoding("utf8");
    child.stderr?.on("data", (chunk: string) => {
      stderr += chunk;
    });
    child.on("error", reject);
    child.on("close", (code) => {
      resolve({ code: code ?? 1, stderr });
    });
  });
}

function startWatcher(
  command: string,
  args: readonly string[],
  cwd: string,
  report: (message: string) => void,
): ChildProcess {
  const child = spawn(command, [...args], {
    cwd,
    windowsHide: true,
    stdio: ["ignore", "pipe", "pipe"],
  });
  const forward = (chunk: Buffer | string): void => {
    const text = String(chunk).trimEnd();
    if (text !== "") {
      report(text);
    }
  };
  child.stdout?.on("data", forward);
  child.stderr?.on("data", forward);
  child.on("error", (error: Error) => {
    report(`could not start the FlexUI compiler: ${error.message}`);
  });
  return child;
}

/**
 * Creates the Vite plugin.
 *
 * ```ts
 * // vite.config.ts
 * import { defineConfig } from "vite";
 * import { flexui } from "@flexui/vite";
 *
 * export default defineConfig({ plugins: [flexui({ minify: true })] });
 * ```
 */
export function flexui(options: FlexUIOptions = {}): FlexUIPlugin {
  let root = options.root ?? process.cwd();
  let minify = options.minify ?? false;
  let serve = false;
  let watcher: ChildProcess | null = null;

  const stopWatcher = (): void => {
    if (watcher !== null) {
      watcher.kill();
      watcher = null;
    }
  };

  return {
    name: "flexui",
    enforce: "pre",

    configResolved(config: ViteResolvedConfigLike): void {
      if (options.root === undefined && typeof config.root === "string") {
        root = config.root;
      }
      serve = config.command === "serve";
      minify = options.minify ?? (config.build?.minify !== false);
    },

    async buildStart(this: VitePluginContextLike): Promise<void> {
      if (serve) {
        // `configureServer` runs the watcher, which performs the first build.
        return;
      }
      const command = resolveCommand(root, options.command);
      const args = compilerArgs("build", options, root, minify);
      try {
        const outcome = await runOnce(command, args, root);
        if (outcome.code !== 0) {
          const message = formatFailure(args, outcome.stderr);
          if (typeof this.error === "function") {
            this.error(message);
          }
          throw new Error(message);
        }
      } catch (error) {
        const message = error instanceof Error ? error.message : String(error);
        if (typeof this.error === "function") {
          this.error(message);
        }
        throw error;
      }
    },

    configureServer(server: ViteDevServerLike): void {
      if (options.watch === false) {
        return;
      }
      const command = resolveCommand(root, options.command);
      const args = compilerArgs("watch", options, root, minify);
      const report = (message: string): void => {
        const target = server.config?.logger;
        if (message.includes("error") || message.includes("failed")) {
          (target?.error ?? console.error)(message);
        } else {
          (target?.info ?? console.info)(message);
        }
      };
      watcher = startWatcher(command, args, root, report);
      server.httpServer?.once("close", stopWatcher);
    },

    closeBundle(): void {
      stopWatcher();
    },
  };
}

export default flexui;
