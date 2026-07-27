import "maplibre-gl/dist/maplibre-gl.css";

import Styles from "./map.module.sass";

import Modal from "@components/modal/modal";
import IconButton from "@components/icon_button/icon-button";

import {
    createMapPinImageData,
    getMapPinImageId,
    MAP_PIN_IMAGE_PIXEL_RATIO,
} from "./map-pin";

import {
    isLocalizationEnabled,
    searchQuery,
    selectedDistanceKm,
    selectedPointStatuses,
    selectedPointTypes,
    userLocation,
    type PointStatus,
    type PointType,
    type UserLocation,
} from "../../store";

import { useCallback, useEffect, useMemo, useRef, useState } from "preact/hooks";

import type { GetMapPointsResponse } from "@api/responses";
import type { MapPoint } from "@api/responses/get-map-points-response";
import type {
    GeoJSONSource,
    Map as MapLibreMap,
    MapGeoJSONFeature,
    MapLayerMouseEvent,
} from "maplibre-gl";
import Button from "@/shared/components/button/button";

type MapProps = {
    initialReturnPoints: GetMapPointsResponse;
    selectedPoint?: string;
};

type MapStatus = "loading" | "ready" | "error";
type UserLocationMarkerStatus = "hidden" | "updating" | "visible" | "error";

type MapPointMetadata = {
    name: string;
    address: string;
    postalCode: string;
    pointType?: PointType;
    status?: PointStatus;
};

type DisplayMapPoint = MapPoint & MapPointMetadata;

type MapPointFeatureCollection = {
    type: "FeatureCollection";
    features: Array<{
        type: "Feature";
        id: string;
        properties: {
            returnPointId: string;
            iconId: string;
        };
        geometry: {
            type: "Point";
            coordinates: [number, number];
        };
    }>;
};

type UserLocationFeatureCollection = {
    type: "FeatureCollection";
    features: Array<{
        type: "Feature";
        properties: Record<string, never>;
        geometry: {
            type: "Point";
            coordinates: [number, number];
        };
    }>;
};

const MAP_STYLE_URL = "https://tiles.openfreemap.org/styles/positron";
const DEFAULT_CENTER: [number, number] = [19.1451, 51.9194];
const EARTH_RADIUS_KM = 6371;
const POINTS_SOURCE_ID = "return-points";
const CLUSTERS_LAYER_ID = "return-point-clusters";
const CLUSTER_COUNT_LAYER_ID = "return-point-cluster-count";
const MAP_PINS_LAYER_ID = "return-point-map-pins";
const HOVERED_MAP_PIN_LAYER_ID = "return-point-map-pin-hover";
const USER_LOCATION_SOURCE_ID = "user-location";
const USER_LOCATION_HALO_LAYER_ID = "user-location-halo";
const USER_LOCATION_POINT_LAYER_ID = "user-location-point";
const MAP_PIN_MIN_ZOOM = 9;
const POLISH_MAP_LABEL_EXPRESSION = [
    "coalesce",
    ["get", "name:pl"],
    ["get", "name"],
    ["get", "name_en"],
];

const normalizeSearchText = (value: string): string => {
    return value
        .normalize("NFD")
        .replace(/[\u0300-\u036f]/g, "")
        .toLocaleLowerCase("pl")
        .trim();
};

const isValidMapPoint = (point: MapPoint): boolean => {
    return Number.isFinite(point.latitude) && Number.isFinite(point.longitude);
};

const getUserLocationGeoJson = (
    location: UserLocation | null,
): UserLocationFeatureCollection => {
    if (
        !location
        || !Number.isFinite(location.latitude)
        || !Number.isFinite(location.longitude)
        || location.latitude < -90
        || location.latitude > 90
        || location.longitude < -180
        || location.longitude > 180
    ) {
        return {
            type: "FeatureCollection",
            features: [],
        };
    }

    return {
        type: "FeatureCollection",
        features: [{
            type: "Feature",
            properties: {},
            geometry: {
                type: "Point",
                coordinates: [location.longitude, location.latitude],
            },
        }],
    };
};

const getDistanceKm = (
    firstLatitude: number,
    firstLongitude: number,
    secondLatitude: number,
    secondLongitude: number,
): number => {
    const toRadians = (degrees: number): number => degrees * Math.PI / 180;
    const latitudeDelta = toRadians(secondLatitude - firstLatitude);
    const longitudeDelta = toRadians(secondLongitude - firstLongitude);
    const firstLatitudeRadians = toRadians(firstLatitude);
    const secondLatitudeRadians = toRadians(secondLatitude);

    const haversine = Math.sin(latitudeDelta / 2) ** 2
        + Math.cos(firstLatitudeRadians)
        * Math.cos(secondLatitudeRadians)
        * Math.sin(longitudeDelta / 2) ** 2;

    return EARTH_RADIUS_KM * 2 * Math.atan2(Math.sqrt(haversine), Math.sqrt(1 - haversine));
};

const getDisplayMapPoint = (point: MapPoint): DisplayMapPoint => {
    const pointWithOptionalMetadata = point as MapPoint & Partial<MapPointMetadata>;

    return {
        ...point,
        name: pointWithOptionalMetadata.name?.trim() || point.shortName,
        address: pointWithOptionalMetadata.address?.trim() || "",
        postalCode: pointWithOptionalMetadata.postalCode?.trim() || "",
        pointType: pointWithOptionalMetadata.pointType,
        status: pointWithOptionalMetadata.status,
    };
};

const getAnimationDuration = (): number => {
    return window.matchMedia("(prefers-reduced-motion: reduce)").matches ? 0 : 450;
};

const localizeMapLabels = (map: MapLibreMap): void => {
    map.getStyle().layers.forEach((layer) => {
        if (layer.type !== "symbol") return;

        const layout = layer.layout as Record<string, unknown> | undefined;
        const textField = layout?.["text-field"];
        if (!textField || !JSON.stringify(textField).includes("name")) return;

        map.setLayoutProperty(layer.id, "text-field", POLISH_MAP_LABEL_EXPRESSION);
    });
};

const fitMapToPoints = (map: MapLibreMap, points: readonly DisplayMapPoint[]): void => {
    if (points.length === 0) return;

    if (points.length === 1) {
        const [point] = points;
        map.easeTo({
            center: [point.longitude, point.latitude],
            zoom: Math.max(map.getZoom(), 14),
            duration: getAnimationDuration(),
        });
        return;
    }

    const longitudes = points.map((point) => point.longitude);
    const latitudes = points.map((point) => point.latitude);
    const bounds: [[number, number], [number, number]] = [
        [Math.min(...longitudes), Math.min(...latitudes)],
        [Math.max(...longitudes), Math.max(...latitudes)],
    ];

    map.fitBounds(bounds, {
        padding: 64,
        maxZoom: 14,
        duration: getAnimationDuration(),
    });
};

const centerMapOnPoint = (
    map: MapLibreMap,
    point: Pick<MapPoint, "latitude" | "longitude">,
): void => {
    map.easeTo({
        center: [point.longitude, point.latitude],
        zoom: Math.max(map.getZoom(), 20),
        duration: getAnimationDuration(),
    });
};

type DialogData = {
    isOpen: boolean;
    pointId: string;
    title: string;
    subtitle: string;
    logoBgColor: string;
    logoTextColor: string;
    shortName: string;
};

export default function Map({ initialReturnPoints, selectedPoint }: MapProps) {
    const mapContainerRef = useRef<HTMLDivElement>(null);
    const mapRef = useRef<MapLibreMap | null>(null);
    const sourceUpdateIdRef = useRef(0);
    const userLocationUpdateIdRef = useRef(0);
    const [mapStatus, setMapStatus] = useState<MapStatus>("loading");
    const [userLocationMarkerStatus, setUserLocationMarkerStatus] =
        useState<UserLocationMarkerStatus>("hidden");

    const [dialogData, setDialogData] = useState<DialogData>({
        isOpen: false,
        title: "",
        pointId: "",
        subtitle: "",
        logoBgColor: "",
        logoTextColor: "",
        shortName: "",
    });

    const currentSearchQuery = searchQuery.value;
    const currentPointTypes = selectedPointTypes.value;
    const currentPointStatuses = selectedPointStatuses.value;
    const currentDistanceKm = selectedDistanceKm.value;
    const localizationEnabled = isLocalizationEnabled.value;
    const currentLocation = userLocation.value;

    const displayPoints = useMemo(() => {
        return (initialReturnPoints.points ?? []).map(getDisplayMapPoint);
    }, [initialReturnPoints.points]);

    const selectedDisplayPoint = useMemo(() => {
        if (selectedPoint === undefined) return undefined;

        return displayPoints.find((point) => (
            point.returnPointId === selectedPoint && isValidMapPoint(point)
        ));
    }, [displayPoints, selectedPoint]);

    const filterMetadataUnavailable = (
        currentPointTypes.length > 0
        && !displayPoints.some((point) => point.pointType !== undefined)
    ) || (
        currentPointStatuses.length > 0
        && !displayPoints.some((point) => point.status !== undefined)
    );

    const filteredPoints = useMemo(() => {
        const normalizedQuery = normalizeSearchText(currentSearchQuery);

        return displayPoints.filter((point) => {
            if (!isValidMapPoint(point)) return false;

            const pointType = point.pointType;
            const pointStatus = point.status;

            if (
                currentPointTypes.length > 0
                && (pointType === undefined || !currentPointTypes.includes(pointType))
            ) return false;

            if (
                currentPointStatuses.length > 0
                && (pointStatus === undefined || !currentPointStatuses.includes(pointStatus))
            ) return false;

            if (normalizedQuery) {
                const searchableText = normalizeSearchText([
                    point.name,
                    point.shortName,
                    point.address,
                    point.postalCode,
                ].filter(Boolean).join(" "));

                if (!searchableText.includes(normalizedQuery)) return false;
            }

            if (localizationEnabled && currentLocation) {
                const distanceKm = getDistanceKm(
                    currentLocation.latitude,
                    currentLocation.longitude,
                    point.latitude,
                    point.longitude,
                );

                if (distanceKm > currentDistanceKm) return false;
            }

            return true;
        });
    }, [
        displayPoints,
        currentSearchQuery,
        currentPointTypes,
        currentPointStatuses,
        currentDistanceKm,
        localizationEnabled,
        currentLocation,
    ]);

    const filteredPointsGeoJson = useMemo<MapPointFeatureCollection>(() => ({
        type: "FeatureCollection",
        features: filteredPoints.map((point) => ({
            type: "Feature",
            id: point.returnPointId,
            properties: {
                returnPointId: point.returnPointId,
                iconId: getMapPinImageId(
                    point.bgColor,
                    point.textColor,
                    point.shortName,
                ),
            },
            geometry: {
                type: "Point",
                coordinates: [point.longitude, point.latitude],
            },
        })),
    }), [filteredPoints]);

    const userLocationGeoJson = useMemo(() => {
        return getUserLocationGeoJson(localizationEnabled ? currentLocation : null);
    }, [localizationEnabled, currentLocation]);

    const syncRenderedFeatureCounts = useCallback((): void => {
        const map = mapRef.current;
        const mapRoot = mapContainerRef.current?.closest<HTMLElement>("[data-map]");
        if (
            !map
            || !mapRoot
            || !map.getLayer(MAP_PINS_LAYER_ID)
            || !map.getLayer(CLUSTERS_LAYER_ID)
        ) return;

        const renderedPointIds = new Set(
            map.queryRenderedFeatures({ layers: [MAP_PINS_LAYER_ID] })
                .map((feature) => feature.properties?.returnPointId)
                .filter((returnPointId): returnPointId is string => (
                    typeof returnPointId === "string"
                )),
        );
        const renderedClusterIds = new Set(
            map.queryRenderedFeatures({ layers: [CLUSTERS_LAYER_ID] })
                .map((feature) => feature.properties?.cluster_id)
                .filter((clusterId): clusterId is number => typeof clusterId === "number"),
        );

        mapRoot.dataset.renderedMarkers = String(renderedPointIds.size);
        mapRoot.dataset.renderedClusters = String(renderedClusterIds.size);
        mapRoot.dataset.domMarkers = "0";
    }, []);

    useEffect(() => {
        const mapContainer = mapContainerRef.current;
        if (!mapContainer) return;

        setMapStatus("loading");

        let isDisposed = false;
        let resizeObserver: ResizeObserver | undefined;

        const initializeMap = async (): Promise<void> => {
            try {
                const mapLibreModule = await import("maplibre-gl");
                if (isDisposed) return;

                const map = new mapLibreModule.Map({
                    container: mapContainer,
                    style: MAP_STYLE_URL,
                    center: DEFAULT_CENTER,
                    zoom: 5.5,
                    minZoom: 4,
                    maxZoom: 18,
                    attributionControl: { compact: true },
                });

                mapRef.current = map;
                const mapRoot = mapContainer.closest<HTMLElement>("[data-map]");
                const updateCameraData = (): void => {
                    if (!mapRoot) return;

                    const center = map.getCenter();
                    mapRoot.dataset.mapZoom = map.getZoom().toFixed(2);
                    mapRoot.dataset.mapCenter = `${center.lng.toFixed(5)},${center.lat.toFixed(5)}`;
                };

                map.once("load", () => {
                    if (isDisposed) return;

                    localizeMapLabels(map);

                    const mapPinImages = new globalThis.Map<string, {
                        bgColor: string;
                        textColor: string;
                        shortName: string;
                    }>();
                    displayPoints.forEach((point) => {
                        mapPinImages.set(
                            getMapPinImageId(
                                point.bgColor,
                                point.textColor,
                                point.shortName,
                            ),
                            {
                                bgColor: point.bgColor,
                                textColor: point.textColor,
                                shortName: point.shortName,
                            },
                        );
                    });
                    mapPinImages.forEach((image, imageId) => {
                        if (map.hasImage(imageId)) return;

                        map.addImage(
                            imageId,
                            createMapPinImageData(
                                image.bgColor,
                                image.textColor,
                                image.shortName,
                            ),
                            { pixelRatio: MAP_PIN_IMAGE_PIXEL_RATIO },
                        );
                    });
                    if (mapRoot) mapRoot.dataset.mapPinImages = String(mapPinImages.size);

                    map.addSource(POINTS_SOURCE_ID, {
                        type: "geojson",
                        data: {
                            type: "FeatureCollection",
                            features: [],
                        },
                        cluster: true,
                        clusterRadius: 38,
                        clusterMaxZoom: MAP_PIN_MIN_ZOOM - 1,
                    });

                    map.addLayer({
                        id: CLUSTERS_LAYER_ID,
                        type: "circle",
                        source: POINTS_SOURCE_ID,
                        maxzoom: MAP_PIN_MIN_ZOOM,
                        filter: ["has", "point_count"],
                        paint: {
                            "circle-color": [
                                "step",
                                ["get", "point_count"],
                                "#2f8f46",
                                100,
                                "#19723d",
                                750,
                                "#0e4328",
                            ],
                            "circle-radius": [
                                "step",
                                ["get", "point_count"],
                                20,
                                100,
                                26,
                                750,
                                32,
                            ],
                            "circle-stroke-color": "#ffffff",
                            "circle-stroke-width": 2,
                        },
                    });

                    map.addLayer({
                        id: CLUSTER_COUNT_LAYER_ID,
                        type: "symbol",
                        source: POINTS_SOURCE_ID,
                        maxzoom: MAP_PIN_MIN_ZOOM,
                        filter: ["has", "point_count"],
                        layout: {
                            "text-field": ["get", "point_count_abbreviated"],
                            "text-font": ["Noto Sans Bold"],
                            "text-size": 13,
                            "text-allow-overlap": true,
                        },
                        paint: {
                            "text-color": "#ffffff",
                            "text-halo-color": "rgba(0, 0, 0, .2)",
                            "text-halo-width": 1,
                        },
                    });

                    map.addLayer({
                        id: MAP_PINS_LAYER_ID,
                        type: "symbol",
                        source: POINTS_SOURCE_ID,
                        minzoom: MAP_PIN_MIN_ZOOM,
                        filter: ["!", ["has", "point_count"]],
                        layout: {
                            "icon-image": ["get", "iconId"],
                            "icon-anchor": "bottom",
                            "icon-allow-overlap": true,
                            "icon-ignore-placement": true,
                        },
                    });

                    map.addLayer({
                        id: HOVERED_MAP_PIN_LAYER_ID,
                        type: "symbol",
                        source: POINTS_SOURCE_ID,
                        minzoom: MAP_PIN_MIN_ZOOM,
                        filter: [
                            "all",
                            ["!", ["has", "point_count"]],
                            ["==", ["get", "returnPointId"], ""],
                        ],
                        layout: {
                            "icon-image": ["get", "iconId"],
                            "icon-anchor": "bottom",
                            "icon-size": 1.25,
                            "icon-allow-overlap": true,
                            "icon-ignore-placement": true,
                        },
                    });

                    map.addSource(USER_LOCATION_SOURCE_ID, {
                        type: "geojson",
                        data: {
                            type: "FeatureCollection",
                            features: [],
                        },
                    });

                    map.addLayer({
                        id: USER_LOCATION_HALO_LAYER_ID,
                        type: "circle",
                        source: USER_LOCATION_SOURCE_ID,
                        paint: {
                            "circle-radius": 16,
                            "circle-color": "#1a73e8",
                            "circle-opacity": .2,
                            "circle-stroke-color": "#1a73e8",
                            "circle-stroke-opacity": .35,
                            "circle-stroke-width": 1,
                        },
                    });

                    map.addLayer({
                        id: USER_LOCATION_POINT_LAYER_ID,
                        type: "circle",
                        source: USER_LOCATION_SOURCE_ID,
                        paint: {
                            "circle-radius": 7,
                            "circle-color": "#1a73e8",
                            "circle-stroke-color": "#ffffff",
                            "circle-stroke-width": 3,
                        },
                    });

                    const handleClusterClick = (event: MapLayerMouseEvent): void => {
                        const cluster = event.features?.[0];
                        if (!cluster || cluster.geometry.type !== "Point") return;

                        const clusterId = Number(cluster.properties?.cluster_id);
                        if (!Number.isFinite(clusterId)) return;

                        const source = map.getSource(POINTS_SOURCE_ID) as GeoJSONSource | undefined;
                        if (!source) return;

                        void source.getClusterExpansionZoom(clusterId).then((zoom) => {
                            if (isDisposed) return;

                            map.easeTo({
                                center: cluster.geometry.coordinates as [number, number],
                                zoom: Math.min(zoom, map.getMaxZoom()),
                                duration: getAnimationDuration(),
                            });
                        });
                    };

                    const getClosestMapPin = (
                        event: MapLayerMouseEvent,
                    ): MapGeoJSONFeature | undefined => {
                        let closestPoint: MapGeoJSONFeature | undefined;
                        let closestDistanceSquared = Number.POSITIVE_INFINITY;

                        event.features?.forEach((point) => {
                            if (point.geometry.type !== "Point") return;

                            const markerAnchor = map.project(
                                point.geometry.coordinates as [number, number],
                            );
                            const horizontalDistance = event.point.x - markerAnchor.x;
                            const verticalDistance = event.point.y - (markerAnchor.y - 32);
                            const distanceSquared = horizontalDistance ** 2 + verticalDistance ** 2;

                            if (distanceSquared >= closestDistanceSquared) return;

                            closestPoint = point;
                            closestDistanceSquared = distanceSquared;
                        });

                        return closestPoint;
                    };
                    const handleMapPinClick = (event: MapLayerMouseEvent): void => {
                        const point = getClosestMapPin(event);

                        if (!point || point.geometry.type !== "Point") return;

                        const [longitude, latitude] = point.geometry.coordinates;

                        centerMapOnPoint(map, {
                            latitude,
                            longitude,
                        });

                        setTimeout(() => {
                            const returnPointId = point.properties?.returnPointId;

                            if (typeof returnPointId !== "string") return;

                            setDialogData({
                                isOpen: true,
                                pointId: returnPointId,
                                title: "Żabka",
                                subtitle: "ul. Szczecińska 18",
                                logoBgColor: initialReturnPoints.points.find((p) => p.returnPointId === returnPointId)?.bgColor ?? "#000000",
                                logoTextColor: initialReturnPoints.points.find((p) => p.returnPointId === returnPointId)?.textColor ?? "#ffffff",
                                shortName: initialReturnPoints.points.find((p) => p.returnPointId === returnPointId)?.shortName ?? "",
                            });
                        }, getAnimationDuration() + 10);
                    };
                    const showInteractiveCursor = (): void => {
                        map.getCanvas().style.cursor = "pointer";
                    };
                    const resetInteractiveCursor = (): void => {
                        map.getCanvas().style.cursor = "";
                    };
                    let hoveredMapPinId: string | null = null;
                    const setHoveredMapPin = (returnPointId: string | null): void => {
                        if (hoveredMapPinId === returnPointId) return;

                        hoveredMapPinId = returnPointId;
                        map.setFilter(HOVERED_MAP_PIN_LAYER_ID, [
                            "all",
                            ["!", ["has", "point_count"]],
                            ["==", ["get", "returnPointId"], returnPointId ?? ""],
                        ]);
                        if (mapRoot) mapRoot.dataset.hoveredMapPin = returnPointId ?? "";
                    };
                    const handleMapPinHover = (event: MapLayerMouseEvent): void => {
                        const returnPointId = getClosestMapPin(event)?.properties?.returnPointId;
                        if (typeof returnPointId !== "string") return;

                        setHoveredMapPin(returnPointId);
                        showInteractiveCursor();
                    };
                    const resetMapPinHover = (): void => {
                        setHoveredMapPin(null);
                        resetInteractiveCursor();
                    };

                    map.on("click", CLUSTERS_LAYER_ID, handleClusterClick);
                    map.on("click", MAP_PINS_LAYER_ID, handleMapPinClick);
                    map.on("mouseenter", CLUSTERS_LAYER_ID, showInteractiveCursor);
                    map.on("mouseleave", CLUSTERS_LAYER_ID, resetInteractiveCursor);
                    map.on("mouseenter", MAP_PINS_LAYER_ID, handleMapPinHover);
                    map.on("mousemove", MAP_PINS_LAYER_ID, handleMapPinHover);
                    map.on("mouseleave", MAP_PINS_LAYER_ID, resetMapPinHover);
                    map.on("idle", syncRenderedFeatureCounts);

                    if (mapRoot) {
                        mapRoot.dataset.clustering = "enabled";
                        mapRoot.dataset.markerRenderer = "maplibre-webgl";
                        mapRoot.dataset.domMarkers = "0";
                    }
                    updateCameraData();
                    setMapStatus("ready");
                });

                map.on("moveend", updateCameraData);

                map.on("error", () => {
                    if (!isDisposed && !map.loaded()) setMapStatus("error");
                });

                resizeObserver = new ResizeObserver(() => map.resize());
                resizeObserver.observe(mapContainer);
            } catch {
                if (!isDisposed) setMapStatus("error");
            }
        };

        void initializeMap();

        return () => {
            isDisposed = true;
            resizeObserver?.disconnect();

            mapRef.current?.remove();
            mapRef.current = null;
        };
    }, [displayPoints, syncRenderedFeatureCounts, initialReturnPoints.points]);

    useEffect(() => {
        const map = mapRef.current;
        if (mapStatus !== "ready" || !map) return;

        const pointsSource = map.getSource(POINTS_SOURCE_ID) as GeoJSONSource | undefined;
        if (!pointsSource) return;

        const sourceUpdateId = ++sourceUpdateIdRef.current;
        const mapRoot = mapContainerRef.current?.closest<HTMLElement>("[data-map]");
        if (mapRoot) {
            mapRoot.dataset.renderedMarkers = "0";
            mapRoot.dataset.renderedClusters = "0";
        }
        void pointsSource.setData(filteredPointsGeoJson, true)
            .then(() => {
                if (
                    sourceUpdateId === sourceUpdateIdRef.current
                    && mapRef.current === map
                ) syncRenderedFeatureCounts();
            })
            .catch(() => {
                if (mapRef.current === map) setMapStatus("error");
            });

        if (filteredPoints.length > 0) {
            fitMapToPoints(map, filteredPoints);
        } else if (localizationEnabled && currentLocation) {
            map.easeTo({
                center: [currentLocation.longitude, currentLocation.latitude],
                zoom: Math.max(map.getZoom(), 12),
                duration: getAnimationDuration(),
            });
        }
    }, [
        filteredPoints,
        filteredPointsGeoJson,
        mapStatus,
        localizationEnabled,
        currentLocation,
        syncRenderedFeatureCounts,
    ]);

    useEffect(() => {
        const map = mapRef.current;
        if (mapStatus !== "ready" || !map || !selectedDisplayPoint) return;

        centerMapOnPoint(map, selectedDisplayPoint);
    }, [mapStatus, selectedDisplayPoint]);

    useEffect(() => {
        const map = mapRef.current;
        if (mapStatus !== "ready" || !map) return;

        const locationSource = map.getSource(USER_LOCATION_SOURCE_ID) as GeoJSONSource | undefined;
        if (!locationSource) {
            setUserLocationMarkerStatus("error");
            return;
        }

        const updateId = ++userLocationUpdateIdRef.current;
        setUserLocationMarkerStatus("updating");

        void locationSource.setData(userLocationGeoJson, true)
            .then(() => {
                if (
                    mapRef.current !== map
                    || updateId !== userLocationUpdateIdRef.current
                ) return;

                setUserLocationMarkerStatus(
                    userLocationGeoJson.features.length > 0 ? "visible" : "hidden",
                );
            })
            .catch(() => {
                if (
                    mapRef.current === map
                    && updateId === userLocationUpdateIdRef.current
                ) setUserLocationMarkerStatus("error");
            });
    }, [mapStatus, userLocationGeoJson]);

    const handleZoomIn = (): void => {
        mapRef.current?.zoomIn({ duration: getAnimationDuration() });
    };

    const handleZoomOut = (): void => {
        mapRef.current?.zoomOut({ duration: getAnimationDuration() });
    };

    const handleCenterOnLocation = (): void => {
        if (!currentLocation) return;

        mapRef.current?.easeTo({
            center: [currentLocation.longitude, currentLocation.latitude],
            zoom: Math.max(mapRef.current.getZoom(), 14),
            duration: getAnimationDuration(),
        });
    };

    const isMapReady = mapStatus === "ready";
    const canCenterOnLocation = isMapReady && localizationEnabled && Boolean(currentLocation);

    return (
        <section
            class={Styles["map"]}
            aria-label="Interaktywna mapa punktów zwrotu"
            aria-busy={mapStatus === "loading"}
            data-map=""
            data-map-status={mapStatus}
            data-visible-markers={filteredPoints.length}
            data-clustering="enabled"
            data-marker-renderer="maplibre-webgl"
            data-dom-markers="0"
            data-map-pin-images="0"
            data-rendered-markers="0"
            data-rendered-clusters="0"
            data-user-location-status={userLocationMarkerStatus}
            data-user-location-visible={userLocationMarkerStatus === "visible" ? "true" : "false"}
        >
            <div
                ref={mapContainerRef}
                class={Styles["map__canvas"]}
                aria-label="Mapa punktów zwrotu"
                data-map-canvas=""
            />

            <div
                class={`${Styles["map__controls"]} ${Styles["map__controls--zoom"]}`}
                role="group"
                aria-label="Sterowanie przybliżeniem mapy"
                data-map-zoom-controls=""
            >
                <IconButton
                    class={Styles["map__control-button"]}
                    icon="mdi:plus"
                    aria-label="Przybliż mapę"
                    disabled={!isMapReady}
                    data-map-zoom-in=""
                    onClick={handleZoomIn}
                />

                <IconButton
                    class={Styles["map__control-button"]}
                    icon="mdi:minus"
                    aria-label="Oddal mapę"
                    disabled={!isMapReady}
                    data-map-zoom-out=""
                    onClick={handleZoomOut}
                />
            </div>

            <div class={`${Styles["map__controls"]} ${Styles["map__controls--location"]}`}>
                <IconButton
                    class={Styles["map__control-button"]}
                    icon="mdi:crosshairs-gps"
                    aria-label="Wyśrodkuj mapę na mojej lokalizacji"
                    title={!localizationEnabled
                        ? "Włącz lokalizację, aby wyśrodkować mapę"
                        : undefined}
                    disabled={!canCenterOnLocation}
                    data-map-center-location=""
                    onClick={handleCenterOnLocation}
                />
            </div>

            {mapStatus === "loading" && (
                <div class={Styles["map__message"]} role="status">
                    Ładowanie mapy…
                </div>
            )}

            {mapStatus === "error" && (
                <div class={Styles["map__message"]} role="alert">
                    Nie udało się załadować mapy. Sprawdź połączenie z internetem.
                </div>
            )}

            {isMapReady && filteredPoints.length === 0 && (
                <div class={Styles["map__message"]} role="status">
                    {filterMetadataUnavailable
                        ? "Dane punktów nie zawierają informacji wymaganych przez wybrany filtr."
                        : "Brak punktów spełniających wybrane kryteria."}
                </div>
            )}

            <Modal 
              bodyClass={Styles["map__modal-body"]}
              title={dialogData.title} 
              subtitle={dialogData.subtitle} 
              open={dialogData.isOpen} 
              onClose={() => {setDialogData({ ...dialogData, isOpen: false });}}
            >
              <div 
                class={Styles["map__modal-body__card"]}
                style={{"--brand-bg-color": dialogData.logoBgColor, "--brand-text-color": dialogData.logoTextColor}}
              >
                <div class={Styles["map__modal-body__card__header"]}>
                  <span class={Styles["map__modal-body__card__header__brand"]}>{dialogData.shortName}</span>

                  <div class={Styles["map__modal-body__card__header__content"]}>
                    <span class={Styles["map__modal-body__card__header__content__badge"]}>Brak informacji</span>
                    <h4>Żabka</h4>
                    <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Perferendis cum, fugit maiores deleniti eveniet, porro rerum sint, alias aliquid asperiores voluptate! Minus cupiditate!</p>
                  </div>
                </div>
              </div>

              <div class={Styles["map__modal-body__footer"]}>
                <Button class={`${Styles["map__modal-body__footer__button"]} ${Styles["map__modal-body__footer__button--navigation"]}`} variant="contained">Wyznacz trasę</Button>

                <Button class={Styles["map__modal-body__footer__button"]} variant="contained">Szczegóły</Button>
              </div>
            </Modal>
        </section>
    );
}