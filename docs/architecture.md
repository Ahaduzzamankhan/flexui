# Architecture

## Pipeline

```
source files
  → scanner      discover files (globs, ignore rules, deterministic order)
  → extractor    markup/script state machines produce raw candidates
  → parser       bracket-aware split into variants + utility
  → ast          structural rules, not CSS strings
  → resolver     utilities registry + variants → rules
  → generator    sort, group by at-rule context, emit CSS
  → optimizer    deduplicate rules and declarations, drop empty rules
  → writer       buffered output
```

Two invariants shape the design:

1. **CSS text is produced in exactly one place.** Parsing and resolution never
   concatenate CSS. A rule is data (`ast.Rule`), so it can be sorted, deduped
   and tested without string surgery.
2. **Output order is a total order.** Every rule carries a utility order number
   and a variant rank, and the generator sorts by
   (group-rule rank, group-rule key, utility order, variant rank, variant count,
   class name, declaration signature). Sorting is `SliceStable`, so equal keys
   keep resolution order, which itself follows the sorted candidate list.

## Packages

| Package | Responsibility |
| --- | --- |
| `internal/scanner` | Walks the project once, prunes build directories, matches `**`-aware globs, returns sorted absolute paths. |
| `internal/extractor` | Tag/attribute state machine for markup (HTML, JSX, Vue, Svelte, Astro, templating files) and a string/comment/template state machine for scripts. |
| `internal/parser` | Splits a candidate on top-level colons respecting `[]`, `()`, `{}`, quotes and escapes; handles `!important`, negatives, arbitrary properties. |
| `internal/ast` | `Decl`, `Block`, `AppliedVariant`, `Rule`. |
| `internal/themes` | Spacing, colour palettes, type scales, radii, shadows, breakpoints, easings. |
| `internal/utilities` | Utility registry: exact names and longest-prefix matches, each with a stable cascade order. |
| `internal/variants` | Variant resolution into selector templates and at-rules, with ordering ranks. |
| `internal/resolver` | Combines both, applies `!important`, produces `ast.Rule`. |
| `internal/generator` | Selector construction (with CSS identifier escaping), at-rule grouping, preflight, formatting, minified or pretty. |
| `internal/optimizer` | Rule/declaration deduplication and conservative minification of third-party CSS. |
| `internal/cache` | Content-addressed `.css` + metadata entries. |
| `internal/watcher` | Snapshot/diff polling watcher. |
| `internal/compiler` | Options, pipeline, cache keys, input directives, stats. |
| `internal/cli` | Commands, flag/config merging, watch loop. |
| `internal/config` | `flexui.config.json` loading with strict fields and BOM tolerance. |

Nothing outside `cmd/flexui` depends on a global; the compiler is usable as a
library through `compiler.Build`.

## Selector construction

A variant is a template containing `&`. Templates are applied in written order,
each replacing `&` with the accumulated expression; the class name replaces the
final `&`:

```
md:group-hover:p-4
  start:            &
  group-hover:      .group:hover &
  final selector:   .group:hover .md\:group-hover\:p-4
  media wrapper:    @media (min-width: 768px)
```

Selectors escape punctuation with backslash and use hexadecimal escapes for a
leading digit, per CSS syntax.

## Determinism

- Files are sorted before reading.
- Candidates are de-duplicated in a set, then sorted.
- Resolution runs concurrently but writes into an index-addressed slice, so
  scheduling cannot affect order.
- Generation sorts by a total key and emits groups in that order.
- Maps are never iterated when producing output.

The test suite asserts byte-equality between two builds and between builds with
different candidate orders.

## Concurrency

`compiler.parallel` runs a fixed pool of `runtime.NumCPU()` workers for reading,
extraction and resolution. Concurrency is only used where the work is
independent and large enough to pay for the scheduling; extraction per file and
resolution per candidate both qualify, selector construction does not.

## Caching

The cache key hashes the compiler version, the options that affect output, the
passthrough CSS, the safelist and every scanned file's relative path and
contents. A hit skips extraction, resolution and generation entirely and returns
the stored stylesheet plus its metadata, so reported statistics stay accurate.

## Layers

Generated CSS uses two cascade layers:

```css
@layer flexui.base, flexui.utilities;
```

Preflight lives in `flexui.base`, utilities in `flexui.utilities`. Unlayered
author CSS therefore always wins over generated utilities, and passthrough CSS is
appended after the layers.

## Testing

`go test ./...` covers the parser, extractor, scanner/glob matcher, utility
values, variants, generator (escaping, ordering, minified output), optimizer,
cache, watcher, config and end-to-end builds including cache invalidation and
determinism.
