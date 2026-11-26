/**
 * Zero-config entry point.
 *
 * Importing this module (for example `<script type="module" src=".../auto.js">`)
 * initialises the runtime once the document has been parsed, and exposes the
 * handle on `window.flexui` so inline scripts can drive it.
 */
import { init, type InitOptions } from "./index.js";
import type { Cleanup } from "./types.js";
/** Initialises the runtime once; repeated calls return the same cleanup. */
export declare function start(options?: InitOptions): Cleanup;
/** Tears the runtime down and forgets the handle. */
export declare function stop(): void;
export { init };
export type { InitOptions };
