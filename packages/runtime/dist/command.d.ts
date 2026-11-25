import { type Cleanup } from "./types.js";
export interface CommandPaletteOptions {
    /** Close the palette after a command is chosen. Defaults to true. */
    closeOnSelect?: boolean;
}
/**
 * Keyboard-first command palette over existing listbox markup.
 *
 * Item attributes: data-fui-command-item, data-value and optional
 * data-keywords. Selection dispatches a bubbling fui:command event.
 */
export declare class CommandPalette {
    readonly element: HTMLElement;
    private readonly input;
    private readonly list;
    private readonly empty;
    private readonly dialog;
    private readonly options;
    private readonly disposer;
    private entries;
    private filtered;
    private activeIndex;
    constructor(element: HTMLElement, options?: CommandPaletteOptions);
    get isOpen(): boolean;
    /** Re-reads items from the DOM. Call after adding or removing commands. */
    refresh(): void;
    open(): void;
    close(): void;
    toggle(): void;
    dispose(): void;
    private reset;
    private filter;
    private setActive;
    private handleKeydown;
    private select;
}
/** Binds [data-fui-command] palettes, Ctrl/Cmd+K and command openers. */
export declare function initCommandPalettes(root?: ParentNode): Cleanup;
