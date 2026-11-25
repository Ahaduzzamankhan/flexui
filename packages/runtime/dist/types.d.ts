/** A function that removes listeners or reverses a side effect. */
export type Cleanup = () => void;
/** Collects cleanups and disposes them in reverse order. */
export interface Disposer {
    add(cleanup: Cleanup): void;
    dispose(): void;
    readonly size: number;
}
export declare function createDisposer(): Disposer;
/** Returns a process-local, DOM-safe unique id. */
export declare function uid(prefix?: string): string;
/** Adds an event listener and returns its remover. */
export declare function on(target: EventTarget, type: string, handler: EventListener, options?: AddEventListenerOptions | boolean): Cleanup;
export declare function qs<T extends Element = HTMLElement>(root: ParentNode, selector: string): T | null;
export declare function qsa<T extends Element = HTMLElement>(root: ParentNode, selector: string): T[];
export declare function isHTMLElement(value: unknown): value is HTMLElement;
/**
 * Resolves an element from an id, "#id" selector, or any CSS selector.
 * Returns null instead of throwing for malformed selectors.
 */
export declare function getElement(target: string | Element | null | undefined): HTMLElement | null;
/** Reads a data attribute or normal attribute, returning null when absent. */
export declare function attr(element: Element, name: string): string | null;
/** Dispatches a bubbling custom event and returns false when it was cancelled. */
export declare function emit(target: EventTarget, type: string, detail?: unknown): boolean;
/** Returns the closest ancestor matching a selector, starting at the event target. */
export declare function closestFromEvent(event: Event, selector: string): HTMLElement | null;
