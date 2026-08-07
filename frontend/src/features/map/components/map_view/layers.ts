import {
	SOURCE_ID,
	POINT_LAYER,
	CLUSTER_LAYER,
	MAP_FONT_BOLD,
	MAP_PIN_MIN_ZOOM,
	MAP_FONT_REGULAR,
	POINT_HOVER_LAYER,
	POINT_ID_PROPERTY,
	POINT_LABEL_LAYER,
	CLUSTER_COUNT_LAYER,
	POINT_SELECTED_LAYER,
	USER_LOCATION_DOT_LAYER,
	USER_LOCATION_SOURCE_ID,
	POLISH_MAP_LABEL_EXPRESSION,
	USER_LOCATION_ACCURACY_LAYER,
} from "./constants";

import type { AddLayerObject, GeoJSONSource, Map as MapLibreMap } from "maplibre-gl";

type SymbolLayout = NonNullable<Extract<AddLayerObject, { type: "symbol" }>["layout"]>;

const iconLayout = (size: number): SymbolLayout => ({
	"icon-image": ["get", "icon"],
	"icon-size": size,
	"icon-anchor": "bottom",
	"icon-allow-overlap": true,
	"icon-ignore-placement": true,
});

export const localizeMapLabels = (map: MapLibreMap): void => {
	for (const layer of map.getStyle().layers) {
		if (layer.type !== "symbol") continue;

		const textField = (layer.layout as Record<string, unknown> | undefined)?.["text-field"];

		if (!textField || !JSON.stringify(textField).includes("name")) continue;

		map.setLayoutProperty(
			layer.id,
			"text-field",
			POLISH_MAP_LABEL_EXPRESSION as SymbolLayout["text-field"],
		);

		map.setLayoutProperty(layer.id, "text-font", MAP_FONT_REGULAR as SymbolLayout["text-font"]);
	}
}

export const installReturnPointLayers = (map: MapLibreMap): void => {
	map.addSource(SOURCE_ID, {
		type: "geojson",
		data: { type: "FeatureCollection", features: [] },
		cluster: true,
		clusterRadius: 38,
		clusterMaxZoom: MAP_PIN_MIN_ZOOM - 1,
	});

	map.addLayer({
		id: CLUSTER_LAYER,
		type: "circle",
		source: SOURCE_ID,
		maxzoom: MAP_PIN_MIN_ZOOM,
		filter: ["has", "point_count"],
		paint: {
			"circle-color": ["step", ["get", "point_count"], "#2f8f46", 100, "#19723d", 750, "#0e4328"],
			"circle-radius": ["step", ["get", "point_count"], 20, 100, 26, 750, 32],
			"circle-stroke-width": 2,
			"circle-stroke-color": "#ffffff",
		},
	});

	map.addLayer({
		id: CLUSTER_COUNT_LAYER,
		type: "symbol",
		source: SOURCE_ID,
		maxzoom: MAP_PIN_MIN_ZOOM,
		filter: ["has", "point_count"],
		layout: {
			"text-field": ["get", "point_count_abbreviated"],
			"text-font": MAP_FONT_BOLD,
			"text-size": 13,
			"text-allow-overlap": true,
		},
		paint: {
			"text-color": "#ffffff",
			"text-halo-color": "rgb(0 0 0 / 0.2)",
			"text-halo-width": 1,
		},
	});

	map.addLayer({
		id: POINT_LAYER,
		type: "symbol",
		source: SOURCE_ID,
		minzoom: MAP_PIN_MIN_ZOOM,
		filter: ["!", ["has", "point_count"]],
		layout: iconLayout(1),
	});

	map.addLayer({
		id: POINT_LABEL_LAYER,
		type: "symbol",
		source: SOURCE_ID,
		minzoom: MAP_PIN_MIN_ZOOM,
		filter: ["!", ["has", "point_count"]],
		layout: {
			"text-field": ["get", "name"],
			"text-font": MAP_FONT_REGULAR,
			"text-size": 13,
			"text-offset": [0, 1.25],
			"text-anchor": "top",
			"text-optional": true,
		},
		paint: {
			"text-color": "#0e2a1b",
			"text-halo-color": "#ffffff",
			"text-halo-width": 2,
			"text-opacity": [
				"case",
				[
					"any",
					["boolean", ["feature-state", "selected"], false],
					["boolean", ["feature-state", "hover"], false],
				],
				1,
				0,
			],
		},
	});

	map.addLayer({
		id: POINT_SELECTED_LAYER,
		type: "symbol",
		source: SOURCE_ID,
		minzoom: MAP_PIN_MIN_ZOOM,
		filter: ["==", ["get", POINT_ID_PROPERTY], ""],
		layout: iconLayout(1.14),
	});

	map.addLayer({
		id: POINT_HOVER_LAYER,
		type: "symbol",
		source: SOURCE_ID,
		minzoom: MAP_PIN_MIN_ZOOM,
		filter: ["==", ["get", POINT_ID_PROPERTY], ""],
		layout: iconLayout(1.08),
	});
}

export const setPointFeatureFlag = (
	map: MapLibreMap,
	id: string,
	flag: "selected" | "hover",
	value: boolean,
): void => {
	map.setFeatureState({ source: SOURCE_ID, id }, { [flag]: value });
}

export const applyHighlightFilters = (
	map: MapLibreMap,
	selectedId: string | undefined,
	hoveredId: string | undefined,
): void => {
	map.setFilter(POINT_SELECTED_LAYER, ["==", ["get", POINT_ID_PROPERTY], selectedId ?? ""]);
	map.setFilter(POINT_HOVER_LAYER, ["==", ["get", POINT_ID_PROPERTY], hoveredId ?? ""]);

	const excluded = [selectedId, hoveredId].filter((id) => id !== undefined);

	map.setFilter(
		POINT_LAYER,
		excluded.length
			? [
					"all",
					["!", ["has", "point_count"]],
					["!", ["in", ["get", POINT_ID_PROPERTY], ["literal", excluded]]],
				]
			: ["!", ["has", "point_count"]],
	);
}

export const showUserLocation = (
	map: MapLibreMap,
	center: [number, number],
	accuracy: number,
): void => {
	const data = {
		type: "Feature" as const,
		properties: { accuracy },
		geometry: { type: "Point" as const, coordinates: center },
	};

	const existing = map.getSource(USER_LOCATION_SOURCE_ID) as GeoJSONSource | undefined;

	if (existing) {
		existing.setData(data);

		return;
	}

	map.addSource(USER_LOCATION_SOURCE_ID, { type: "geojson", data });

	map.addLayer({
		id: USER_LOCATION_ACCURACY_LAYER,
		type: "circle",
		source: USER_LOCATION_SOURCE_ID,
		paint: {
			"circle-radius": 16,
			"circle-color": "#1a73e8",
			"circle-opacity": 0.2,
			"circle-stroke-color": "#1a73e8",
			"circle-stroke-opacity": 0.35,
			"circle-stroke-width": 1,
		},
	});

	map.addLayer({
		id: USER_LOCATION_DOT_LAYER,
		type: "circle",
		source: USER_LOCATION_SOURCE_ID,
		paint: {
			"circle-radius": 7,
			"circle-color": "#1a73e8",
			"circle-stroke-width": 3,
			"circle-stroke-color": "#ffffff",
		},
	});
}
