/**
 * Zero-config entry point.
 *
 * Importing this module (for example `<script type="module" src=".../auto.js">`)
 * initialises the runtime once the document has been parsed, and exposes the
 * handle on `window.flexui` so inline scripts can drive it.
 */
import { init } from "./index.js";
let dispose = null;
/** Initialises the runtime once; repeated calls return the same cleanup. */
export function start(options = {}) {
    if (dispose) {
        return dispose;
    }
    dispose = init(document, options);
    return dispose;
}
/** Tears the runtime down and forgets the handle. */
export function stop() {
    dispose?.();
    dispose = null;
}
if (typeof window !== "undefined" && typeof document !== "undefined") {
    const boot = () => {
        start();
    };
    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", boot, { once: true });
    }
    else {
        boot();
    }
    window.flexui = {
        start,
        stop,
        init,
    };
}
export { init };
