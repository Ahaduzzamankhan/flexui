export function createDisposer() {
    const cleanups = [];
    let disposed = false;
    return {
        add(cleanup) {
            if (disposed) {
                cleanup();
                return;
            }
            cleanups.push(cleanup);
        },
        dispose() {
            if (disposed) {
                return;
            }
            disposed = true;
            for (let index = cleanups.length - 1; index >= 0; index -= 1) {
                cleanups[index]?.();
            }
            cleanups.length = 0;
        },
        get size() {
            return cleanups.length;
        },
    };
}
let idCounter = 0;
/** Returns a process-local, DOM-safe unique id. */
export function uid(prefix = "fui") {
    idCounter += 1;
    return `${prefix}-${idCounter.toString(36)}`;
}
/** Adds an event listener and returns its remover. */
export function on(target, type, handler, options) {
    target.addEventListener(type, handler, options);
    return () => {
        target.removeEventListener(type, handler, options);
    };
}
export function qs(root, selector) {
    return root.querySelector(selector);
}
export function qsa(root, selector) {
    return Array.from(root.querySelectorAll(selector));
}
export function isHTMLElement(value) {
    return typeof HTMLElement !== "undefined" && value instanceof HTMLElement;
}
/**
 * Resolves an element from an id, "#id" selector, or any CSS selector.
 * Returns null instead of throwing for malformed selectors.
 */
export function getElement(target) {
    if (isHTMLElement(target)) {
        return target;
    }
    if (typeof target !== "string") {
        return null;
    }
    const value = target.trim();
    if (value === "") {
        return null;
    }
    const bySelector = (selector) => {
        try {
            return document.querySelector(selector);
        }
        catch {
            return null;
        }
    };
    if (value.startsWith("#") || value.startsWith(".") || value.startsWith("[")) {
        return bySelector(value);
    }
    return document.getElementById(value) ?? bySelector(value);
}
/** Reads a data attribute or normal attribute, returning null when absent. */
export function attr(element, name) {
    return element.getAttribute(name);
}
/** Dispatches a bubbling custom event and returns false when it was cancelled. */
export function emit(target, type, detail) {
    return target.dispatchEvent(new CustomEvent(type, {
        bubbles: true,
        cancelable: true,
        composed: true,
        detail,
    }));
}
/** Returns the closest ancestor matching a selector, starting at the event target. */
export function closestFromEvent(event, selector) {
    const target = event.target;
    if (!(target instanceof Element)) {
        return null;
    }
    return target.closest(selector);
}
