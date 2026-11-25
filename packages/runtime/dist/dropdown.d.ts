import { type Cleanup } from "./types.js";
export interface DropdownOptions {
    /** Close after a menu item is selected. Defaults to true. */
    closeOnSelect?: boolean;
}
/**
 * A menu button controller: click/Enter/Space/Arrow keys, roving focus,
 * Escape, outside dismissal and aria-expanded wiring.
 */
export declare class Dropdown {
    readonly trigger: HTMLElement;
    readonly menu: HTMLElement;
    private readonly options;
    private readonly disposer;
    private openState;
    constructor(trigger: HTMLElement, menu: HTMLElement, options?: DropdownOptions);
    get isOpen(): boolean;
    open(focus?: "first" | "last" | "none"): void;
    close(restoreFocus?: boolean): void;
    toggle(): void;
    dispose(): void;
    private items;
    private handleTriggerKey;
    private handleMenuKey;
    private handleMenuClick;
}
/** Binds dropdowns declared with data-fui-dropdown="<target>". */
export declare function initDropdowns(root?: ParentNode): Cleanup;
