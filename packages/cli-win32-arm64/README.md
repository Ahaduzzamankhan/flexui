# flexui-cli-windows-arm64

The FlexUI compiler for Windows on ARM, packaged as a single executable for
devices such as Snapdragon-based laptops.

`flexui-compiler` lists this package as an optional dependency, so a normal
install on Windows ARM64 pulls it in automatically and verifies its checksum
before the first build.

```bash
npm install --save-dev flexui-compiler
npx flexui build --minify
```

To wire it up by hand:

```bash
npm install flexui-cli-windows-arm64
set FLEXUI_BINARY=node_modules\flexui-cli-windows-arm64\flexui.exe
```

Nothing else is required on the machine; the executable contains the whole
compiler.

## License

MIT
