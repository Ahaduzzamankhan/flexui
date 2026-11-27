# flexui-compiler

The FlexUI 2.0 compiler. It scans your markup and scripts, keeps only the
utilities you actually use, and writes one deterministic stylesheet.

The compiler is a native binary built in Go. This package finds the right
prebuilt binary for your machine, checks its SHA-256 against the bundled
manifest, and runs it. There is nothing to compile and no install script to
trust.

```bash
npm install --save-dev flexui-compiler
npx flexui init
npx flexui build --minify
```

The installed command is `flexui`, so `npx flexui` and
`node_modules/.bin/flexui` both work.

## Commands

| Command | What it does |
| --- | --- |
| `flexui init` | Write `flexui.config.json` and an entry stylesheet. |
| `flexui build` | Compile the project once and write the stylesheet. |
| `flexui watch` | Recompile when source files change. |
| `flexui inspect` | Print the rules a build would emit, without writing a file. |
| `flexui clean` | Remove the build cache and generated output. |
| `flexui version` | Print the compiler version. |

Flags: `--minify`, `--watch`, `--config <path>`, `--input <path>`,
`--output <path>`, `--content <glob>`, `--safelist <classes>`,
`--no-cache`, `--verbose`.

## Configuration

`flexui init` creates `flexui.config.json`:

```json
{
  "content": ["src/**/*.{html,jsx,tsx,vue,svelte,astro,md}"],
  "input": "src/styles.css",
  "output": "dist/flexui.css",
  "theme": { "colors": { "brand": "#4f46e5" } }
}
```

Themes, components, animations, variants and the content globs all live in this
file; the [project documentation](https://github.com/Ahaduzzamankhan/flexui#readme)
walks through each block.

## TypeScript API

```ts
import { build, inspect, version } from "flexui-compiler";

await build({ minify: true, output: "dist/flexui.css" });
const result = await inspect();
console.log(result.stdout);
console.log(await version());
```

| Export | Purpose |
| --- | --- |
| `build(options)` | Compile the project and write the stylesheet. |
| `watch(options)` | Start watch mode. |
| `inspect(options)` | Return the resolved rules. |
| `version(options)` | Read the compiler version. |
| `run(args, options)` | Run the compiler with arbitrary arguments, capturing output. |
| `runInherit(args, options)` | Run with inherited stdio and return the exit code. |
| `resolveBinary()` | Path and verification state of the resolved binary. |
| `optArgs(options)` | Map options to compiler flags. |

`RunOptions` covers `cwd`, `root`, `config`, `input`, `output`, `content`,
`safelist`, `minify`, `verbose`, `noCache` and `env`. A failed build rejects with
a `FlexUIError` that carries `exitCode`, `stdout` and `stderr`.

## Platform packages

The binary arrives through an optional dependency that matches your operating
system and CPU:

| Platform | Package |
| --- | --- |
| Windows x64 | `flexui-cli-windows-x64` |
| Windows ARM64 | `flexui-cli-windows-arm64` |
| Linux x64 | `flexui-cli-linux-x64` |
| Linux ARM64 | `flexui-cli-linux-arm64` |
| macOS x64 | `flexui-cli-darwin-x64` |
| macOS ARM64 | `flexui-cli-darwin-arm64` |

`FLEXUI_BINARY` points the wrapper at any `flexui` executable when you keep the
compiler somewhere else, for example a system-wide install or a CI cache.

## Requirements

Node.js 18 or newer, or Bun. The package has no runtime dependencies.

## License

MIT
