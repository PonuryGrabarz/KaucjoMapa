<template src="./color-field.svelte.html"></template>

<style lang="sass" src="./color-field.sass"></style>

<script lang="ts">
	import { untrack } from "svelte";
	import { normalizeHexColor } from "@shared/helpers/color";

	import type { HTMLInputAttributes } from "svelte/elements";

	export type ColorFieldProps = Omit<HTMLInputAttributes, "type" | "value" | "children"> & {
		label: string;
		value?: string;
	};

	const DEFAULT_COLOR = "#1f7a3c";

	let {
		label,
		value = $bindable(DEFAULT_COLOR),
		id,
		disabled = false,
		class: className,
		...rest
	}: ColorFieldProps = $props();

	const fallbackId = $props.id();

	const inputId = $derived(id ?? fallbackId);
	const hexId = $derived(`${inputId}-hex`);

	let draft = $state(value);
	let invalid = $state(false);

	$effect(() => {
		const next = value;

		untrack(() => {
			if (normalizeHexColor(draft) === next) return;

			draft = next;
			invalid = false;
		});
	});

	const handleSwatchInput = (event: Event & { currentTarget: HTMLInputElement }) => {
		value = event.currentTarget.value;
	}

	const handleHexInput = (event: Event & { currentTarget: HTMLInputElement }) => {
		draft = event.currentTarget.value;

		const normalized = normalizeHexColor(draft);

		invalid = normalized === undefined;

		if (normalized !== undefined) value = normalized;
	}

	const handleHexBlur = () => {
		if (!invalid) return;

		draft = value;
		invalid = false;
	}
</script>
