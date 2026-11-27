# flexui-cli-windows-x64

The FlexUI compiler for 64-bit Windows, packaged as a single executable.

You normally never install this directly. `flexui-compiler` declares it as an
optional dependency, npm picks it when you are on Windows x64, and the wrapper
verifies the binary's SHA-256 before running it.

```bash
npm install --save-dev flexui-compiler
npx flexui build
```

If you manage binaries yourself, install this package and point the wrapper at
the executable:

```bash
npm install flexui-cli-windows-x64
set FLEXUI_BINARY=node_modules\flexui-cli-windows-x64\flexui.exe
```

The binary is self-contained: no Go toolchain, no runtime, no DLLs.

## License

MIT
