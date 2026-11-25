import { type Cleanup } from "./types.js";
export interface DialogOptions {
    /** Close when the backdrop is clicked. Defaults to true. */
    closeOnBackdrop?: boolean;
    /** Return focus to the opener when the dialog closes. Defaults to true. */
    restoreFocus?: boolean;
    onOpen?: (dialog: HTMLDialogElement) => void;
    onClose?: (dialog: HTMLDialogElement, reason: string) => void;
}
/**
 * Accessible controller for a native <dialog>.
 *
 * showModal() supplies the top layer and inert background. This class adds
 * backdrop dismissal, scroll locking, focus restoration and DOM events.
 */
export declare class Dialog {
    readonly element: HTMLDialogElement;
    private readonly options;
    private readonly disposer;
    private lastFocus;
    private openState;
    private pendingReason;
    private releaseTrap;
    constructor(element: HTMLDialogElement, options?: DialogOptions);
    get isOpen(): boolean;
    open(): void;
    close(reason?: string): void;
    dispose(): void;
    private adoptOpen;
    private markOpen;
    private handleClosed;
    private handleClick;
}
/** Returns an existing controller, if one was bound. */
export declare function getDialog(element: HTMLDialogElement): Dialog | undefined;
/** Creates or reuses the controller for a dialog element. */
export declare function bindDialog(element: HTMLDialogElement, options?: DialogOptions): Dialog;
/**
 * Delegated open/close listeners for data-fui-open and data-fui-close.
 * This is exported so custom overlays can share the same trigger protocol.
 */
export declare function bindOverlayTriggers(): Cleanup;
/** Binds every dialog/drawer under root and the shared trigger delegation. */
export declare function initDialogs(root?: ParentNode): Cleanup;
