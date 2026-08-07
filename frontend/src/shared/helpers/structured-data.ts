import { hourRangeBounds, openingHoursDays, openingHoursEntries } from "@shared/helpers/opening-hours";

import type { ReturnPoint } from "@shared/api/responses";

const SCHEMA_DAYS: Record<string, string> = {
	mon: "Monday",
	tue: "Tuesday",
	wed: "Wednesday",
	thu: "Thursday",
	fri: "Friday",
	sat: "Saturday",
	sun: "Sunday",
};

const openingHoursSpecification = (point: ReturnPoint) =>
	openingHoursEntries(point.opening_hours).flatMap((entry) => {
		const [key, value] = entry;

		const bounds = hourRangeBounds(value);

		if (!bounds) return [];

		const days = openingHoursDays(key)
			.map((day) => SCHEMA_DAYS[day])
			.filter((day): day is string => day !== undefined);

		if (days.length === 0) return [];

		return [
			{
				"@type": "OpeningHoursSpecification",
				dayOfWeek: days,
				opens: bounds.opens,
				closes: bounds.closes,
			},
		];
	});

export const returnPointStructuredData = (point: ReturnPoint, url: string) => {
	const hours = openingHoursSpecification(point);

	return {
		"@context": "https://schema.org",
		"@type": "Place",
		"@id": url,
		name: point.name,
		url,
		...(point.description ? { description: point.description } : {}),
		...(point.address
			? { address: { "@type": "PostalAddress", streetAddress: point.address, addressCountry: "PL" } }
			: {}),
		geo: {
			"@type": "GeoCoordinates",
			latitude: point.latitude,
			longitude: point.longitude,
		},
		...(hours.length > 0 ? { openingHoursSpecification: hours } : {}),
	};
}
