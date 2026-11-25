# FlexUI 2.0

A utility-first CSS compiler, design-token system and (eventually) full UI
ecosystem, written in Go with zero external dependencies.

FlexUI scans your markup and scripts, resolves only the utilities you actually
use and writes one deterministic, standards-compliant stylesheet. The compiler
is a native binary; you never install Go to use it.

- **Status:** 2.0.0-alpha.1. The compiler, component CSS layer and optional
  zero-dependency runtime are implemented and tested. Integrations and
  benchmarks remain - see [docs/roadmap.md](docs/roadmap.md).
- **License:** MIT

## Install

```bash
npm install flexui
npx flexui init
npx flexui build
```

The npm package ships precompiled binaries for Windows x64/ARM64, Linux
x64/ARM64 and macOS x64/ARM64. If your platform package is missing, the wrapper
reports exactly what it looked for and how to fix it; you can also point
`FLEXUI_BINARY` at any `flexui` binary.

## Use

`flexui init` writes `flexui.config.json` and an entry stylesheet:

```css
/* src/styles.css */
@flexui base;
@flexui utilities;
```

Then build:

```bash
npx flexui build --minify
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

**Runtime** - `@flexui/runtime` is an optional zero-dependency ESM package for
dialog, drawer, dropdown, popover, tabs, toast, command palette, focus
management and theme switching. It progressively enhances native `<dialog>` and
the popover API, restores focus, locks scroll, and emits DOM events. Import the
whole runtime or a single module. See [docs/runtime.md](docs/runtime.md).

**Examples** - `examples/basic` is a complete marketing site (hero, bento
features, pricing, FAQ, footer) built entirely from FlexUI utilities.
`examples/components` is an interactive component lab that exercises every
runtime feature in one page.

## Not implemented yet

These are specified and deliberately not faked:

- Animation engine (staggered and scroll-driven animation) - Phase 5.
- Vite/React/Vue/Svelte integrations - Phase 7.
- Theme configuration overrides and CSS-variable token output - Phase 4.
- Signal-based watch (current watch mode polls) and per-file incremental
  compilation - Phase 8.
- Published benchmarks against other frameworks - Phase 8. No performance claim
  is made until measurements exist.

Known limitation worth calling out: `filter` utilities (`blur-*`, `grayscale-*`)
do not yet compose with each other, because each emits the `filter` property
directly. Transforms *do* compose, through CSS custom properties.

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
  runtime/             npm package "@flexui/runtime" (TypeScript, ESM)
  cli-<platform>/      optional platform packages with native binaries
examples/basic/        marketing site example, compiled by FlexUI
examples/components/   component and runtime lab
scripts/               serve.ts preview server, build-binaries.sh
docs/                  documentation
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
- [Runtime](docs/runtime.md)
- [Utility reference](docs/utilities.md)
- [Variants](docs/variants.md)
- [Architecture](docs/architecture.md)
- [Roadmap](docs/roadmap.md)
