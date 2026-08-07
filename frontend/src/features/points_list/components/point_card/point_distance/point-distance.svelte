<template src="./point-distance.svelte.html"></template>

<script lang="ts">
	import { distanceKm } from "@shared/helpers/map";
	import { distanceLabel } from "@shared/helpers/format";
	import { geolocationAccess } from "@shared/helpers/geolocation-access.svelte";

	export type PointDistanceProps = {
		latitude: number;
		longitude: number;
		class?: string;
	};

	const { latitude, longitude, class: className }: PointDistanceProps = $props();

	const label = $derived.by(() => {
		const position = geolocationAccess.position;

		if (!position) return "";

		return distanceLabel(distanceKm([position.lon, position.lat], [longitude, latitude]));
	});
</script>
