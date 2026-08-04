<template src="./button.svelte.html"></template>

<style lang="sass" src="./button.sass"></style>

<script lang="ts">
	import { ripple } from "@shared/effects/ripple";

	import type { Snippet } from "svelte";
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from "svelte/elements";

	type Variant =
		| "primary"
		| "secondary"
		| "tonal"
		| "text"
		| "danger"
		| "menu"
		| "menu-danger"
		| "bottomnav";

	export type ButtonProps = Omit<HTMLButtonAttributes, "type"> & {
		variant?: Variant;
		type?: "button" | "submit" | "reset";
		href?: string;
		target?: HTMLAnchorAttributes["target"];
		rel?: HTMLAnchorAttributes["rel"];
		children: Snippet;
	};

	let {
		variant = "primary",
		type = "button",
		disabled = false,
		href,
		target,
		rel,
		form,
		class: className,
		"aria-pressed": ariaPressed,
		children,
		...rest
	}: ButtonProps = $props();

	const classes = $derived(["btn", `btn--${variant}`, className]);
</script>
