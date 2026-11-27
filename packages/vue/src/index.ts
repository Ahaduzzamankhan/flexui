/**
 * Vue bindings for `flexui-runtime`.
 *
 * Composables bind on mount and dispose on unmount, so Dialog/Dropdown state is
 * destroyed with the component that owns it instead of leaking listeners.
 * Template refs use `ref="..."` on the element and the returned `ref`.
 */

import { onBeforeUnmount, onMounted, ref, shallowRef, type Ref, type ShallowRef } from "vue";

import { bindDialog, type Dialog, type DialogOptions } from "flexui-runtime/dialog";
import { bindDropdown, type Dropdown, type DropdownOptions } from "flexui-runtime/dropdown";
import { bindPopover, type PopoverController, type PopoverOptions } from "flexui-runtime/popover";
import { bindTabs, type Tabs } from "flexui-runtime/tabs";
import {
  getResolvedTheme,
  getTheme,
  setTheme,
  toggleTheme,
  type ResolvedTheme,
  type ThemeChangeDetail,
  type ThemeMode,
} from "flexui-runtime/theme";

export interface UseDialogResult {
  /** `ref` for the `<dialog>` element. */
  element: Ref<HTMLDialogElement | null>;
  controller: ShallowRef<Dialog | null>;
  isOpen: Ref<boolean>;
  open(): void;
  close(reason?: string): void;
}

/** Binds a native `<dialog class="fui-dialog" data-fui-dialog>`. */
export function useDialog(options: DialogOptions = {}): UseDialogResult {
  const element = ref<HTMLDialogElement | null>(null);
  const controller = shallowRef<Dialog | null>(null);
  const isOpen = ref(false);

  onMounted(() => {
    const target = element.value;
    if (target === null) {
      return;
    }
    const instance = bindDialog(target, options);
    controller.value = instance;
    isOpen.value = instance.isOpen;
    instance.element.addEventListener("fui:open", handleOpen);
    instance.element.addEventListener("fui:close", handleClose);
  });

  onBeforeUnmount(() => {
    const instance = controller.value;
    if (instance === null) {
      return;
    }
    instance.element.removeEventListener("fui:open", handleOpen);
    instance.element.removeEventListener("fui:close", handleClose);
    instance.dispose();
    controller.value = null;
  });

  function handleOpen(): void {
    isOpen.value = true;
  }
  function handleClose(): void {
    isOpen.value = false;
  }

  return {
    element,
    controller,
    isOpen,
    open: () => controller.value?.open(),
    close: (reason?: string) => controller.value?.close(reason),
  };
}

export interface UseDropdownResult {
  /** `ref` for the trigger button. */
  trigger: Ref<HTMLElement | null>;
  controller: ShallowRef<Dropdown | null>;
  isOpen: Ref<boolean>;
  open(focus?: "first" | "last" | "none"): void;
  close(restoreFocus?: boolean): void;
  toggle(): void;
}

/** Binds a trigger to a menu element. */
export function useDropdown(menu: HTMLElement | null, options: DropdownOptions = {}): UseDropdownResult {
  const trigger = ref<HTMLElement | null>(null);
  const controller = shallowRef<Dropdown | null>(null);
  const isOpen = ref(false);

  const handleOpen = (): void => {
    isOpen.value = true;
  };
  const handleClose = (): void => {
    isOpen.value = false;
  };

  onMounted(() => {
    const source = trigger.value;
    if (source === null || menu === null) {
      return;
    }
    const instance = bindDropdown(source, menu, options);
    controller.value = instance;
    isOpen.value = instance.isOpen;
    source.addEventListener("fui:open", handleOpen);
    source.addEventListener("fui:close", handleClose);
  });

  onBeforeUnmount(() => {
    const instance = controller.value;
    if (instance === null) {
      return;
    }
    instance.trigger.removeEventListener("fui:open", handleOpen);
    instance.trigger.removeEventListener("fui:close", handleClose);
    instance.dispose();
    controller.value = null;
  });

  return {
    trigger,
    controller,
    isOpen,
    open: (focus?: "first" | "last" | "none") => controller.value?.open(focus),
    close: (restoreFocus?: boolean) => controller.value?.close(restoreFocus),
    toggle: () => controller.value?.toggle(),
  };
}

export interface UsePopoverResult {
  trigger: Ref<HTMLElement | null>;
  controller: ShallowRef<PopoverController | null>;
  isOpen: Ref<boolean>;
  open(): void;
  close(): void;
  toggle(): void;
}

/** Binds a trigger to a `popover` panel element. */
export function usePopover(panel: HTMLElement | null, options: PopoverOptions = {}): UsePopoverResult {
  const trigger = ref<HTMLElement | null>(null);
  const controller = shallowRef<PopoverController | null>(null);
  const isOpen = ref(false);

  const handleOpen = (): void => {
    isOpen.value = true;
  };
  const handleClose = (): void => {
    isOpen.value = false;
  };

  onMounted(() => {
    const source = trigger.value;
    if (source === null || panel === null) {
      return;
    }
    const instance = bindPopover(source, panel, options);
    controller.value = instance;
    isOpen.value = instance.isOpen;
    source.addEventListener("fui:popoveropen", handleOpen);
    source.addEventListener("fui:popoverclose", handleClose);
  });

  onBeforeUnmount(() => {
    const instance = controller.value;
    if (instance === null) {
      return;
    }
    instance.trigger.removeEventListener("fui:popoveropen", handleOpen);
    instance.trigger.removeEventListener("fui:popoverclose", handleClose);
    instance.dispose();
    controller.value = null;
  });

  return {
    trigger,
    controller,
    isOpen,
    open: () => controller.value?.open(),
    close: () => controller.value?.close(),
    toggle: () => controller.value?.toggle(),
  };
}

export interface UseTabsResult {
  /** `ref` for the `[data-fui-tabs]` container. */
  element: Ref<HTMLElement | null>;
  controller: ShallowRef<Tabs | null>;
  activeIndex: Ref<number>;
  activate(index: number, focus?: boolean): void;
}

/** Binds an ARIA tablist and tracks its active index. */
export function useTabs(): UseTabsResult {
  const element = ref<HTMLElement | null>(null);
  const controller = shallowRef<Tabs | null>(null);
  const activeIndex = ref(0);

  const handle = (event: Event): void => {
    activeIndex.value = (event as CustomEvent<{ index: number }>).detail.index;
  };

  onMounted(() => {
    const target = element.value;
    if (target === null) {
      return;
    }
    const instance = bindTabs(target);
    controller.value = instance;
    activeIndex.value = instance.active?.index ?? 0;
    target.addEventListener("fui:tabchange", handle);
  });

  onBeforeUnmount(() => {
    const target = element.value;
    if (target !== null) {
      target.removeEventListener("fui:tabchange", handle);
    }
    controller.value?.dispose();
    controller.value = null;
  });

  return {
    element,
    controller,
    activeIndex,
    activate: (index: number, focus?: boolean) => controller.value?.activate(index, focus),
  };
}

export interface UseThemeResult {
  mode: Ref<ThemeMode>;
  resolved: Ref<ResolvedTheme>;
  setMode(mode: ThemeMode): void;
  toggle(): void;
}

/** Reads and changes the FlexUI theme mode. */
export function useTheme(): UseThemeResult {
  const mode = ref<ThemeMode>(getTheme());
  const resolved = ref<ResolvedTheme>(getResolvedTheme());

  const handle = (event: Event): void => {
    const detail = (event as CustomEvent<ThemeChangeDetail>).detail;
    mode.value = detail.mode;
    resolved.value = detail.resolved;
  };

  onMounted(() => {
    document.addEventListener("fui:themechange", handle);
  });
  onBeforeUnmount(() => {
    document.removeEventListener("fui:themechange", handle);
  });

  return {
    mode,
    resolved,
    setMode: (next: ThemeMode) => setTheme(next),
    toggle: () => toggleTheme(),
  };
}
