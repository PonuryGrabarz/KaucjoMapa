<template src="./search.svelte.html"></template>

<style lang="sass" src="./search.sass"></style>

<script lang="ts">
	import SearchField from "@shared/components/search_field/search-field.svelte";

	import { events } from "@shared/events/app-events";
	import { pointCity } from "@shared/helpers/point-city";
	import { POLAND_CENTER } from "../../map_view/constants";
	import { mapState } from "../../../stores/map-state.svelte";
	import { returnPoints } from "../../../stores/return-points.svelte";

	import type { ReturnPoint } from "@shared/api/responses";
	import type { SearchSuggestion } from "@shared/components/search_field/search-field-types";

	type SuggestionTarget =
		| { kind: "city"; center: [number, number] }
		| { kind: "point"; point: ReturnPoint };

	const MIN_QUERY_LENGTH = 2;
	const MAX_CITY_SUGGESTIONS = 3;
	const MAX_POINT_SUGGESTIONS = 5;

	const CITY_ICON = "material-symbols:location-city-rounded";
	const POINT_ICON = "material-symbols:location-on-outline-rounded";

	let searchQuery = $state("");

	const searchIndex = $derived.by(() => {
		const entries: { point: ReturnPoint; city: string; cityKey: string; haystack: string }[] = [];

		const totals: Record<string, { lon: number; lat: number; count: number }> = {};

		for (const point of returnPoints.all) {
			const city = pointCity(point);

			entries.push({
				point,
				city,
				cityKey: city.toLocaleLowerCase("pl"),
				haystack: `${point.name} ${point.address ?? ""}`.toLocaleLowerCase("pl"),
			});

			if (!city) continue;

			const total = (totals[city] ??= { lon: 0, lat: 0, count: 0 });

			total.lon += point.longitude;
			total.lat += point.latitude;
			total.count += 1;
		}

		const cityCenters = new Map(
			Object.entries(totals).map(
				([city, total]): [string, [number, number]] => [
					city,
					[total.lon / total.count, total.lat / total.count],
				],
			),
		);

		return { entries, cityCenters };
	});

	const matches = $derived.by(() => {
		const needle = searchQuery.trim().toLocaleLowerCase("pl");

		if (needle.length < MIN_QUERY_LENGTH) return { suggestions: [], targets: new Map() };

		const { entries, cityCenters } = searchIndex;

		const cities: string[] = [];
		const points: ReturnPoint[] = [];

		for (const entry of entries) {
			if (
				cities.length < MAX_CITY_SUGGESTIONS
				&& entry.city
				&& entry.cityKey.includes(needle)
				&& !cities.includes(entry.city)
			) {
				cities.push(entry.city);
			}

			if (points.length < MAX_POINT_SUGGESTIONS && entry.haystack.includes(needle))
				points.push(entry.point);

			if (cities.length === MAX_CITY_SUGGESTIONS && points.length === MAX_POINT_SUGGESTIONS) break;
		}

		const suggestions: SearchSuggestion[] = [
			...cities.map((city) => ({
				key: `city:${city}`,
				label: city,
				secondary: "Miasto",
				icon: CITY_ICON,
			})),

			...points.map((point) => ({
				key: `point:${point.id}`,
				label: point.name,
				secondary: point.address ?? "",
				icon: POINT_ICON,
			})),
		];

		const targets = new Map<string, SuggestionTarget>([
			...cities.map(
				(city): [string, SuggestionTarget] => [
					`city:${city}`,
					{ kind: "city", center: cityCenters.get(city) ?? POLAND_CENTER },
				],
			),

			...points.map(
				(point): [string, SuggestionTarget] => [
					`point:${point.id}`,
					{ kind: "point", point },
				],
			),
		]);

		return { suggestions, targets };
	});

	const suggestions = $derived(matches.suggestions);

	const handleSelect = (suggestion: SearchSuggestion) => {
		const target = matches.targets.get(suggestion.key);

		if (!target) return;

		mapState.searchNoResults = false;

		if (target.kind === "city") {
			void events.emit("map:center_on_location", {
				lat: target.center[1],
				lng: target.center[0],
			});

			return;
		}

		void events.emit("map:center_on_location", {
			lat: target.point.latitude,
			lng: target.point.longitude,
		});

		void events.emit("map:show_point_by_id", { id: target.point.id });
	}

	const handleQueryChange = () => {
		mapState.searchNoResults = false;
	}

	const handleSubmit = (event: SubmitEvent) => {
		event.preventDefault();

		const first = suggestions[0];

		if (!first) {
			mapState.searchNoResults = searchQuery.trim().length > 0;

			return;
		}

		handleSelect(first);

		searchQuery = first.label;
	}
</script>
