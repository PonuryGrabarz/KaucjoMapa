<template src="./map-view.svelte.html"></template>

<style lang="sass" src="./map-view.sass"></style>

<script lang="ts">
	import "maplibre-gl/dist/maplibre-gl.css";

	import Button from "@shared/components/button/button.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { untrack } from "svelte";
	import { supportsWebgl } from "./helpers";
	import { MediaQuery } from "svelte/reactivity";
	import { pluralPl } from "@shared/helpers/format";
	import { events } from "@shared/events/app-events";
	import { prefersReducedMotion } from "svelte/motion";
	import { bindMapInteractions } from "./interactions";
	import { mapState } from "../../stores/map-state.svelte";
	import { returnPoints } from "../../stores/return-points.svelte";
	import { ensureMarkerImages, pointCoordinates, toFeatureCollection } from "./geojson";
	import { applySearchParamsToMapState, writeMapStateToUrl } from "../../stores/map-state-url";

	import {
		RETURN_POINT_TYPE_LABELS,
		POINT_STATUS_LABELS,
	} from "@shared/helpers/return-point-labels";

	import type { Snippet } from "svelte";

	import {
		SOURCE_ID,
		POLAND_CENTER,
		MAP_STYLE_URL,
		MAP_LOAD_TIMEOUT_MS,
	} from "./constants";

	import {
		showUserLocation,
		localizeMapLabels,
		setPointFeatureFlag,
		applyHighlightFilters,
		installReturnPointLayers,
	} from "./layers";

	import type { ReturnPoint } from "@shared/api/responses";

	import type {
		LngLat,
		GeoJSONSource,
		Map as MapLibreMap,
		Popup as MapLibrePopup,
	} from "maplibre-gl";

	type MapLibreModule = typeof import("maplibre-gl");

	if (typeof window !== "undefined") applySearchParamsToMapState(window.location.search);

	const isDesktop = new MediaQuery("min-width: 768px", true);

	let mapAttempt = $state(0);
	let mapLoaded = $state(false);
	let isOffline = $state(false);
	let webglUnsupported = $state(false);
	let pickingLocation = $state(false);
	let visiblePointCount = $state(0);
	let mapError = $state<string | undefined>(undefined);
	let hoveredPointId = $state<string | undefined>(undefined);

	let map: MapLibreMap | undefined;
	let maplibre: MapLibreModule | undefined;
	let hoverPopup: MapLibrePopup | undefined;
	let focusedPointId: string | undefined;
	let appliedHoveredId: string | undefined;
	let appliedSelectedId: string | undefined;

	const filteredPoints = $derived(mapState.apply(returnPoints.all));

	const pointLabel = (count: number) => pluralPl(count, "punkt", "punkty", "punktów");

	const viewportCountLabel = $derived(
		mapLoaded
			? `${visiblePointCount} ${pointLabel(visiblePointCount)} w widocznym obszarze`
			: `${filteredPoints.length} ${pointLabel(filteredPoints.length)}`,
	);

	const filteredCountLabel = $derived(
		`Filtry dopuszczają ${filteredPoints.length} ${pointLabel(filteredPoints.length)}`,
	);

	const easeDuration = (duration: number) => (prefersReducedMotion.current ? 0 : duration);

	const updateVisibleCount = () => {
		if (!map || !mapLoaded) {
			visiblePointCount = filteredPoints.length;

			return;
		}

		const bounds = map.getBounds();

		visiblePointCount = filteredPoints.filter((point) =>
			bounds.contains(pointCoordinates(point)),
		).length;
	}

	const syncHighlight = (selected: string | undefined, hovered: string | undefined) => {
		if (!map || !mapLoaded) return;

		if (appliedSelectedId !== selected) {
			if (appliedSelectedId) setPointFeatureFlag(map, appliedSelectedId, "selected", false);

			if (selected) setPointFeatureFlag(map, selected, "selected", true);

			appliedSelectedId = selected;
		}

		if (appliedHoveredId !== hovered) {
			if (appliedHoveredId) setPointFeatureFlag(map, appliedHoveredId, "hover", false);

			if (hovered) setPointFeatureFlag(map, hovered, "hover", true);

			appliedHoveredId = hovered;
		}

		applyHighlightFilters(map, selected, hovered);
	}

	const removeHoverPopup = () => {
		hoverPopup?.remove();
	}

	const showHoverPopup = (point: ReturnPoint, lngLat: LngLat) => {
		if (!map || !maplibre) return;

		hoverPopup ??= new maplibre.Popup({
			offset: 30,
			closeButton: false,
			closeOnClick: false,
			className: "map-view__hover-popup",
		});

		const content = document.createDocumentFragment();
		const name = document.createElement("strong");
		const meta = document.createElement("span");

		name.textContent = point.name;
		meta.textContent = `${POINT_STATUS_LABELS[point.current_status]} · ${RETURN_POINT_TYPE_LABELS[point.type]}`;

		content.append(name, meta);

		hoverPopup.setLngLat(lngLat).setDOMContent(content).addTo(map);
	}

	const focusPoint = (id: string) => {
		const point = returnPoints.find(id);

		if (!point || !map) {
			mapState.selectionSettled = true;

			return;
		}

		const instance = map;

		const duration = easeDuration(420);

		const settle = () => {
			if (mapState.selectedReturnPointId === id) mapState.selectionSettled = true;
		};

		instance.once("moveend", settle);

		setTimeout(settle, duration + 120);

		instance.easeTo({
			center: pointCoordinates(point),
			zoom: Math.max(instance.getZoom(), 14),
			offset: [0, isDesktop.current ? 0 : -120],
			duration,
		});
	}

	const selectPoint = (id: string) => {
		hoveredPointId = undefined;

		removeHoverPopup();

		mapState.selectedReturnPointId = id;
	}

	const hoverPoint = (id: string | undefined, lngLat: LngLat) => {
		if (id === undefined || id === mapState.selectedReturnPointId) {
			hoveredPointId = undefined;

			removeHoverPopup();

			return;
		}

		hoveredPointId = id;

		const point = returnPoints.find(id);

		if (point) showHoverPopup(point, lngLat);
	}

	const pickLocation = (lngLat: LngLat) => {
		pickingLocation = false;

		void events.emit("map:location_picked", { longitude: lngLat.lng, latitude: lngLat.lat });
	}

	const requestLocation = () => {
		if (mapState.locationState === "requesting") return;

		if (!("geolocation" in navigator)) {
			mapState.locationState = "unavailable";

			return;
		}

		mapState.locationState = "requesting";

		navigator.geolocation.getCurrentPosition(
			(position) => {
				const center: [number, number] = [position.coords.longitude, position.coords.latitude];

				mapState.locationState = "active";
				mapState.referenceCenter = center;

				if (!map || !mapLoaded) return;

				showUserLocation(map, center, position.coords.accuracy);

				map.easeTo({ center, zoom: 13, duration: easeDuration(600) });
			},
			(error) => {
				mapState.locationState = error.code === error.PERMISSION_DENIED ? "denied" : "unavailable";
			},
			{ enableHighAccuracy: true, timeout: 9000, maximumAge: 120_000 },
		);
	}

	const retryMap = () => {
		mapError = undefined;
		mapAttempt += 1;
	}

	const initMap = (container: HTMLElement) => {
		isOffline = !navigator.onLine;

		if (!supportsWebgl()) {
			webglUnsupported = true;

			return;
		}

		let disposed = false;
		let unbindInteractions: (() => void) | undefined;
		let loadTimeout: ReturnType<typeof setTimeout> | undefined;

		const armLoadTimeout = () => {
			clearTimeout(loadTimeout);

			if (disposed || mapLoaded || mapError !== undefined || document.hidden) return;

			loadTimeout = setTimeout(() => {
				if (!mapLoaded)
					mapError = "Mapa ładuje się zbyt długo. Sprawdź połączenie i spróbuj ponownie.";
			}, MAP_LOAD_TIMEOUT_MS);
		};

		document.addEventListener("visibilitychange", armLoadTimeout);

		void Promise.all([import("maplibre-gl"), document.fonts.ready])
			.then(([module]) => {
				if (disposed) return;

				maplibre = module;

				const instance = new module.Map({
					container,
					style: MAP_STYLE_URL,
					center: POLAND_CENTER,
					zoom: 5.5,
					minZoom: 4,
					maxZoom: 18,
					attributionControl: false,
					cooperativeGestures: true,
					canvasContextAttributes: { preserveDrawingBuffer: true },
				});

				map = instance;

				instance.addControl(new module.AttributionControl({ compact: true }), "bottom-right");
				instance.addControl(new module.NavigationControl({ showCompass: false }), "bottom-right");

				armLoadTimeout();

				instance.once("load", () => {
					if (disposed) return;

					clearTimeout(loadTimeout);

					localizeMapLabels(instance);
					installReturnPointLayers(instance);

					unbindInteractions = bindMapInteractions(instance, {
						isPicking: () => pickingLocation,
						reducedMotion: () => prefersReducedMotion.current,
						onPointSelect: selectPoint,
						onPointHover: hoverPoint,
						onLocationPick: pickLocation,
						onViewportChange: updateVisibleCount,
					});

					mapLoaded = true;
				});
			})
			.catch(() => {
				if (!disposed) mapError = "Nie udało się uruchomić mapy. Spróbuj ponownie.";
			});

		return () => {
			disposed = true;

			clearTimeout(loadTimeout);

			document.removeEventListener("visibilitychange", armLoadTimeout);

			unbindInteractions?.();
			removeHoverPopup();

			map?.remove();

			map = undefined;
			maplibre = undefined;
			hoverPopup = undefined;
			mapLoaded = false;
			appliedHoveredId = undefined;
			appliedSelectedId = undefined;
		};
	}

	$effect(() => {
		const visible = filteredPoints;
		const ready = mapLoaded;

		if (!ready || !map) return;

		const source = map.getSource(SOURCE_ID) as GeoJSONSource | undefined;

		if (!source) return;

		ensureMarkerImages(map, visible);

		source.setData(toFeatureCollection(visible));

		appliedHoveredId = undefined;
		appliedSelectedId = undefined;

		untrack(() => {
			const selected = mapState.selectedReturnPointId;

			if (selected !== undefined && !visible.some((point) => point.id === selected)) {
				mapState.selectedReturnPointId = undefined;

				removeHoverPopup();
			}

			syncHighlight(mapState.selectedReturnPointId, hoveredPointId);
			updateVisibleCount();
		});
	});

	$effect(() => {
		const selected = mapState.selectedReturnPointId;
		const hovered = hoveredPointId;

		if (!mapLoaded) return;

		syncHighlight(selected, hovered);
	});

	$effect(() => {
		const id = mapState.selectedReturnPointId;
		const ready = mapLoaded;
		const unavailable = webglUnsupported || mapError !== undefined;

		if (id === undefined) {
			focusedPointId = undefined;

			return;
		}

		if (unavailable) {
			mapState.selectionSettled = true;

			return;
		}

		if (!ready || !map || id === focusedPointId) return;

		focusedPointId = id;

		untrack(() => focusPoint(id));
	});

	$effect(() => {
		writeMapStateToUrl();
	});

	$effect(() => {
		const unsubscribeCenter = events.on("map:center_on_location", ({ data }) => {
			if (!map) return;

			map.easeTo({
				center: [data.lng, data.lat],
				zoom: Math.max(map.getZoom(), 12),
				duration: easeDuration(600),
			});
		});

		const unsubscribeShowPoint = events.on("map:show_point_by_id", ({ data }) =>
			selectPoint(data.id),
		);

		const unsubscribeLocation = events.on("map:request_location", requestLocation);

		const unsubscribeBeginPick = events.on("map:begin_location_pick", () => {
			pickingLocation = true;
		});

		return () => {
			unsubscribeCenter();
			unsubscribeShowPoint();
			unsubscribeLocation();
			unsubscribeBeginPick();
		};
	});

</script>
