/** A function that removes listeners or reverses a side effect. */
export type Cleanup = () => void;

/** Collects cleanups and disposes them in reverse order. */
export interface Disposer {
  add(cleanup: Cleanup): void;
  dispose(): void;
  readonly size: number;
}

export function createDisposer(): Disposer {
  const cleanups: Cleanup[] = [];
  let disposed = false;

  return {
    add(cleanup: Cleanup): void {
      if (disposed) {
        cleanup();
        return;
      }
      cleanups.push(cleanup);
    },
    dispose(): void {
      if (disposed) {
        return;
      }
      disposed = true;
      for (let index = cleanups.length - 1; index >= 0; index -= 1) {
        cleanups[index]?.();
      }
      cleanups.length = 0;
    },
    get size(): number {
      return cleanups.length;
    },
  };
}

let idCounter = 0;

/** Returns a process-local, DOM-safe unique id. */
export function uid(prefix = "fui"): string {
  idCounter += 1;
  return `${prefix}-${idCounter.toString(36)}`;
}

/** Adds an event listener and returns its remover. */
export function on(
  target: EventTarget,
  type: string,
  handler: EventListener,
  options?: AddEventListenerOptions | boolean,
): Cleanup {
  target.addEventListener(type, handler, options);
  return () => {
    target.removeEventListener(type, handler, options);
  };
}

export function qs<T extends Element = HTMLElement>(root: ParentNode, selector: string): T | null {
  return root.querySelector<T>(selector);
}

export function qsa<T extends Element = HTMLElement>(root: ParentNode, selector: string): T[] {
  return Array.from(root.querySelectorAll<T>(selector));
}

export function isHTMLElement(value: unknown): value is HTMLElement {
  return typeof HTMLElement !== "undefined" && value instanceof HTMLElement;
}

/**
 * Resolves an element from an id, "#id" selector, or any CSS selector.
 * Returns null instead of throwing for malformed selectors.
 */
export function getElement(target: string | Element | null | undefined): HTMLElement | null {
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

  const bySelector = (selector: string): HTMLElement | null => {
    try {
      return document.querySelector<HTMLElement>(selector);
    } catch {
      return null;
    }
  };

  if (value.startsWith("#") || value.startsWith(".") || value.startsWith("[")) {
    return bySelector(value);
  }
  return document.getElementById(value) ?? bySelector(value);
}

/** Reads a data attribute or normal attribute, returning null when absent. */
export function attr(element: Element, name: string): string | null {
  return element.getAttribute(name);
}

/** Dispatches a bubbling custom event and returns false when it was cancelled. */
export function emit(target: EventTarget, type: string, detail?: unknown): boolean {
  return target.dispatchEvent(
    new CustomEvent(type, {
      bubbles: true,
      cancelable: true,
      composed: true,
      detail,
    }),
  );
}

/** Returns the closest ancestor matching a selector, starting at the event target. */
export function closestFromEvent(event: Event, selector: string): HTMLElement | null {
  const target = event.target;
  if (!(target instanceof Element)) {
    return null;
  }
  return target.closest<HTMLElement>(selector);
}
