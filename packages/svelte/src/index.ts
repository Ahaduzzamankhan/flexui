/**
 * Svelte actions for `@flexui/runtime`.
 *
 * ```svelte
 * <script>
 *   import { dialog, dropdown, tabs } from "@flexui/svelte";
 * </script>
 *
 * <dialog class="fui-dialog" use:dialog={{ closeOnBackdrop: true }}>…</dialog>
 * <button class="fui-btn" use:dropdown={"#account-menu"}>Account</button>
 * <div class="fui-tabs" data-fui-tabs use:tabs>…</div>
 * ```
 *
 * Actions need no Svelte imports, so the package has no build-time dependency
 * on the framework and works with Svelte 4 and 5 component syntax alike.
 * Parameters are read once, when the action runs.
 */

import { bindDialog, type DialogOptions } from "@flexui/runtime/dialog";
import { bindDropdown, type DropdownOptions } from "@flexui/runtime/dropdown";
import { bindPopover, type PopoverOptions } from "@flexui/runtime/popover";
import { bindTabs } from "@flexui/runtime/tabs";

/** The shape Svelte uses for an action's return value. */
export interface Action<Params = undefined> {
  update?(params: Params): void;
  destroy?(): void;
}

function noop(): void {
  // Actions without a parameter still need a destroy hook.
}

/** Binds a native `<dialog>`. */
export function dialog(node: HTMLDialogElement, options: DialogOptions = {}): Action<DialogOptions> {
  const controller = bindDialog(node, options);
  return { destroy: () => controller.dispose() };
}

/** Binds a trigger to a menu element or an id. */
export function dropdown(
  node: HTMLElement,
  menu: string | HTMLElement | null,
  options: DropdownOptions = {},
): Action<string | HTMLElement | null> {
  if (menu === null) {
    return { destroy: noop };
  }
  const target = typeof menu === "string" ? document.getElementById(menu.replace(/^#/, "")) : menu;
  if (target === null) {
    return { destroy: noop };
  }
  const controller = bindDropdown(node, target, options);
  return { destroy: () => controller.dispose() };
}

/** Binds a trigger to a `popover` panel element or an id. */
export function popover(
  node: HTMLElement,
  panel: string | HTMLElement | null,
  options: PopoverOptions = {},
): Action<string | HTMLElement | null> {
  if (panel === null) {
    return { destroy: noop };
  }
  const target = typeof panel === "string" ? document.getElementById(panel.replace(/^#/, "")) : panel;
  if (target === null) {
    return { destroy: noop };
  }
  const controller = bindPopover(node, target, options);
  return { destroy: () => controller.dispose() };
}

/** Binds an ARIA tablist container. */
export function tabs(node: HTMLElement): Action<undefined> {
  const controller = bindTabs(node);
  return { destroy: () => controller.dispose() };
}

export { bindDialog, bindDropdown, bindPopover, bindTabs };
export type { DialogOptions, DropdownOptions, PopoverOptions };
