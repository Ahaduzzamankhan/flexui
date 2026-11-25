# Roadmap

Phases follow the product specification. A phase is only marked complete when it
is implemented, tested and documented — no placeholder features.

## Phase 1 — Compiler foundation ✅

- Go module, monorepo layout, zero third-party dependencies.
- CLI: `build`, `watch`, `inspect`, `init`, `clean`, `version`.
- Scanner with `**` globs and build-directory pruning.
- Markup and script extractors built on state machines.
- Bracket-aware candidate parser, AST/IR, utility registry (≈300 utilities),
  variant resolver, generator, optimizer.
- Content-addressed cache, polling watch mode, deterministic output.
- Precompiled binaries for six platforms with checksums and an npm wrapper.
- 90+ unit and integration assertions.

## Phase 2 — Utility framework (partly done)

Implemented in Phase 1 because the compiler needs real utilities to be
verifiable: spacing, layout, flexbox, grid, sizing, position, colour, borders,
effects, transforms, interaction, typography scales.

Still open:

- Split ratios (`columns-*`), multi-column utilities.
- `divide` opacity and `place-*` completeness review.
- Broader logical-property coverage for RTL-first layouts.

## Phase 3 — Variant engine (partly done)

Implemented: responsive, interaction, state, structural, group/peer,
accessibility, theme, `data-*`, `aria-*`, `has-*`, `supports-*`, arbitrary
selectors, container queries, composition and deterministic ordering.

Still open:

- Named group/peer targeting beyond a single name.
- `max-*` and range breakpoint variants.
- User-defined variants from configuration.
- Layout presets (`presets/`): dashboard, sidebar, docs, e-commerce shells.

## Phase 4 — Theme and typography engine

- Config-driven theme overrides (spacing, colour, type, radius, shadow,
  breakpoints, animation tokens).
- CSS custom property output for tokens, with runtime theme switching.
- System theme (`prefers-color-scheme`) strategy alongside the class strategy.
- Variable font and font-feature configuration, local font declarations.

## Phase 5 — Animation engine

- Built-in animations (fade, slide, zoom, rotate, bounce, ping, spin, shake,
  wobble, float, shimmer, blur, flip).
- Animation controls: duration, delay, iteration, direction, fill mode,
  play state, easing, composition.
- Custom and arbitrary keyframes with deduplication.
- Staggered children, scroll-driven animations, view transitions where
  supported.
- Reduced-motion handling and cubic-bezier/spring presets.

## Phase 6 — Component system ✅

Foundation, forms, navigation, overlays, data display and complex patterns, each
with semantic HTML, FlexUI CSS, accessible behaviour, responsive layout, theme
support and documented customisation. Components do not require a large runtime;
overlays progressively enhance.

Implemented as plain classes in `@layer flexui.components`, emitted only when
used: buttons, badges, cards, stats, avatars, forms, navigation, tabs, tables,
timelines, trees, calendars, charts, dialogs, drawers, popovers, toasts, command
palettes and layout shells. See [components](components.md).

## Phase 7 — Runtime and integrations

- ✅ Optional zero-dependency vanilla JS runtime, importable per module
  (dialog, drawer, dropdown, popover, tabs, toast, command palette, focus
  management, theme switching). See [runtime](runtime.md).
- Vite plugin with incremental compilation and error reporting.
- React, Vue and Svelte helpers where needed.

## Phase 8 — Optimisation and release

- Per-file incremental compilation and cache pruning.
- Keyframe deduplication and further CSS size work.
- Benchmarks against other utility frameworks on realistic projects (cold,
  warm, incremental build time; peak memory; binary, CSS and runtime size;
  large-project scaling). No performance claim will be published before that.
- npm publishing pipeline, release automation, documentation site.

## Known limitations today

- Filters do not compose; transforms do.
- `ring` and `shadow` compose only through the shared `box-shadow` expression.
- Watch mode polls (250 ms) rather than using filesystem events.
- Dark mode is class-based only.
- No layout preset package yet; the example uses utilities directly.
