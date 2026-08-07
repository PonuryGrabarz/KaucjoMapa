<template src="./distance-range.svelte.html"></template>

<style lang="sass" src="./distance-range.sass"></style>

<script lang="ts">
	import { geolocationAccess } from "@shared/helpers/geolocation-access.svelte";

	export type DistanceRangeProps = {
		legend?: string;
		maxKm?: number;
		min?: number;
		max?: number;
		minName?: string;
		maxName?: string;
		latName?: string;
		lonName?: string;
		emitPosition?: boolean;
	};

	let {
		legend = "Odległość od lokalizacji",
		maxKm = 100,
		min = $bindable(0),
		max = $bindable(maxKm),
		minName = "min_distance",
		maxName = "max_distance",
		latName = "lat",
		lonName = "lon",
		emitPosition = true,
	}: DistanceRangeProps = $props();

	const position = $derived(geolocationAccess.position);

	const enabled = $derived(geolocationAccess.granted && position !== undefined);

	const percent = (value: number) => `${(value / maxKm) * 100}%`;

	const unavailableMessage = $derived(
		geolocationAccess.permission === "denied"
			? "Przeglądarka odmówiła dostępu do lokalizacji, więc filtr odległości jest wyłączony."
			: geolocationAccess.permission === "unsupported"
				? "Ta przeglądarka nie udostępnia lokalizacji, więc filtr odległości jest wyłączony."
				: "Filtr odległości włączy się, gdy zezwolisz przeglądarce na dostęp do lokalizacji.",
	);

	const setMin = (value: number) => {
		min = Math.min(value, max);
	}

	const setMax = (value: number) => {
		max = Math.max(value, min);
	}
</script>
