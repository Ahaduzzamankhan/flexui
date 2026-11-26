import { describe, expect, test } from "bun:test";

import { compilerArgs, flexui, formatFailure, resolveCommand } from "../src/index.js";

describe("compilerArgs", () => {
  test("passes flags before the project directory", () => {
    const args = compilerArgs(
      "build",
      {
        config: "config/flexui.json",
        input: "src/styles.css",
        output: "dist/app.css",
        content: ["./index.html", "./src/**/*.tsx"],
        safelist: ["md:flex", "dark:bg-slate-900"],
      },
      "/project",
      true,
    );

    expect(args).toEqual([
      "build",
      "--config",
      "config/flexui.json",
      "--input",
      "src/styles.css",
      "--output",
      "dist/app.css",
      "--content",
      "./index.html,./src/**/*.tsx",
      "--safelist",
      "md:flex,dark:bg-slate-900",
      "--minify",
      "/project",
    ]);
  });

  test("omits unset options and empty lists", () => {
    expect(compilerArgs("watch", {}, "/project", false)).toEqual(["watch", "/project"]);
    expect(compilerArgs("build", { content: "", safelist: [] }, "/p", false)).toEqual([
      "build",
      "/p",
    ]);
  });

  test("accepts a single content glob", () => {
    const args = compilerArgs("build", { content: "src/**/*.vue" }, "/p", false);
    expect(args).toContain("src/**/*.vue");
    expect(args.indexOf("--content")).toBeLessThan(args.indexOf("/p"));
  });
});

describe("formatFailure", () => {
  test("uses the compiler diagnostic when present", () => {
    expect(formatFailure(["build", "/p"], "  read input CSS: missing\n")).toBe(
      "read input CSS: missing",
    );
  });

  test("falls back to the command line", () => {
    expect(formatFailure(["build", "/p"], "   ")).toBe("flexui build /p exited without a diagnostic");
  });
});

describe("resolveCommand", () => {
  test("prefers an explicit command", () => {
    expect(resolveCommand("/project", "./bin/flexui")).toBe("./bin/flexui");
  });

  test("falls back to PATH when nothing local exists", () => {
    const previous = process.env["FLEXUI_BINARY"];
    delete process.env["FLEXUI_BINARY"];
    try {
      expect(resolveCommand("/definitely/missing/project")).toBe("flexui");
    } finally {
      if (previous !== undefined) {
        process.env["FLEXUI_BINARY"] = previous;
      }
    }
  });

  test("honours FLEXUI_BINARY", () => {
    const previous = process.env["FLEXUI_BINARY"];
    process.env["FLEXUI_BINARY"] = "/tmp/flexui-test-binary";
    try {
      expect(resolveCommand("/project")).toBe("/tmp/flexui-test-binary");
    } finally {
      if (previous === undefined) {
        delete process.env["FLEXUI_BINARY"];
      } else {
        process.env["FLEXUI_BINARY"] = previous;
      }
    }
  });
});

describe("flexui()", () => {
  test("registers as a pre plugin with the expected hooks", () => {
    const plugin = flexui();
    expect(plugin.name).toBe("flexui");
    expect(plugin.enforce).toBe("pre");
    expect(typeof plugin.configResolved).toBe("function");
    expect(typeof plugin.buildStart).toBe("function");
    expect(typeof plugin.configureServer).toBe("function");
    expect(typeof plugin.closeBundle).toBe("function");
  });

  test("tracks serve mode so buildStart does not double compile", async () => {
    const plugin = flexui({ root: "/project" });
    plugin.configResolved?.({ root: "/project", command: "serve" });
    // In serve mode buildStart returns immediately without spawning anything,
    // even though no binary exists in this test environment.
    await expect(plugin.buildStart?.call({})).resolves.toBeUndefined();
  });
});
