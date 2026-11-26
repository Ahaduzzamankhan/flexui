import { type Cleanup } from "./types.js";
interface NativePopover extends HTMLElement {
    showPopover: () => void;
    hidePopover: () => void;
    togglePopover: (force?: boolean) => boolean;
}
export interface PopoverOptions {
    /** Pixels between trigger and panel. Defaults to 8. */
    offset?: number;
}
/**
 * Positions a native popover against its trigger and keeps aria-expanded in
 * sync. Browsers without the popover API get a hidden-attribute fallback with
 * the same keyboard and outside-click behaviour.
 */
export declare class PopoverController {
    readonly trigger: HTMLElement;
    readonly panel: NativePopover;
    private readonly options;
    private readonly disposer;
    private readonly native;
    private openState;
    constructor(trigger: HTMLElement, panel: HTMLElement, options?: PopoverOptions);
    get isOpen(): boolean;
    open(): void;
    close(): void;
    toggle(): void;
    dispose(): void;
    private syncOpen;
    private position;
}
/** Returns the controller bound to a trigger, if any. */
export declare function getPopover(trigger: HTMLElement): PopoverController | undefined;
/** Creates or reuses the controller for a trigger and its panel. */
export declare function bindPopover(trigger: HTMLElement, panel: HTMLElement, options?: PopoverOptions): PopoverController;
/** Binds popovers declared with data-fui-popover="<target>". */
export declare function initPopovers(root?: ParentNode): Cleanup;
export {};
