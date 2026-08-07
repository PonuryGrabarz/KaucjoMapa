import { isHexColor } from "@shared/helpers/color";

export const DEFAULT_MARKER_BACKGROUND = "#2f8f46";
export const DEFAULT_MARKER_TEXT = "#ffffff";

export const MARKER_LABEL_MAX_LENGTH = 3;

const resolveColor = (value: string | null | undefined, fallback: string): string =>
	isHexColor(value) ? value!.trim() : fallback;

export const resolveMarkerBackground = (value: string | null | undefined): string =>
	resolveColor(value, DEFAULT_MARKER_BACKGROUND);

export const resolveMarkerText = (value: string | null | undefined): string =>
	resolveColor(value, DEFAULT_MARKER_TEXT);

export const resolveMarkerLabel = (shortName: string | null | undefined): string =>
	shortName?.trim().slice(0, MARKER_LABEL_MAX_LENGTH) || "?";
