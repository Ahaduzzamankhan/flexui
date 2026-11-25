import { type ThemeOptions } from "./theme.js";
import { type Cleanup } from "./types.js";
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
export declare function init(root?: ParentNode, options?: InitOptions): Cleanup;
