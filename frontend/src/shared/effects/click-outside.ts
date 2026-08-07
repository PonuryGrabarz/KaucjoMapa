import type { Action } from "svelte/action";

export type ClickOutsideOptions = {
	onOutside: () => void;
	ignore?: HTMLElement | null;
};

export const clickOutside: Action<HTMLElement, ClickOutsideOptions> = (node, options) => {
	let current = options;

	const handlePointerDown = (event: PointerEvent) => {
		if (!(event.target instanceof Node)) return;

		if (node.contains(event.target)) return;

		if (current.ignore?.contains(event.target)) return;

		current.onOutside();
	};

	document.addEventListener("pointerdown", handlePointerDown, { capture: true });

	return {
		update(next) {
			current = next;
		},
		destroy() {
			document.removeEventListener("pointerdown", handlePointerDown, { capture: true });
		},
	};
};
