import { type Cleanup } from "./types.js";
export type DrawerSide = "left" | "right";
/** Reads the side declared by data-fui-side, defaulting to right. */
export declare function drawerSide(element: HTMLElement): DrawerSide;
/** Binds every drawer under root, including its open/close triggers. */
export declare function initDrawers(root?: ParentNode): Cleanup;
