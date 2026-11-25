import { type Cleanup } from "./types.js";
/** Selector for elements that can receive keyboard focus. */
export declare const FOCUSABLE_SELECTOR: string;
/** Returns focusable descendants in DOM order, excluding hidden/aria-hidden nodes. */
export declare function focusable(container: ParentNode): HTMLElement[];
/**
 * Focuses the autofocus element, the first focusable descendant, or the
 * container itself (making it programmatically focusable when needed).
 */
export declare function focusFirst(container: HTMLElement): void;
/**
 * Keeps Tab and Shift+Tab inside a container. Native modal dialogs already do
 * this; the helper is for popovers, sheets and custom overlay implementations.
 */
export declare function trapFocus(container: HTMLElement): Cleanup;
