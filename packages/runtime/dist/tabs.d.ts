import { type Cleanup } from "./types.js";
export interface TabChangeDetail {
    tabs: Tabs;
    tab: HTMLElement;
    panel: HTMLElement | null;
    index: number;
}
/** Accessible tabs with roving tabindex and automatic activation. */
export declare class Tabs {
    readonly root: HTMLElement;
    private readonly disposer;
    private readonly tabs;
    private readonly panels;
    private activeIndex;
    constructor(root: HTMLElement);
    get active(): {
        tab: HTMLElement;
        panel: HTMLElement | null;
        index: number;
    } | null;
    activate(index: number, focus?: boolean, notify?: boolean): void;
    dispose(): void;
    private handleKeydown;
}
/** Binds every [data-fui-tabs] container under root. */
export declare function initTabs(root?: ParentNode): Cleanup;
