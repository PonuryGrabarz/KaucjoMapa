<template src="./text-field.svelte.html"></template>

<style lang="sass" src="./text-field.sass"></style>

<script lang="ts">
	import IconButton from "@shared/components/icon_button/icon-button.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import type { Snippet } from "svelte";
	import type { HTMLInputAttributes, HTMLTextareaAttributes } from "svelte/elements";

	export type TextFieldProps = Omit<
		HTMLInputAttributes & HTMLTextareaAttributes,
		"value" | "children"
	> & {
		label: string;
		value?: string;
		hint?: string;
		error?: string;
		icon?: string;
		labelEnd?: Snippet;
		labelHidden?: boolean;
		multiline?: boolean;
	};

	const REVEAL_ICON = "material-symbols:visibility-outline-rounded";
	const HIDE_ICON = "material-symbols:visibility-off-outline-rounded";

	let {
		label,
		value = $bindable(""),
		hint,
		error,
		icon,
		labelEnd,
		labelHidden = false,
		multiline = false,
		type = "text",
		id,
		class: className,
		...rest
	}: TextFieldProps = $props();

	let revealed = $state(false);

	const fallbackId = $props.id();

	const controlId = $derived(id ?? fallbackId);
	const hintId = $derived(`${controlId}-hint`);
	const errorId = $derived(`${controlId}-error`);

	const isPassword = $derived(type === "password");

	const describedBy = $derived(
		[error ? errorId : undefined, hint ? hintId : undefined].filter(Boolean).join(" ") || undefined,
	);

	const inputAttributes: HTMLInputAttributes = $derived(
		isPassword
			? {
					...rest,
					type: revealed ? "text" : "password",
					spellcheck: false,
					autocapitalize: "off",
				}
			: { ...rest, type },
	);
</script>
