<template src="./filter-group.svelte.html"></template>

<style lang="sass" src="./filter-group.sass"></style>

<script lang="ts">
	import Checkbox from "@shared/components/checkbox/checkbox.svelte";

	import type { FilterOption } from "./filter-group-types";

	export type FilterGroupProps = {
		legend: string;
		name: string;
		options: readonly FilterOption[];
		checked?: readonly string[];
		onToggle?: (value: string, checked: boolean) => void;
		columns?: boolean;
		submitEmpty?: boolean;
	};

	const {
		legend,
		name,
		options,
		checked = [],
		onToggle,
		columns = false,
		submitEmpty = false,
	}: FilterGroupProps = $props();
</script>
