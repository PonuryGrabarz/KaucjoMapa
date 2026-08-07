<template src="./checkbox.svelte.html"></template>

<style lang="sass" src="./checkbox.sass"></style>

<script lang="ts">
	import type { HTMLInputAttributes } from "svelte/elements";

	export type CheckboxProps = Omit<HTMLInputAttributes, "type" | "checked" | "children"> & {
		label: string;
		checked?: boolean;
	};

	let { label, checked = $bindable(false), id, class: className, ...rest }: CheckboxProps = $props();

	const fallbackId = $props.id();

	const inputId = $derived(id ?? fallbackId);
</script>
