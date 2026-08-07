<template src="./popover.svelte.html"></template>

<style lang="sass" src="./popover.sass"></style>

<script lang="ts">
	import { autoFocus } from "@shared/effects/auto-focus";
	import { trapFocus } from "@shared/effects/trap-focus";
	import { clickOutside } from "@shared/effects/click-outside";

	import type { Snippet } from "svelte";
	import type { HTMLAttributes } from "svelte/elements";

	export type PopoverProps = Omit<
		HTMLAttributes<HTMLDivElement>,
		"children" | "role" | "aria-label"
	> & {
		open?: boolean;
		"aria-label": string;
		modal?: boolean;
		role?: "dialog";
		align?: "start" | "end";
		width?: string;
		padding?: string;
		triggerElement?: HTMLElement | null;
		onclose?: () => void;
		children: Snippet;
	};

	let {
		open = $bindable(false),
		modal = true,
		role,
		align = "end",
		width,
		padding,
		triggerElement = null,
		onclose,
		class: className,
		children,
		...rest
	}: PopoverProps = $props();

	const resolvedRole = $derived(role ?? (modal ? "dialog" : undefined));

	const close = () => {
		if (!open) return;

		open = false;

		onclose?.();
	};

	const handleWindowKeydown = (event: KeyboardEvent) => {
		if (!open || event.key !== "Escape") return;

		event.stopPropagation();

		close();
	};

	const handleFocusOut = (event: FocusEvent) => {
		if (modal) return;

		const next = event.relatedTarget;

		if (!(next instanceof Node)) return;

		if (event.currentTarget instanceof HTMLElement && event.currentTarget.contains(next)) return;

		if (triggerElement?.contains(next)) return;

		close();
	};
</script>
