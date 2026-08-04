<template src="./menu-item.svelte.html"></template>

<style lang="sass" src="./menu-item.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";

	import type { Snippet } from "svelte";
	import type { HTMLButtonAttributes } from "svelte/elements";

	export type MenuItemProps = Omit<HTMLButtonAttributes, "type"> & {
		href?: string;
		danger?: boolean;
		separatorBefore?: boolean;
		children: Snippet;
	};

	let {
		href,
		danger = false,
		separatorBefore = false,
		class: className,
		children,
		...rest
	}: MenuItemProps = $props();
</script>
