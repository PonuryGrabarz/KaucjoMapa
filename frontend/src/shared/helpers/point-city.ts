import type { ReturnPoint } from "@shared/api/responses";

const STREET_MARKER_PATTERN = /\s(?:ul|al|os|pl|ulica|aleja|aleje|osiedle|plac|rondo)\.?\s/i;

export const pointCity = (point: ReturnPoint): string => {
	const address = point.address?.trim();

	if (!address) return "";

	const marker = address.match(STREET_MARKER_PATTERN);

	if (marker?.index !== undefined) return address.slice(0, marker.index).trim();

	const firstDigit = address.search(/\d/);

	const head = (firstDigit > 0 ? address.slice(0, firstDigit) : address).trim().split(/\s+/);

	return (head.length > 1 ? head.slice(0, -1).join(" ") : head[0]) ?? "";
}
