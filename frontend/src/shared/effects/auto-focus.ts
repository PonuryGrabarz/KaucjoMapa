import type { Action } from "svelte/action";

import { FOCUSABLE_SELECTOR } from "./trap-focus";

export const autoFocus: Action<HTMLElement> = (node) => {
	const restoreTo = document.activeElement instanceof HTMLElement ? document.activeElement : null;

	const firstFocusable = node.querySelector<HTMLElement>(FOCUSABLE_SELECTOR);

	(firstFocusable ?? node).focus();

	return {
		destroy() {
			restoreTo?.focus();
		},
	};
};
