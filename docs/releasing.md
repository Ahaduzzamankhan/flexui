# Versioning, CI and releases

## One version, one file

[`version.json`](../version.json) is the only place the version is written by
hand:

```json
{
  "name": "flexui",
  "version": "2.0.0-alpha.2",
  "channel": "alpha"
}
```

Everything else resolves it:

| Consumer | Resolver |
| --- | --- |
| Go compiler and build scripts | `go run ./cmd/version [-field version\|tag\|name\|channel]` |
| npm packages, release notes, tags | `bun scripts/version.ts` |
| Release binaries | `-ldflags "-X .../internal/compiler.Version=<version>"` in `scripts/build-binaries.sh` |

`internal/compiler/version_test.go` fails when `compiler.Version` and
`version.json` disagree, so a plain `go build` can never silently ship a stale
version.

## Bumping a version

1. Edit `version.json`.
2. Add a `## <version>` section to [`CHANGELOG.md`](../CHANGELOG.md). Releases
   without a section are refused, so release notes always exist.
3. Run `bun run version:sync`. This rewrites `"version"` in the root and every
   `packages/*/package.json`, plus internal `flexui-*` dependency pins.
4. Run `bun run version:check`. It reports any manifest or internal dependency
   that still disagrees and whether the changelog section exists.
5. Commit the changes.

Update the default in `internal/compiler/compiler.go` only if you want
`go run` to report the new version without linking flags; the version test tells
you when it is out of date.

## Local release preparation

```bash
bun run release              # build binaries, write dist-release/, no tag
bun run release:publish      # tag, push, gh release create
```

`scripts/release.ts`:

1. reads `version.json`;
2. refuses to continue if any package version, internal dependency pin or
   changelog section is missing;
3. cross-compiles all six binaries (`scripts/build-binaries.sh`, version
   injected with `-ldflags`);
4. stages every binary under `dist-release/assets/` with a unique name
   (`flexui-linux-x64`, `flexui-win32-x64.exe`, …), because a GitHub release
   identifies assets by file name and six files called `flexui` would collide;
5. writes `dist-release/RELEASE_NOTES.md` (the changelog section plus the asset
   list) and `dist-release/SHA256SUMS`;
6. with `--publish`, checks that the working tree is clean, creates the
   annotated tag `v<version>`, pushes it, and runs
   `gh release create ... --notes-file dist-release/RELEASE_NOTES.md` with the
   staged binaries, `checksums.json` and `SHA256SUMS` attached.

Nothing is tagged or pushed without `--publish`.

## CI

`.github/workflows/ci.yml` runs on every push to `main` and every pull request:

- **Go compiler** — `gofmt` check, `go vet`, `go test ./...`, `go build ./...`.
- **TypeScript packages** — `bun install --frozen-lockfile`,
  `bun run version:check`, `bun run typecheck`, and the runtime and Vite plugin
  tests.
- **Cross-compiled binaries** — builds all six platforms, asserts that the
  produced binary reports the version from `version.json`, and uploads the
  binaries and checksum manifest as a workflow artifact.

## Release workflow

`.github/workflows/release.yml` runs when a `v*` tag is pushed:

1. **Verify** — the pushed tag must equal `bun scripts/version.ts tag`, then
   version consistency, vet, Go tests, typecheck and TypeScript tests run.
2. **Release** — builds the binaries, writes notes from `CHANGELOG.md`, and
   creates the GitHub release with the binaries, `checksums.json` and
   `SHA256SUMS` attached.
3. **npm** — publishes with `scripts/publish.ts` when `NPM_TOKEN` is set. The
   script validates every package (name, non-private, version match, build
   output, CLI checksums) before publishing anything, then publishes the
   platform packages first, the runtime and framework integrations next, and the
   CLI last, so `optionalDependencies` resolve by the time the CLI is installed.
   Prereleases use an explicit dist-tag (`alpha`, from `version.json`'s channel)
   because npm refuses to publish a prerelease without one.

The npm job authenticates in one of two ways:

- **Trusted publishing (preferred)** — when the repository is registered as a
  trusted publisher on npmjs.com, the job exchanges its GitHub OIDC token for a
  short-lived credential, attaches provenance automatically and needs no secret.
- **`NPM_TOKEN`** — a granular access token with publish rights, used when OIDC
  is not configured.

Either path runs plain `npm publish`. Staged publishing is a separate, explicit
flow (`npm stage publish`, or a trusted publisher configured as stage-only); this
repository never uses it, so a release version becomes installable as soon as the
registry accepts the tarball. If a brand-new package is held back by npm's review
of first-time packages, the placeholder `0.0.0-stage` appears until that review
finishes, and the remaining packages are unaffected.

Without a credential the release still publishes binaries on GitHub and the npm
job reports that it skipped.

Validate locally without publishing:

```bash
bun run build:js && bun run build:runtime && bun run build:integrations
bun run publish:npm:dry-run
```

## Verifying a release

```bash
./packages/cli-linux-x64/flexui version       # must print the released version
sha256sum -c dist-release/SHA256SUMS          # from the repository root
bun run version:check
```

The compile itself is reproducible: the same sources, configuration and flags
produce byte-identical CSS, and the cache key covers the compiler binary
fingerprint, configuration and every scanned file.
