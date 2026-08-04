import { POINT_ID_PROPERTY } from "./constants";
import { drawMarkerImage, markerImageKey, markerStyleOf } from "./helpers";

import type { ReturnPoint } from "@/shared/api/responses";
import type { FeatureCollection, Point } from "geojson";
import type { Map as MapLibreMap } from "maplibre-gl";

export const pointCoordinates = (point: ReturnPoint): [number, number] => [
	point.longitude,
	point.latitude,
];

export const ensureMarkerImages = (map: MapLibreMap, points: readonly ReturnPoint[]): void => {
	for (const point of points) {
		const style = markerStyleOf(point);
		const key = markerImageKey(style);

		if (map.hasImage(key)) continue;

		map.addImage(key, drawMarkerImage(style), { pixelRatio: 2 });
	}
}

export const toFeatureCollection = (points: readonly ReturnPoint[]): FeatureCollection<Point> => ({
	type: "FeatureCollection",
	features: points.map((point) => ({
		type: "Feature",
		id: point.id,
		geometry: { type: "Point", coordinates: pointCoordinates(point) },
		properties: {
			[POINT_ID_PROPERTY]: point.id,
			name: point.name,
			icon: markerImageKey(markerStyleOf(point)),
		},
	})),
});
