<template src="./select-field.svelte.html"></template>

<style lang="sass" src="./select-field.sass"></style>

<script lang="ts" generics="TValue extends string">
	import type { HTMLSelectAttributes } from "svelte/elements";

	type SelectFieldOption<TOption extends string> = {
		value: TOption;
		label: string;
	};

	type SelectFieldProps = Omit<HTMLSelectAttributes, "value" | "children"> & {
		label: string;
		options: readonly SelectFieldOption<TValue>[];
		value?: TValue;
		hint?: string;
		labelHidden?: boolean;
	};

	let {
		label,
		options,
		value = $bindable(),
		hint,
		labelHidden = false,
		id,
		class: className,
		...rest
	}: SelectFieldProps = $props();

	const fallbackId = $props.id();

	const controlId = $derived(id ?? fallbackId);
	const hintId = $derived(`${controlId}-hint`);
</script>
