<template src="./dialog.svelte.html"></template>

<style lang="sass" src="./dialog.sass"></style>

<script lang="ts">
	import IconButton from "@shared/components/icon_button/icon-button.svelte";

	import type { Snippet } from "svelte";
	import type { HTMLDialogAttributes } from "svelte/elements";

	export type DialogProps = Omit<
		HTMLDialogAttributes,
		"children" | "open" | "aria-label" | "aria-labelledby" | "onclose" | "oncancel"
	> & {
		open?: boolean;
		title: string;
		variant?: "modal" | "sheet";
		maxWidth?: number;
		closeLabel?: string;
		dismissOnBackdrop?: boolean;
		onclose?: () => void;
		children: Snippet;
		footer?: Snippet;
	};

	let {
		open = $bindable(false),
		title,
		variant = "modal",
		maxWidth = 500,
		closeLabel = "Zamknij",
		dismissOnBackdrop = true,
		onclose,
		children,
		footer,
		class: className,
		...rest
	}: DialogProps = $props();

	const titleId = $props.id();

	let dialogElement: HTMLDialogElement | undefined = $state();

	const classes = $derived(["dialog", `dialog--${variant}`, className]);

	$effect(() => {
		if (!dialogElement) return;

		if (open && !dialogElement.open) dialogElement.showModal();
		else if (!open && dialogElement.open) dialogElement.close();
	});

	const requestClose = () => dialogElement?.close();

	const handleClose = () => {
		open = false;

		onclose?.();
	}

	const handleClick = (event: MouseEvent) => {
		if (dismissOnBackdrop && event.target === dialogElement) requestClose();
	}
</script>
