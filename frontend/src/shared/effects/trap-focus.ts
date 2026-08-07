import type { Action } from "svelte/action";

export const FOCUSABLE_SELECTOR =
	'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])';

const isVisible = (element: HTMLElement): boolean =>
	element.checkVisibility?.() ?? element.offsetParent !== null;

const focusableElements = (node: HTMLElement): HTMLElement[] =>
	Array.from(node.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR)).filter(isVisible);

export const trapFocus: Action<HTMLElement, boolean | undefined> = (node, enabled = true) => {
	let isEnabled = enabled ?? true;

	const handleKeydown = (event: KeyboardEvent) => {
		if (!isEnabled || event.key !== "Tab") return;

		const elements = focusableElements(node);

		if (elements.length === 0) return;

		const first = elements[0]!;
		const last = elements.at(-1)!;

		const active = document.activeElement;

		if (event.shiftKey && active === first) {
			event.preventDefault();
			last.focus();
		} else if (!event.shiftKey && active === last) {
			event.preventDefault();
			first.focus();
		}
	}

	node.addEventListener("keydown", handleKeydown);

	return {
		update(next) {
			isEnabled = next ?? true;
		},
		destroy() {
			node.removeEventListener("keydown", handleKeydown);
		},
	};
};
