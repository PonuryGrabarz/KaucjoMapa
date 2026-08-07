export const SOURCE_ID = "return-points";

export const POINT_LAYER = "return-point-symbols";
export const CLUSTER_LAYER = "return-point-clusters";
export const POINT_LABEL_LAYER = "return-point-label";
export const POINT_HOVER_LAYER = "return-point-symbols-hover";
export const CLUSTER_COUNT_LAYER = "return-point-cluster-count";
export const POINT_SELECTED_LAYER = "return-point-symbols-selected";

export const USER_LOCATION_SOURCE_ID = "user-location";
export const USER_LOCATION_DOT_LAYER = "user-location-dot";
export const USER_LOCATION_ACCURACY_LAYER = "user-location-accuracy";

export const POINT_ID_PROPERTY = "returnPointId";

export const MAP_STYLE_URL = "https://tiles.openfreemap.org/styles/positron";

export const MAP_PIN_MIN_ZOOM = 9;

export const POLAND_CENTER: [number, number] = [19.1451, 51.9194];

export const MAP_LOAD_TIMEOUT_MS = 9000;

export const POLISH_MAP_LABEL_EXPRESSION = [
  "coalesce",
  ["get", "name:pl"],
  ["get", "name"],
  ["get", "name_en"],
];

export const MAP_FONT_REGULAR = ["Noto Sans Regular"];
export const MAP_FONT_BOLD = ["Noto Sans Bold"];
