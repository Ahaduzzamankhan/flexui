# flexui-cli-darwin-x64

The FlexUI compiler for Intel Macs.

You get it automatically by installing `flexui-compiler` on macOS x64: npm
selects this optional dependency, the wrapper checks the packaged SHA-256, and
the binary runs directly with no Rosetta translation and no Go install.

```bash
npm install --save-dev flexui-compiler
npx flexui build
```

Using it on its own:

```bash
npm install flexui-cli-darwin-x64
export FLEXUI_BINARY="$PWD/node_modules/flexui-cli-darwin-x64/flexui"
```

Because the executable is unsigned by Apple, the npm wrapper ships its checksum
in `flexui-compiler/checksums.json` so integrity is verifiable without a
notarisation round trip.

## License

MIT
