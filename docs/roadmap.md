# Roadmap

Phases follow the product specification. A phase is marked complete when it is
implemented, tested and documented; nothing here is a placeholder.

## Phase 1 — Compiler foundation ✅

- Go module, monorepo layout, zero third-party dependencies.
- CLI: `build`, `watch`, `inspect`, `init`, `clean`, `version`.
- Scanner with `**` globs and build-directory pruning; markup and script
  extractors built on state machines.
- Bracket-aware candidate parser, AST/IR, utility registry, variant resolver,
  generator, optimizer.
- Deterministic output, content-addressed cache, polling watch mode.
- Precompiled binaries for six platforms with checksums and an npm wrapper.

## Phase 2 — Utility framework ✅

Spacing, layout, flexbox, grid, sizing, position, colour, borders, effects,
transforms, interaction, typography scales, multi-column (`columns-*`,
`break-*`), logical placement (`place-self-*`, `place-content-*`) and semantic
colour aliases that follow the active theme.

Remaining gaps, tracked below: `divide-*` opacity modifiers and broader
logical-property coverage for RTL-first layouts.

## Phase 3 — Variant engine ✅

Responsive breakpoints, `max-*` and arbitrary `min-[…]`/`max-[…]` ranges,
interaction, state, structural, group/peer, accessibility, `dark`/`light`,
`data-[…]`, `aria-[…]`, `has-[…]`, `supports-[…]`, arbitrary selectors, container
queries, composition, configuration-defined custom variants and deterministic
ordering.

Remaining gaps: layout presets (`presets/`: dashboard, sidebar, docs,
e-commerce shells) and named group/peer targeting beyond a single name.

## Phase 4 — Theme and typography engine ✅

- Every design token is emitted as a CSS custom property in `flexui.theme`.
- Semantic roles (`--flexui-surface`, `--flexui-fg`, `--flexui-line`, …) drive
  components and utilities, so theme changes need no rebuild.
- Config overrides for colour, spacing, type, radius, shadow, opacity, z-index,
  breakpoints, easing, duration, aspect ratio, max-width and container sizes.
- Dark-mode strategies: `class` (default), `media` and `both`.
- Local `@font-face` output and font-family/weight/line-height/letter-spacing
  scales.
- Reduced-motion strategy is configurable (`safe` / `off`).

## Phase 5 — Animation engine ✅

- 25 built-in animations (fade, slide, zoom, rotate, bounce, ping, spin, shake,
  wobble, float, shimmer, blur, flip and directional slides).
- Controls: duration, delay, easing, iteration, direction, fill mode, play state
  and composition, plus `animate-none`.
- Staggered children, view-timeline and scroll-timeline utilities, custom and
  arbitrary keyframes from configuration, deduplicated `@keyframes`.
- One-shot reduced-motion override and cubic-bezier/spring easing presets.

## Phase 6 — Component system ✅

150+ built-in classes in `@layer flexui.components` for foundation, forms,
navigation, overlays, data display and layout shells, emitted only when used.

The component engine is open: `components` in `flexui.config.json` defines any
number of components with base declarations, variants, parts, states and
`extends` (custom or built-in), validated at load time and resolved through the
same pipeline as built-ins. A custom component may replace a built-in of the
same name. See [custom components](custom-components.md).

## Phase 7 — Runtime and integrations ✅

- Optional zero-dependency ESM runtime: dialog, drawer, dropdown, popover, tabs,
  toast, command palette, focus management, theme switching, zero-config
  `flexui-runtime/auto` entry, and binder/getter APIs for framework wrappers.
- Framework helpers: `flexui-react` hooks, `flexui-vue` composables and
  `flexui-svelte` actions.
- `flexui-vite` plugin for dev (watch) and build, with compiler error
  reporting and Vite-native configuration.

## Phase 8 — Optimisation and release ✅

- Incremental builds: `.flexui/manifest.json` stores per-file size, mtime,
  content hash and candidates, so a rebuild after one edit reads one file. The
  manifest prunes deleted files, validates its version and is bypassed by
  `--no-cache`. Measured in [benchmarks](benchmarks.md).
- Whole-build cache plus deterministic output; keyframes are emitted once.
- `version.json` is the single version source, resolved by `cmd/version` (Go) and
  `scripts/version.ts` (npm scripts, tags, release notes).
- GitHub Actions: CI (format, vet, tests, typecheck, cross-compiled binaries) and
  a tag-driven release that builds binaries, attaches uniquely named assets,
  writes the changelog as release notes, and publishes npm packages when
  `NPM_TOKEN` is configured.

Remaining: comparison benchmarks against other frameworks on equivalent
projects, a generated documentation site, and filesystem-event watch mode.

## Known limitations

- `filter` utilities (`blur-*`, `grayscale-*`) do not compose with each other;
  transforms do, through CSS custom properties.
- `ring` and `shadow` compose only through the shared `box-shadow` expression.
- Watch mode polls every 250 ms rather than subscribing to filesystem events.
- The incremental manifest trusts size + modification time; a rewrite that
  preserves both is not detected (`--no-cache` bypasses it).
- Warm and incremental builds still scan and stat every file; the numbers are
  in [benchmarks](benchmarks.md).
- No layout preset package yet; the examples compose utilities directly.
