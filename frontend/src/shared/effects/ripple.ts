import type { Action } from "svelte/action";

export const ripple: Action<HTMLElement> = (node) => {
  const handlePointerDown = (event: PointerEvent) => {
		if (event.button !== 0) return;

		if (node.matches(":disabled, [aria-disabled='true']")) return;

		if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;

		const rect = node.getBoundingClientRect();

		const diameter = Math.max(rect.width, rect.height) * 1.4;

		const span = document.createElement("span");

		span.className = "ripple-ink";

		span.style.width = `${diameter}px`;
		span.style.height = `${diameter}px`;

		span.style.left = `${event.clientX - rect.left - diameter / 2}px`;
		span.style.top = `${event.clientY - rect.top - diameter / 2}px`;

		const previousPosition = getComputedStyle(node).position;

		if (previousPosition === "static") node.style.position = "relative";

		node.style.overflow = node.style.overflow || "hidden";

		node.append(span);

		span.addEventListener("animationend", () => span.remove(), { once: true });
	}

	node.addEventListener("pointerdown", handlePointerDown);
  
	return {
		destroy() {
			node.removeEventListener("pointerdown", handlePointerDown);
		},
	};
};
