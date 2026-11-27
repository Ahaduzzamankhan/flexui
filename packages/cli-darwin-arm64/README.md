# flexui-cli-darwin-arm64

The FlexUI compiler for Apple silicon.

Installing `flexui-compiler` on an M-series Mac pulls this package in as an
optional dependency, verifies the binary checksum, and runs it natively.

```bash
npm install --save-dev flexui-compiler
npx flexui build --minify
```

Standalone use:

```bash
npm install flexui-cli-darwin-arm64
export FLEXUI_BINARY="$PWD/node_modules/flexui-cli-darwin-arm64/flexui"
```

The executable is self-contained and needs no developer tools on the machine.

## License

MIT
