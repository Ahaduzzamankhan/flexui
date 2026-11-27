# flexui-v2

One command installs everything FlexUI needs. The package pulls in the compiler
and the native binary that matches your operating system and CPU, so you never
pick a platform package by hand.

```bash
npm install flexui-v2
```

Or globally, when you want the `flexui` command available everywhere:

```bash
npm install -g flexui-v2
```

After that:

```bash
npx flexui init
npx flexui build
flexui version        # global install
```

## What it installs

| Piece | Package |
| --- | --- |
| The `flexui` command and its TypeScript API | `flexui-compiler` |
| The compiler binary for your machine | `flexui-cli-windows-x64`, `flexui-cli-windows-arm64`, `flexui-cli-linux-x64`, `flexui-cli-linux-arm64`, `flexui-cli-darwin-x64` or `flexui-cli-darwin-arm64` |

npm selects exactly one platform package: each declares its `os` and `cpu`, so
Windows, Linux and macOS — on x64 or ARM64 — all resolve to the right executable
with no scripts and no downloads during install.

The wrapper verifies the binary's SHA-256 against the manifest shipped with
`flexui-compiler` before the first build, so a truncated or tampered download is
refused instead of executed.

## Requirements

Node.js 18 or newer, or Bun. The compiler itself is a native binary and does not
need Go, a C toolchain or any system library on the machine using it.

## Using it in a project

```bash
npm install --save-dev flexui-v2
npx flexui init
npx flexui build
```

The installed command is `flexui`. For Vite projects, add `flexui-vite`:

```bash
npm install --save-dev flexui-vite
```

```ts
// vite.config.ts
import { flexui } from "flexui-vite";
export default { plugins: [flexui()] };
```

## Platforms

Windows x64, Windows ARM64, Linux x64, Linux ARM64, macOS x64 and macOS ARM64 —
one package covers all of them.

## Documentation

- [Getting started](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/getting-started.md)
- [CLI reference](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/cli.md)
- [Repository](https://github.com/Ahaduzzamankhan/flexui#readme)

## License

MIT
