<template src="./icon-button.svelte.html"></template>

<style lang="sass" src="./icon-button.sass"></style>

<script lang="ts">
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { ripple } from "@shared/effects/ripple";
  
	import type { HTMLButtonAttributes } from "svelte/elements";

	export type IconButtonProps = Omit<HTMLButtonAttributes, "type" | "aria-label"> & {
		icon: string;
		"aria-label": string;
		type?: "button" | "submit" | "reset";
		iconSize?: number | string;
	};

	let {
		icon,
		iconSize = 24,
		type = "button",
		class: className,
		...rest
	}: IconButtonProps = $props();

	const classes = $derived(["icon-btn", className]);
</script>
