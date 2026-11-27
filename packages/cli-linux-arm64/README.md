# flexui-cli-linux-arm64

The FlexUI compiler for 64-bit ARM Linux, for servers and development machines
running on ARM.

`flexui-compiler` depends on this package optionally, so a normal install on
Linux ARM64 resolves it without extra steps and verifies its checksum before the
first run.

```bash
npm install --save-dev flexui-compiler
npx flexui init
npx flexui build
```

Manual setup:

```bash
npm install flexui-cli-linux-arm64
export FLEXUI_BINARY="$PWD/node_modules/flexui-cli-linux-arm64/flexui"
```

The binary is statically linked and needs no runtime libraries.

## License

MIT
