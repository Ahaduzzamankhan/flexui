import { initCommandPalettes } from "./command.js";
import { initDialogs } from "./dialog.js";
import { initDropdowns } from "./dropdown.js";
import { initPopovers } from "./popover.js";
import { initTabs } from "./tabs.js";
import { initTheme, type ThemeOptions } from "./theme.js";
import { initToasts } from "./toast.js";
import { createDisposer, type Cleanup } from "./types.js";

export * from "./command.js";
export * from "./dialog.js";
export * from "./drawer.js";
export * from "./dropdown.js";
export * from "./focus.js";
export * from "./popover.js";
export * from "./tabs.js";
export * from "./theme.js";
export * from "./toast.js";
export * from "./types.js";

export interface InitOptions {
  theme?: ThemeOptions;
}

/**
 * Initializes every runtime feature found under root.
 *
 * Markup is progressive-enhancement friendly: dialogs and drawers are native
 * <dialog> elements, popovers use the native popover API, and tabs are normal
 * ARIA tabs. Calling init() twice is safe; avoid it anyway in production.
 */
export function init(root: ParentNode = document, options: InitOptions = {}): Cleanup {
  const disposer = createDisposer();

  disposer.add(initTheme(options.theme ?? { root }));
  disposer.add(initDialogs(root));
  disposer.add(initDropdowns(root));
  disposer.add(initTabs(root));
  disposer.add(initPopovers(root));
  disposer.add(initCommandPalettes(root));
  disposer.add(initToasts(root));

  return () => disposer.dispose();
}
