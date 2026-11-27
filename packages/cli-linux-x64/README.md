# flexui-cli-linux-x64

The FlexUI compiler for 64-bit Linux, shipped as a statically linked binary.

Installing `flexui-compiler` pulls this package in on Linux x64 through its
optional dependencies, checks the SHA-256 recorded in the wrapper, and runs it.
You do not need Go, a C toolchain or any system library.

```bash
npm install --save-dev flexui-compiler
npx flexui build
```

Direct use, for Docker images or CI caches:

```bash
npm install flexui-cli-linux-x64
export FLEXUI_BINARY="$PWD/node_modules/flexui-cli-linux-x64/flexui"
"$FLEXUI_BINARY" version
```

The same binary powers `flexui watch` in long-running dev containers.

## License

MIT
