/**
 * Zero-config entry point.
 *
 * Importing this module (for example `<script type="module" src=".../auto.js">`)
 * initialises the runtime once the document has been parsed, and exposes the
 * handle on `window.flexui` so inline scripts can drive it.
 */

import { init, type InitOptions } from "./index.js";
import type { Cleanup } from "./types.js";

let dispose: Cleanup | null = null;

/** Initialises the runtime once; repeated calls return the same cleanup. */
export function start(options: InitOptions = {}): Cleanup {
  if (dispose) {
    return dispose;
  }
  dispose = init(document, options);
  return dispose;
}

/** Tears the runtime down and forgets the handle. */
export function stop(): void {
  dispose?.();
  dispose = null;
}

if (typeof window !== "undefined" && typeof document !== "undefined") {
  const boot = (): void => {
    start();
  };
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot, { once: true });
  } else {
    boot();
  }

  (window as Window & { flexui?: unknown }).flexui = {
    start,
    stop,
    init,
  };
}

export { init };
export type { InitOptions };
