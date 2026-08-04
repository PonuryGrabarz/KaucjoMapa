import {
	SOURCE_ID,
	POINT_LAYER,
	CLUSTER_LAYER,
	POINT_HOVER_LAYER,
	POINT_ID_PROPERTY,
	POINT_SELECTED_LAYER,
} from "./constants";

import type {
	LngLat,
	MapMouseEvent,
	GeoJSONSource,
	MapGeoJSONFeature,
	Map as MapLibreMap,
	MapLayerMouseEvent,
	Point as ScreenPoint,
} from "maplibre-gl";

const ICON_CENTER_OFFSET_PX = 32;

export interface MapInteractionHandlers {
	isPicking: () => boolean;
	reducedMotion: () => boolean;
	onPointSelect: (id: string) => void;
	onPointHover: (id: string | undefined, lngLat: LngLat) => void;
	onLocationPick: (lngLat: LngLat) => void;
	onViewportChange: () => void;
}

const featurePointId = (feature: MapGeoJSONFeature | undefined): string | undefined => {
	const id: unknown = feature?.properties[POINT_ID_PROPERTY];

	return typeof id === "string" ? id : undefined;
}

const closestPointFeature = (
	map: MapLibreMap,
	screenPoint: ScreenPoint,
): MapGeoJSONFeature | undefined => {
	const features = map.queryRenderedFeatures(screenPoint, {
		layers: [POINT_HOVER_LAYER, POINT_SELECTED_LAYER, POINT_LAYER],
	});

	const seen = new Set<string>();

	let closest: MapGeoJSONFeature | undefined;
	let closestDistanceSquared = Number.POSITIVE_INFINITY;

	for (const feature of features) {
		if (feature.geometry.type !== "Point") continue;

		const id = featurePointId(feature);

		if (id === undefined || seen.has(id)) continue;

		seen.add(id);

		const anchor = map.project(feature.geometry.coordinates as [number, number]);

		const horizontal = screenPoint.x - anchor.x;
		const vertical = screenPoint.y - (anchor.y - ICON_CENTER_OFFSET_PX);

		const distanceSquared = horizontal ** 2 + vertical ** 2;

		if (distanceSquared >= closestDistanceSquared) continue;

		closestDistanceSquared = distanceSquared;
		closest = feature;
	}

	return closest;
}

export const bindMapInteractions = (
	map: MapLibreMap,
	handlers: MapInteractionHandlers,
): (() => void) => {
	const canvas = map.getCanvas();

	const setCursor = (cursor: string) => {
		canvas.style.cursor = handlers.isPicking() ? "crosshair" : cursor;
	};

	const handleClusterClick = (event: MapLayerMouseEvent) => {
		if (handlers.isPicking()) return;

		const feature = event.features?.[0];
		const clusterId = Number(feature?.properties["cluster_id"]);

		if (!feature || feature.geometry.type !== "Point" || !Number.isFinite(clusterId)) return;

		const center = feature.geometry.coordinates as [number, number];

		const source = map.getSource(SOURCE_ID) as GeoJSONSource | undefined;

		if (!source) return;

		void source.getClusterExpansionZoom(clusterId).then((zoom) => {
			map.easeTo({ center, zoom, duration: handlers.reducedMotion() ? 0 : 480 });
		});
	};

	const handleClick = (event: MapMouseEvent) => {
		if (handlers.isPicking()) {
			handlers.onLocationPick(event.lngLat);

			return;
		}

		const id = featurePointId(closestPointFeature(map, event.point));

		if (id !== undefined) handlers.onPointSelect(id);
	};

	const handleMouseMove = (event: MapMouseEvent) => {
		if (handlers.isPicking()) {
			handlers.onPointHover(undefined, event.lngLat);
			setCursor("");

			return;
		}

		const id = featurePointId(closestPointFeature(map, event.point));

		handlers.onPointHover(id, event.lngLat);
		setCursor(id === undefined ? "" : "pointer");
	};

	const handleClusterEnter = () => setCursor("zoom-in");
	const handleClusterLeave = () => setCursor("");

	const handleCanvasLeave = () => {
		handlers.onPointHover(undefined, map.getCenter());
		setCursor("");
	};

	const handleMoveEnd = () => handlers.onViewportChange();

	map.on("click", CLUSTER_LAYER, handleClusterClick);
	map.on("click", handleClick);
	map.on("mousemove", handleMouseMove);
	map.on("mouseenter", CLUSTER_LAYER, handleClusterEnter);
	map.on("mouseleave", CLUSTER_LAYER, handleClusterLeave);
	map.on("moveend", handleMoveEnd);

	canvas.addEventListener("mouseleave", handleCanvasLeave);

	return () => {
		map.off("click", CLUSTER_LAYER, handleClusterClick);
		map.off("click", handleClick);
		map.off("mousemove", handleMouseMove);
		map.off("mouseenter", CLUSTER_LAYER, handleClusterEnter);
		map.off("mouseleave", CLUSTER_LAYER, handleClusterLeave);
		map.off("moveend", handleMoveEnd);

		canvas.removeEventListener("mouseleave", handleCanvasLeave);
	};
}
