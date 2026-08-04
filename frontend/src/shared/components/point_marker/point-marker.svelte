<template src="./point-marker.svelte.html"></template>

<style lang="sass" src="./point-marker.sass"></style>

<script lang="ts">
	import {
		resolveMarkerText,
		resolveMarkerLabel,
		resolveMarkerBackground,
	} from "@shared/helpers/marker-appearance";

	import type { HTMLAttributes } from "svelte/elements";

	export type PointMarkerProps = Omit<HTMLAttributes<HTMLSpanElement>, "children"> & {
		shortName?: string | null;
		size?: number;
		pending?: boolean;
		pendingLabel?: string;
		label?: string;
		textColor?: string | null;
		backgroundColor?: string | null;
	};

	const {
		shortName,
		backgroundColor,
		textColor,
		size = 40,
		pending = false,
		pendingLabel = "Oczekuje na weryfikację",
		label,
		class: className,
		...rest
	}: PointMarkerProps = $props();

	const safeSize = $derived(Number.isFinite(size) ? Math.max(1, size) : 40);

	const safeBackground = $derived(resolveMarkerBackground(backgroundColor));
	const safeText = $derived(resolveMarkerText(textColor));

	const shortLabel = $derived(resolveMarkerLabel(shortName));

	const accessibleLabel = $derived.by(() => {
		const name = label?.trim();

		if (!name) return undefined;

		return pending ? `${name}, ${pendingLabel}` : name;
	});
</script>
