# Changelog

All notable changes to FlexUI are documented here. The release workflow reads
the section for the version in `version.json` and uses it as the GitHub release
body, so every release needs an entry here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and
the project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-11-27

FlexUI 2.0 is out, installable with one command.

### Added

- `flexui-v2` installs the compiler and the matching platform binary with a
  single `npm install flexui-v2` (or `npm install -g flexui-v2`). npm picks the
  binary package from its `os` and `cpu` fields, so Windows, Linux and macOS on
  x64 or ARM64 all resolve to the right executable.
- Every npm package ships its own README, and the repository README carries a
  complete installation guide.

### Fixed

- Framework packages import `flexui-runtime` under its published name.
- The wrapper also finds a platform package installed beside the compiler, which
  covers nested and global installs.
- CI validates package contents and publishable tarballs before a release.

## [2.0.0-alpha.3] - 2026-11-27

One command to install everything.

### Added

- `flexui-v2` installs the compiler and the matching platform binary with a
  single `npm install flexui-v2` (or `npm install -g flexui-v2`). npm picks the
  binary package from its `os` and `cpu` fields, so Windows, Linux and macOS on
  x64 or ARM64 all work without choosing a package by hand.

### Fixed

- The wrapper also looks for a platform package installed alongside the
  compiler, which covers nested installs.

## [2.0.0-alpha.2] - 2026-11-27

Fixes npm packaging for the first alpha.

### Fixed

- The npm wrapper now resolves the Windows compiler from
  `flexui-cli-windows-x64` and `flexui-cli-windows-arm64`. The first published
  `flexui-compiler` tarball still asked for the old `win32` package names, so
  `flexui` could not find its binary on Windows.

## [2.0.0-alpha.1] - 2026-11-26

First public alpha. Everything listed is implemented, tested and documented;
nothing is a placeholder.

### Added

- **Compiler** — a single native binary (Go standard library only) that scans
  HTML, JSX/TSX, Vue, Svelte, Astro, Markdown and nine other template formats,
  resolves only the classes found in those sources, and writes one deterministic
  stylesheet. Deterministic ordering is enforced by layer, at-rule, utility,
  variant rank and class name.
- **Utilities** — spacing, layout, flexbox, grid, sizing, position, colour,
  typography, borders, effects, transforms, interaction, multi-column and
  logical placement utilities, with semantic colour aliases that follow the
  active theme.
- **Variants** — responsive breakpoints, `max-*` and `min-[…]`/`max-[…]` ranges,
  interaction and state variants, group/peer, accessibility, `dark`/`light`,
  `data-[…]`, `aria-[…]`, `has-[…]`, `supports-[…]`, arbitrary selectors,
  container queries and configuration-defined custom variants.
- **Theme engine** — every design token is emitted as a CSS custom property in
  `@layer flexui.theme`, with semantic surface/foreground roles, `@font-face`
  output, and `class`, `media` and `both` dark-mode strategies. Theme switching
  needs no rebuild.
- **Animation engine** — 25 built-in animations plus duration, delay, easing,
  iteration, direction, fill, play-state and composition controls, staggered
  children, view-timeline and scroll-timeline utilities, custom keyframes from
  configuration, deduplicated `@keyframes`, and a one-shot reduced-motion
  override.
- **Component system** — 150+ `fui-*` classes in `@layer flexui.components` for
  foundation, forms, navigation, overlays, data display and layout shells. Only
  classes used in scanned sources are emitted.
- **Custom component engine** — define any number of components in
  `flexui.config.json` with base declarations, variants, parts, state selectors
  and `extends` (custom or built-in). Definitions are validated at load time,
  resolve through the same pipeline as built-ins, can replace a built-in of the
  same name, and produce deterministic CSS.
- **Runtime** — optional zero-dependency ESM `flexui-runtime` for dialogs,
  drawers, dropdowns, popovers, tabs, toasts, command palettes, focus trapping
  and theme switching, with a zero-config `flexui-runtime/auto` entry and
  framework-agnostic binder/getter APIs.
- **Integrations** — `flexui-vite` (build and watch, error reporting),
  `flexui-react` hooks, `flexui-vue` composables and `flexui-svelte` actions.
- **Incremental builds** — `.flexui/manifest.json` records each file's size,
  modification time, content hash and extracted candidates, so a rebuild after
  one edit reads one file instead of the whole project. The manifest prunes
  deleted files, validates its own version, and is bypassed by `--no-cache`.
- **Distribution** — precompiled binaries for Windows, Linux and macOS on x64 and
  ARM64 with SHA-256 checksums, resolved by a TypeScript npm wrapper that needs
  no Go toolchain and runs no install scripts. Release assets use unique
  per-platform names, and npm packages publish in dependency order with an
  explicit dist-tag for prereleases.
- **Tooling** — `version.json` as the single version source, `bun run
  version:sync` / `version:check`, `scripts/release.ts` for tagged releases with
  binaries and changelog notes, and GitHub Actions CI and release workflows.

### Changed

- The stylesheet layer statement lists only layers that are actually emitted, so
  a build without preflight no longer declares `flexui.base`.
- Overlay scroll locking uses the `data-fui-overlay-open` attribute and is
  reference counted across stacked dialogs and drawers.

### Fixed

- `min-[700px]` and other bracketed range variants resolved to a breakpoint
  lookup instead of a raw length.
- `initDialogs` and `initDrawers` bound controllers only when disposed, so
  `init()` never actually activated them.
- `initTheme({ storageKey })` persisted theme changes to the default key instead
  of the configured one.

### Known limitations

- `filter` utilities do not compose with each other; transforms do.
- Watch mode polls (250 ms) rather than using filesystem events.
- Warm and incremental builds are still bound by whole-project scanning and
  cache writes; measured numbers are in `docs/benchmarks.md`.
- No published comparison against other frameworks yet.

[2.0.0-alpha.1]: https://github.com/Ahaduzzamankhan/flexui/releases/tag/v2.0.0-alpha.1
