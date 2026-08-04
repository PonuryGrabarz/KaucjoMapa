<template src="./filter-panel.svelte.html"></template>

<style lang="sass" src="./filter-panel.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";
	import Checkbox from "@shared/components/checkbox/checkbox.svelte";
	import IconButton from "@shared/components/icon_button/icon-button.svelte";
	import FilterGroup from "@shared/components/filter_group/filter-group.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";
	import DistanceRange from "@shared/components/distance_range/distance-range.svelte";

	import type { FilterPanelGroup } from "./filter-panel-types";

	export type FilterPanelProps = {
		groups: readonly FilterPanelGroup[];
		open?: boolean;
		title?: string;
		hint?: string;
		availabilityLegend?: string;
		openNow?: boolean;
		openNowName?: string;
		openNowLabel?: string;
		distanceMin?: number;
		distanceMax?: number;
		distanceMaxKm?: number;
		emitPosition?: boolean;
		activeCount?: number;
		keepMounted?: boolean;
		submitEmptyGroups?: boolean;
		footer?: boolean;
		resetLabel?: string;
		applyLabel?: string;
		onreset?: () => void;
	};

	let {
		groups,
		open = $bindable(false),
		title = "Filtry punktów",
		hint = "Zawęź wyniki bez utraty obszaru mapy.",
		availabilityLegend = "Dostępność",
		openNow = $bindable(false),
		openNowName,
		openNowLabel = "Otwarte teraz",
		distanceMin = $bindable(0),
		distanceMax = $bindable(100),
		distanceMaxKm = 100,
		emitPosition = true,
		activeCount = 0,
		keepMounted = false,
		submitEmptyGroups = false,
		footer = true,
		resetLabel = "Wyczyść filtry",
		applyLabel = "Zastosuj",
		onreset,
	}: FilterPanelProps = $props();

	const panelId = $props.id();

	const triggerLabel = $derived(
		activeCount > 0 ? `Filtry, aktywne grupy: ${activeCount}` : "Filtry",
	);

	const closePanel = () => {
		open = false;
	}

	const handleWindowKeydown = (event: KeyboardEvent) => {
		if (open && event.key === "Escape") closePanel();
	}
</script>
