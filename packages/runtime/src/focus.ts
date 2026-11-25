import { on, qsa, type Cleanup } from "./types.js";

/** Selector for elements that can receive keyboard focus. */
export const FOCUSABLE_SELECTOR = [
  "a[href]",
  "area[href]",
  "button:not([disabled])",
  "input:not([disabled]):not([type='hidden'])",
  "select:not([disabled])",
  "textarea:not([disabled])",
  "summary",
  "iframe",
  "object",
  "embed",
  "audio[controls]",
  "video[controls]",
  "[contenteditable='true']",
  "[tabindex]:not([tabindex='-1'])",
].join(",");

/** Returns focusable descendants in DOM order, excluding hidden/aria-hidden nodes. */
export function focusable(container: ParentNode): HTMLElement[] {
  return qsa<HTMLElement>(container, FOCUSABLE_SELECTOR).filter((element) => {
    if (element.getAttribute("aria-hidden") === "true") {
      return false;
    }
    if (element.closest("[hidden], [inert]")) {
      return false;
    }
    return element.getClientRects().length > 0 || element === document.activeElement;
  });
}

/**
 * Focuses the autofocus element, the first focusable descendant, or the
 * container itself (making it programmatically focusable when needed).
 */
export function focusFirst(container: HTMLElement): void {
  const autofocus = container.querySelector<HTMLElement>("[autofocus]");
  const first = autofocus ?? focusable(container)[0];
  if (first) {
    first.focus({ preventScroll: true });
    return;
  }
  if (!container.hasAttribute("tabindex")) {
    container.setAttribute("tabindex", "-1");
  }
  container.focus({ preventScroll: true });
}

/**
 * Keeps Tab and Shift+Tab inside a container. Native modal dialogs already do
 * this; the helper is for popovers, sheets and custom overlay implementations.
 */
export function trapFocus(container: HTMLElement): Cleanup {
  return on(container, "keydown", (event) => {
    const keyboard = event as KeyboardEvent;
    if (keyboard.key !== "Tab") {
      return;
    }

    const items = focusable(container);
    const first = items[0];
    const last = items[items.length - 1];
    if (!first || !last) {
      keyboard.preventDefault();
      focusFirst(container);
      return;
    }

    const active = document.activeElement;
    const outside = !(active instanceof Node) || !container.contains(active);

    if (keyboard.shiftKey && (outside || active === first)) {
      keyboard.preventDefault();
      last.focus();
      return;
    }
    if (!keyboard.shiftKey && (outside || active === last)) {
      keyboard.preventDefault();
      first.focus();
    }
  });
}
