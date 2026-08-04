<template src="./search-field.svelte.html"></template>

<style lang="sass" src="./search-field.sass"></style>

<script lang="ts">
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import type { SearchSuggestion, SearchSuggestionGroup } from "./search-field-types";

	export type SearchFieldProps = {
		label: string;
		name?: string;
		value?: string;
		placeholder?: string;
		suggestions?: readonly SearchSuggestion[];
		source?: readonly SearchSuggestionGroup[];
		maxSuggestions?: number;
		minQueryLength?: number;
		submitOnEnter?: boolean;
		submitOnSelect?: boolean;
		invalid?: boolean;
		emptyMessage?: string;
		onselect?: (suggestion: SearchSuggestion) => void;
		onquerychange?: (query: string) => void;
	};

	let {
		label,
		name = "query",
		value = $bindable(""),
		placeholder,
		suggestions = [],
		source = [],
		maxSuggestions = 8,
		minQueryLength = 2,
		submitOnEnter = true,
		submitOnSelect = true,
		invalid = false,
		emptyMessage = "Nie znaleziono takiego miejsca.",
		onselect,
		onquerychange,
	}: SearchFieldProps = $props();

	const componentId = $props.id();

	const inputId = `${componentId}-input`;
	const listboxId = `${componentId}-listbox`;

	const optionId = (index: number) => `${componentId}-option-${index}`;

	let inputElement = $state<HTMLInputElement | undefined>(undefined);
	let focused = $state(false);
	let activeIndex = $state(-1);

	const sourceIndex = $derived(
		source.flatMap((group) =>
			group.values.map((label) => ({
				label,
				haystack: label.toLocaleLowerCase("pl"),
				key: `${group.keyPrefix}${label}`,
				secondary: group.secondary,
				icon: group.icon,
			})),
		),
	);

	const visibleSuggestions = $derived.by(() => {
		if (sourceIndex.length === 0) return suggestions;

		const needle = value.trim().toLocaleLowerCase("pl");

		if (needle.length < minQueryLength) return [];

		const matched: SearchSuggestion[] = [];

		for (const entry of sourceIndex) {
			if (!entry.haystack.includes(needle)) continue;

			matched.push({
				key: entry.key,
				label: entry.label,
				secondary: entry.secondary,
				icon: entry.icon,
			});

			if (matched.length === maxSuggestions) break;
		}

		return matched;
	});

	const isOpen = $derived(focused && value.trim().length >= minQueryLength);
	const hasSuggestions = $derived(visibleSuggestions.length > 0);

	const choose = (suggestion: SearchSuggestion, viaPointer: boolean) => {
		value = suggestion.label;
		focused = false;
		activeIndex = -1;

		onselect?.(suggestion);

		if (viaPointer && submitOnSelect) inputElement?.form?.requestSubmit();
	}

	const handleInput = () => {
		activeIndex = -1;

		onquerychange?.(value);
	}

	const handleKeydown = (event: KeyboardEvent) => {
		if (event.key === "Escape") {
			focused = false;
			activeIndex = -1;

			return;
		}

		if (event.key === "ArrowDown" || event.key === "ArrowUp") {
			if (!isOpen || !hasSuggestions) return;

			event.preventDefault();

			activeIndex =
				event.key === "ArrowDown"
					? (activeIndex + 1) % visibleSuggestions.length
					: activeIndex <= 0
						? visibleSuggestions.length - 1
						: activeIndex - 1;

			return;
		}

		if (event.key !== "Enter") return;

		const chosen = isOpen ? visibleSuggestions[activeIndex] : undefined;

		if (chosen) {
			choose(chosen, false);

			if (!submitOnEnter || !submitOnSelect) event.preventDefault();

			return;
		}

		if (!submitOnEnter) event.preventDefault();
	}
</script>
