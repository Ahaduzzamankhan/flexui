#!/usr/bin/env bash
# Cross-compiles the FlexUI compiler for every published platform, vendors each
# binary into its platform package and writes the checksum manifest used by the
# npm wrapper for integrity verification.
#
# Contributors need Go for this script. End users never do: they install the
# published binaries through npm.
set -euo pipefail

cd "$(dirname "$0")/.."
ROOT="$(pwd)"
OUT="$ROOT/packages/cli/binaries"
# version.json is the single source of truth; cmd/version resolves it without
# needing jq or any other tool.
VERSION="$(go run ./cmd/version -field version)"
if [ -z "$VERSION" ]; then
  echo "could not resolve the version from version.json" >&2
  exit 1
fi
LDFLAGS="-s -w -X github.com/Ahaduzzamankhan/flexui/internal/compiler.Version=$VERSION"
echo "version $VERSION"

mkdir -p "$OUT"
printf '%s\n' "{" > "$ROOT/packages/cli/checksums.json"

first=1
for target in "windows amd64 win32-x64 flexui.exe" \
              "windows arm64 win32-arm64 flexui.exe" \
              "linux amd64 linux-x64 flexui" \
              "linux arm64 linux-arm64 flexui" \
              "darwin amd64 darwin-x64 flexui" \
              "darwin arm64 darwin-arm64 flexui"; do
  set -- $target
  goos="$1"; goarch="$2"; key="$3"; name="$4"
  dest="$ROOT/packages/cli-$key"
  mkdir -p "$dest"
  echo "building $goos/$goarch -> packages/cli-$key/$name"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build -trimpath -ldflags "$LDFLAGS" -o "$dest/$name" ./cmd/flexui
  mkdir -p "$OUT/$key"
  cp "$dest/$name" "$OUT/$key/$name"
  sum="$(sha256sum "$dest/$name" | cut -d' ' -f1)"
  if [ $first -eq 0 ]; then printf ',\n' >> "$ROOT/packages/cli/checksums.json"; fi
  first=0
  printf '  "%s": { "file": "%s", "sha256": "%s", "version": "%s" }' \
    "$key" "$name" "$sum" "$VERSION" >> "$ROOT/packages/cli/checksums.json"
done
printf '\n}\n' >> "$ROOT/packages/cli/checksums.json"

echo "done. checksums written to packages/cli/checksums.json"
