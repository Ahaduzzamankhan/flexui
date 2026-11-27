# flexui-vue

Vue 3 composables for FlexUI's interactive components. Each composable binds a
[flexui-runtime](https://www.npmjs.com/package/flexui-runtime) controller on
mount and disposes it on unmount, so dialog and menu state lives exactly as long
as the component that owns it.

```bash
npm install flexui-runtime flexui-vue
```

`vue` (3.5 or newer) is a peer dependency.

## Dialog

```vue
<script setup>
import { useDialog } from "flexui-vue";

const { ref, open, close, isOpen } = useDialog({ closeOnBackdrop: true });
</script>

<template>
  <button class="fui-btn" @click="open">Settings</button>
  <dialog ref="ref" class="fui-dialog" data-fui-dialog>
    <h2>Settings</h2>
    <button class="fui-btn" @click="close">Done</button>
  </dialog>
</template>
```

`useDialog(options)` returns `ref`, `open()`, `close()`, `toggle()` and `isOpen`.

## Dropdown

```vue
<script setup>
import { ref as vueRef } from "vue";
import { useDropdown } from "flexui-vue";

const menu = vueRef(null);
const trigger = useDropdown(menu);
</script>

<template>
  <button :ref="trigger.ref" class="fui-btn">Account</button>
  <div ref="menu" role="menu">…</div>
</template>
```

`useDropdown(menu, options)` takes the menu element (a template ref) or a
selector and returns the trigger `ref` with the open state.

## Popover

```vue
const trigger = usePopover(panel);
```

## Tabs

```vue
<script setup>
import { useTabs } from "flexui-vue";

const { ref, activeIndex, select } = useTabs();
</script>

<template>
  <div ref="ref" class="fui-tabs" data-fui-tabs>
    <div role="tablist">…</div>
  </div>
</template>
```

## Theme

```vue
<script setup>
import { useTheme } from "flexui-vue";

const { mode, resolved, setMode } = useTheme();
</script>
```

`setMode("light" | "dark" | "system")` switches the theme without a rebuild, and
`resolved` tells you which theme is actually active.

## Documentation

See the [FlexUI repository](https://github.com/Ahaduzzamankhan/flexui#readme) and
[docs/runtime.md](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/runtime.md)
for markup contracts, options and events.

## License

MIT
