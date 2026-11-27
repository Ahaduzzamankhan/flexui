# flexui-vite

The Vite plugin for FlexUI. It runs the same compiler binary as the CLI, so a
Vite project and a plain CLI project produce byte-identical stylesheets. During
`vite build` the stylesheet is compiled before the bundle is written; during
`vite dev` the plugin starts `flexui watch` and forwards compiler errors to the
Vite logger with file and line information.

```bash
npm install --save-dev flexui-compiler flexui-vite
```

## Usage

```ts
// vite.config.ts
import { defineConfig } from "vite";
import { flexui } from "flexui-vite";

export default defineConfig({
  plugins: [flexui({ minify: true })],
});
```

The plugin resolves the compiler through `flexui-compiler`, so the same platform
packages and checksum verification apply.

## Options

| Option | Type | Default | Purpose |
| --- | --- | --- | --- |
| `minify` | `boolean` | `false` | Minify the generated stylesheet. |
| `input` | `string` | config value | Entry stylesheet. |
| `output` | `string` | config value | Output stylesheet path. |
| `content` | `string[]` | config value | Override the content globs. |
| `safelist` | `string[]` | config value | Classes that are always emitted. |
| `config` | `string` | `flexui.config.json` | Path to the configuration file. |
| `command` | `string` | resolved | Use a specific compiler binary. |
| `watch` | `boolean` | `true` in dev | Recompile while the dev server runs. |

## How it works

- **Dev** — the plugin spawns `flexui watch` once, keeps it running with the dev
  server, and sends `error` events so Vite overlays show compiler failures.
- **Build** — the plugin compiles once in `buildStart`, before Vite reads the
  stylesheet, so the bundle always contains the latest output.
- **Errors** — stderr is captured and formatted with the failing arguments, so a
  broken class or config points at a line instead of a silent empty file.

## Programmatic helpers

`resolveCommand`, `compilerArgs` and `formatFailure` are exported for build tools
that want the same behaviour without the plugin itself.

## Documentation

See [docs/integrations.md](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/integrations.md)
and the [FlexUI repository](https://github.com/Ahaduzzamankhan/flexui#readme).

## License

MIT
