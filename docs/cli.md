# CLI reference

```
flexui <command> [options] [project-dir]
```

`project-dir` defaults to the current directory.

## Commands

| Command | Purpose |
| --- | --- |
| `flexui build` | Compile the utilities used by your sources into one stylesheet. |
| `flexui watch` | Rebuild whenever a scanned file changes. Same options as `build`. |
| `flexui inspect` | Print resolved rules in emission order. Writes nothing. |
| `flexui init` | Create `flexui.config.json` and `src/styles.css`. Never overwrites. |
| `flexui clean` | Remove the build cache directory (`.flexui`). |
| `flexui version` | Print the compiler version and target platform. |

## Options

| Option | Meaning |
| --- | --- |
| `--config <file>` | Config file. Defaults to `flexui.config.json` or `flexui.json` in the project directory. |
| `--input <file>` | Entry stylesheet containing `@flexui` directives. |
| `--output <file>` | Output stylesheet path. Parent directories are created. |
| `--content <globs>` | Comma-separated glob patterns to scan, overriding the config. |
| `--safelist <list>` | Comma-separated classes that must always be emitted, even when no markup references them. |
| `--minify` | Emit compact CSS. |
| `--watch` | Rebuild on change; equivalent to the `watch` command. |
| `--no-cache` | Ignore the build cache. |
| `--verbose` | Print scan and cache details. |

Flags override the config file only when they are actually passed, so a config
value is not silently replaced by a flag default.

## Config file

```json
{
  "content": ["./index.html", "./src/**/*.{ts,tsx}"],
  "output": "dist/flexui.css",
  "input": "src/styles.css",
  "minify": true,
  "safelist": ["dark:bg-slate-950"],
  "base": true
}
```

| Field | Meaning |
| --- | --- |
| `content` | Globs or directories to scan. Supports `**`, `*`, `?` and `{a,b}`. |
| `output` | Output stylesheet. |
| `input` | Entry stylesheet with `@flexui` directives. |
| `minify` | Compact output. |
| `safelist` | Classes that must be emitted regardless of scanning. |
| `base` | Emit the preflight reset. Defaults to `true`. |

Unknown fields are rejected, so a typo is reported instead of ignored. A UTF-8
BOM is tolerated.

## Directives

```css
@flexui base;        /* preflight reset */
@flexui utilities;   /* generated utilities */
```

Any other `@flexui` directive is an error naming the file. Directives are
removed from the output; the rest of the entry file is appended verbatim after
the generated layers.

## Skip rules

`node_modules`, `.git`, `dist`, `build`, `out`, `.next`, `.nuxt`,
`.svelte-kit`, `coverage`, `.cache`, `vendor`, `target`, `tmp` and other build
directories are never traversed.

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Success. |
| `1` | Build or configuration failure. |
| `2` | Usage error (unknown command, bad flag). |

## The build cache

Results are keyed by a hash of the compiler version, the relevant options and
every scanned file's contents, then stored in `.flexui/cache/<key>.css` with a
small metadata sidecar. Changing any scanned file, the config, the safelist or
`--minify` produces a new key, so a stale stylesheet can never be served.
`flexui clean` removes the whole `.flexui` directory.
