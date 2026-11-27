# Benchmarks

These are measurements of FlexUI's own build pipeline. No comparison against
other frameworks is published yet, because that requires an equivalent project
definition per framework and a fair hardware setup; until that exists, no
relative performance claim is made.

## Method

```bash
go test -run='^$' -bench='BenchmarkBuild' -benchmem -benchtime=2s ./internal/compiler
```

- Fixture: 200 HTML files in `src/group*/page*.html`, each using two of five
  realistic class sets (layout, grid, typography, surfaces, component classes).
- `BenchmarkBuildCold` disables both caches (`NoCache`), so every run scans,
  reads, extracts, resolves, generates and optimizes.
- `BenchmarkBuildWarmCache` runs a build that is unchanged, so the whole-build
  cache hits.
- `BenchmarkBuildIncremental` edits one file between iterations; the other 199
  files are reused from the incremental manifest, then the stylesheet is
  regenerated and both caches are written.
- Machine: Intel Core i5-6500 @ 3.20 GHz, Windows, amd64, Go toolchain pinned by
  `go.mod`, `GOMAXPROCS=4`.

## Results

| Benchmark | Time per build | Allocated | Allocations |
| --- | --- | --- | --- |
| Cold (no cache) | 21.7 ms | 1.08 MB | 8,221 |
| Warm (whole-build cache hit) | 7.1 ms | 0.60 MB | 5,832 |
| Incremental (1 of 200 files edited) | 20.4 ms | 1.26 MB | 7,754 |

Reproduce with the command above; numbers vary with hardware and filesystem, so
compare runs on the same machine rather than across machines.

## What the incremental work actually bought

The first implementation wrote one cache file per source file. Measured on this
machine, loading 200 of those entries took **26 ms**, while extracting all 200
files took **0.66 ms** — the per-file disk cache made rebuilds slower, so it was
replaced.

The current design keeps one `.flexui/manifest.json` holding each file's size,
modification time, content hash and extracted candidates. Against the earlier
per-file implementation:

| Scenario | Per-file cache | Single manifest |
| --- | --- | --- |
| Warm build (200 files) | 17.8 ms | 7.1 ms |
| Incremental (1 edit) | 47.8 ms | 20.4 ms |

The warm path is now dominated by scanning and stat-ing files plus reading the
manifest, not by parsing. The incremental path still pays for CSS regeneration,
the build-cache write and the manifest write, which is why it is close to a cold
build at this project size; the number of files read falls from 200 to 1.

The manifest identity check uses size and modification time. A rewrite that
preserves both (for example a restored backup timestamp on an identically sized
file) is not detected; `flexui build --no-cache` and `flexui clean` bypass the
manifest entirely.

## Stylesheet size

Sizes depend on the utilities used, not on the framework: each class is emitted
once, in one layer, and unused utilities are never emitted. The example sites'
generated stylesheets are not a meaningful benchmark on their own, so no size
table is published here.
