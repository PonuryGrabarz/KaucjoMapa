<template src="./status-pill.svelte.html"></template>

<style lang="sass" src="./status-pill.sass"></style>

<script module lang="ts">
	import type { PointStatus } from "@shared/helpers/return-point-labels";

	const TONE_CLASSES: Record<PointStatus, string> = {
		dostepny: "status-pill--available",
		nieczynny: "status-pill--off",
		brak_informacji: "status-pill--unknown",
	};
</script>

<script lang="ts">
	import { POINT_STATUS_LABELS } from "@shared/helpers/return-point-labels";

	import type { HTMLAttributes } from "svelte/elements";

	export type StatusPillProps = Omit<HTMLAttributes<HTMLSpanElement>, "children"> & {
		status: PointStatus;
		dot?: boolean;
	};

	const { status, dot = false, class: className, ...rest }: StatusPillProps = $props();

	const toneClass = $derived(TONE_CLASSES[status]);

	const label = $derived(POINT_STATUS_LABELS[status]);
</script>
