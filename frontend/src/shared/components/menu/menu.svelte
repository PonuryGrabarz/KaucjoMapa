<template src="./menu.svelte.html"></template>

<script lang="ts">
	import Popover from "@shared/components/popover/popover.svelte";

	import type { Snippet } from "svelte";
	import type { HTMLAttributes } from "svelte/elements";

	export type MenuProps = Omit<
		HTMLAttributes<HTMLDivElement>,
		"children" | "role" | "aria-label"
	> & {
		open?: boolean;
		"aria-label": string;
		align?: "start" | "end";
		width?: string;
		triggerElement?: HTMLElement | null;
		onclose?: () => void;
		children: Snippet;
	};

	let {
		open = $bindable(false),
		children,
		...rest
	}: MenuProps = $props();
</script>
