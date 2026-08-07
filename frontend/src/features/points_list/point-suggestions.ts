import { pointCity } from "@shared/helpers/point-city";

import type { ReturnPoint } from "@shared/api/responses";
import type { SearchSuggestionGroup } from "@shared/components/search_field/search-field-types";

const CITY_ICON = "material-symbols:location-city-rounded";
const POINT_ICON = "material-symbols:location-on-outline-rounded";

const byCountDescending = (a: [string, number], b: [string, number]): number => b[1] - a[1];

const rankedValues = (values: readonly string[]): string[] => {
	const counts = new Map<string, number>();

	for (const value of values) counts.set(value, (counts.get(value) ?? 0) + 1);

	return [...counts.entries()].sort(byCountDescending).map(([value]) => value);
}

export const buildSearchSuggestions = (
	points: readonly ReturnPoint[],
): SearchSuggestionGroup[] => [
	{
		keyPrefix: "city:",
		secondary: "Miasto",
		icon: CITY_ICON,
		values: rankedValues(points.map(pointCity).filter((city) => city.length > 1)),
	},
	{
		keyPrefix: "name:",
		secondary: "Nazwa punktu",
		icon: POINT_ICON,
		values: rankedValues(points.map((point) => point.name.trim()).filter(Boolean)),
	},
];
