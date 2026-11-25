import { type Cleanup } from "./types.js";
export type ToastVariant = "info" | "success" | "danger";
export interface ToastOptions {
    title: string;
    description?: string;
    variant?: ToastVariant;
    /** Milliseconds before dismissal. Use 0 to keep the toast open. */
    duration?: number;
    dismissible?: boolean;
}
export interface ToastHandle {
    readonly element: HTMLElement;
    dismiss(): void;
}
/** Creates and announces a toast, creating the live region on first use. */
export declare function toast(options: ToastOptions): ToastHandle;
/** Binds static toast dismiss buttons and auto-dismiss timers. */
export declare function initToasts(root?: ParentNode): Cleanup;
