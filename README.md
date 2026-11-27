# FlexUI 2.0

A utility-first CSS compiler, design-token system and (eventually) full UI
ecosystem, written in Go with zero external dependencies.

FlexUI scans your markup and scripts, resolves only the utilities you actually
use and writes one deterministic, standards-compliant stylesheet. The compiler
is a native binary; you never install Go to use it.

- **Status:** 2.0.0-alpha.2. The compiler, utilities, variants, theme and
  animation engines, component layer (including user-defined components),
  optional runtime, framework integrations, incremental builds, benchmarks and
  release automation are implemented and tested. See
  [docs/roadmap.md](docs/roadmap.md) for what remains.
- **Version:** one source of truth in [`version.json`](version.json); tags,
  binaries, npm packages and release notes all resolve it.
- **License:** MIT

## Installation

FlexUI installs with the Node.js toolchain you already have. The compiler is a
native binary, so Go is not part of the setup.

**Requirements:** Node.js 18 or newer, or Bun 1.1 or newer.

### One command

```bash
npm install flexui-v2
```

`flexui-v2` installs the compiler and the binary for your machine in one step.
npm picks the right platform package from its `os` and `cpu` fields, so Windows,
Linux and macOS — on x64 or ARM64 — all work without a platform choice:

```bash
npm install --save-dev flexui-v2     # project
npm install -g flexui-v2             # global `flexui` command
pnpm add -D flexui-v2
yarn add -D flexui-v2
bun add -d flexui-v2
```

When you only want the compiler inside one project, `flexui-compiler` is the
same thing without the umbrella package:

```bash
npm install --save-dev flexui-compiler
```

### Create a project

```bash
npx flexui init
```

With a global install, drop `npx` and call `flexui init` directly.

`init` writes `flexui.config.json` and an entry stylesheet. Send the `content`
globs at your source files, then compile:

```bash
npx flexui build
```

The stylesheet is written to `dist/flexui.css` by default. While you work,
`npx flexui watch` recompiles on every change.

### Verify the install

```bash
npx flexui version
npx flexui inspect
```

`version` prints the compiler version, and `inspect` lists the rules a build
would emit — the quickest way to confirm the binary, the configuration and your
content globs are all connected.

### Framework setups

- **Vite, React, Vue, Svelte** — install `flexui-vite` and add `flexui()` to
  `vite.config.ts`. The plugin compiles when the dev server starts and before
  every production build.
- **Interactive components** — install `flexui-runtime`, plus `flexui-react`,
  `flexui-vue` or `flexui-svelte` when you want framework bindings.

### Updating

```bash
npm install --save-dev flexui-compiler@alpha
```

The `alpha` dist-tag follows the 2.0 line. Every release is described in
[CHANGELOG.md](CHANGELOG.md).

## Use

`flexui init` writes `flexui.config.json` and an entry stylesheet:

```css
/* src/styles.css */
@flexui base;
@flexui utilities;
```

Then build:

```bash
npx flexui-compiler build --minify
```

Output is plain CSS. No JavaScript is required for styling.

```css
@layer flexui.base, flexui.utilities;
@layer flexui.utilities {
  .p-4 { padding: 1rem; }
  @media (min-width: 768px) {
    .md\:flex { display: flex; }
  }
}
```

## What works today

**Compiler pipeline** — scanner → extractor → parser → AST → utility resolver →
variant resolver → generator → optimizer, with a content-addressed build cache,
watch mode and stable output ordering.

- Scans `.html`, `.jsx`, `.tsx`, `.js`, `.ts`, `.vue`, `.svelte`, `.astro`, `.md`,
  `.mdx` and nine more template formats, with structural state machines for
  markup and scripts (no regex-driven markup parsing).
- Bracket-aware candidate parser: nested arbitrary values, arbitrary properties
  (`[mask-type:luminance]`), `!important`, negative values and opacity
  modifiers.
- Utility families: spacing, layout, flexbox, grid, sizing, position, colour,
  typography, effects, borders, transforms and interaction.
- Variants: responsive (`sm`…`2xl`), interaction, state, structural,
  group/peer, accessibility, `dark`/`light`, `data-[…]`, `aria-[…]`, `has-[…]`,
  `supports-[…]`, arbitrary selectors, and container queries (`@sm:`, `@[40rem]:`).
- Deterministic output: rules sort by group-rule rank, utility order, variant
  rank, variant count and class name; the same input always yields the same bytes.
- Type scale and semantic classes (`type-display-xl` … `type-code`).

**CLI** — `build`, `watch`, `inspect`, `init`, `clean`, `version` with
`--minify`, `--watch`, `--config`, `--input`, `--output`, `--content`,
`--safelist`, `--no-cache` and `--verbose`.

**Distribution** — six cross-compiled binaries (Windows, Linux, macOS × x64 and
ARM64) with SHA-256 checksums, resolved by a TypeScript npm wrapper that never
needs Go and never runs an install script. The wrapper also exposes a typed
programmatic API (`build`, `watch`, `inspect`, `version`).

**Component library** - 150+ semantic classes in `@layer flexui.components`:
buttons, cards, forms, navigation, tables, timelines, dialogs, drawers,
popovers, toasts and command palettes. Components use the same design tokens as
utilities, and utilities always override them. Only the classes you actually use
are emitted. See [docs/components.md](docs/components.md).

**Custom components** - the component engine is open. Define components in
`flexui.config.json` with base declarations, variants, parts, state selectors and
`extends` (custom or built-in), and use them as `fui-<name>` classes. There is no
cap on how many you define, they are validated at config load time, only the
classes you use are emitted, and a custom component can replace a built-in. See
[docs/custom-components.md](docs/custom-components.md).

**Animation engine** - 25 built-in animations plus duration, delay, easing,
iteration, direction, fill, play-state and composition controls, staggered
children, view- and scroll-timeline utilities, configuration-defined keyframes,
deduplicated `@keyframes` and a reduced-motion override.

**Theme engine** - design tokens emitted as CSS custom properties with semantic
roles, config overrides and `class`/`media`/`both` dark-mode strategies, so theme
switching needs no rebuild.

**Integrations** - `flexui-vite` for dev and build, `flexui-react` hooks,
`flexui-vue` composables and `flexui-svelte` actions.

**Runtime** - `flexui-runtime` is an optional zero-dependency ESM package for
dialog, drawer, dropdown, popover, tabs, toast, command palette, focus
management and theme switching. It progressively enhances native `<dialog>` and
the popover API, restores focus, locks scroll, and emits DOM events. Import the
whole runtime or a single module. See [docs/runtime.md](docs/runtime.md).

**Examples** - `examples/basic` is a complete marketing site (hero, bento
features, pricing, FAQ, footer) built entirely from FlexUI utilities.
`examples/components` is an interactive component lab that exercises every
runtime feature in one page.

## Release and CI

`version.json` is the only version written by hand. `go run ./cmd/version`
resolves it for Go builds, `scripts/version.ts` resolves it for npm scripts,
tags and release notes, and release binaries get it injected with `-ldflags`.

```bash
bun run version:check        # every package manifest matches version.json
bun run release              # build binaries, stage assets, write notes
bun run release:publish      # tag, push and create the GitHub release
bun run publish:npm:dry-run  # verify publishable tarballs without publishing
```

GitHub Actions runs formatting, vet, tests, typecheck and cross-compiled
binaries on every push and pull request. Pushing a `v*` tag verifies the tag
against `version.json`, builds every platform, attaches uniquely named binaries
plus `checksums.json`/`SHA256SUMS`, and uses the matching `CHANGELOG.md` section
as the release body. npm packages publish from that job when `NPM_TOKEN` is set.
See [docs/releasing.md](docs/releasing.md).

## Roadmap

The next milestones — a generated documentation site, comparison benchmarks and
filesystem-event watching — are tracked in [docs/roadmap.md](docs/roadmap.md).
Release notes for every version live in [CHANGELOG.md](CHANGELOG.md).

## Repository layout

```
cmd/flexui/            compiler entry point
internal/
  scanner/             file discovery and glob matching
  extractor/           markup and script state machines
  parser/              candidate parser
  ast/                 intermediate representation
  themes/              design tokens
  utilities/           utility registry
  variants/            variant resolver
  resolver/            candidate → AST
  generator/           CSS generator and preflight
  optimizer/           deduplication and minification
  cache/               content-addressed build cache
  watcher/             polling file watcher
  compiler/            pipeline orchestration
  cli/                 command line interface
  components/          component CSS registry
packages/
  cli/                 npm package "flexui" (TypeScript: resolver + API)
  runtime/             npm package "flexui-runtime" (TypeScript, ESM)
  react/ vue/ svelte/  framework helpers
  vite-plugin/         npm package "flexui-vite"
  cli-<platform>/      optional platform packages with native binaries
examples/basic/        marketing site example, compiled by FlexUI
examples/components/   component and runtime lab
scripts/               serve.ts, build-binaries.sh, version.ts, release.ts, publish.ts
docs/                  documentation
version.json           single version source
CHANGELOG.md           release notes used by the release workflow
.github/workflows/     CI and tag-driven release
```

## Development

Building the compiler requires Go; using it does not.

```bash
go test ./...                     # Go unit and integration tests
go vet ./...                      # static checks
go build -o bin/flexui ./cmd/flexui

bun install                       # TypeScript toolchain (dev only)
bun run typecheck                 # tsc, strict, runtime + examples + wrapper
bun run build:js                  # emit packages/cli/dist
bun run build:runtime             # emit packages/runtime/dist
bun run build:components          # runtime + example TypeScript

bash scripts/build-binaries.sh    # cross-compile all platforms + checksums
```

Build and preview the example sites:

```bash
cd examples/basic
../../bin/flexui build                   # writes dist/flexui.css
cd ../..

cd examples/components
../../bin/flexui build                   # writes dist/flexui.css
cd ../..

bun scripts/serve.ts . 4173              # http://127.0.0.1:4173/examples/components/
```

`scripts/serve.ts` uses Node APIs only, so it runs under Node or Bun.

## Documentation

- [Getting started](docs/getting-started.md)
- [CLI reference](docs/cli.md)
- [Components](docs/components.md)
- [Custom components](docs/custom-components.md)
- [Runtime](docs/runtime.md)
- [Integrations](docs/integrations.md)
- [Benchmarks](docs/benchmarks.md)
- [Releasing](docs/releasing.md)
- [Utility reference](docs/utilities.md)
- [Variants](docs/variants.md)
- [Architecture](docs/architecture.md)
- [Roadmap](docs/roadmap.md)
