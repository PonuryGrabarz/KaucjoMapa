import {
	resolveMarkerText,
	resolveMarkerLabel,
	resolveMarkerBackground,
} from "@shared/helpers/marker-appearance";

import type { ReturnPoint, ReturnPointStatus } from "@/shared/api/responses";

const STATUS_PIN_COLORS: Record<ReturnPointStatus, string> = {
	dostepny: "#1f7a3c",
	nieczynny: "#5a6b60",
	brak_informacji: "#b75c17",
};

const MARKER_FONT_FALLBACK = "ui-sans-serif, system-ui, sans-serif";

export interface MarkerStyle {
	backgroundColor: string;
	textColor: string;
	shortName: string;
	status: ReturnPointStatus;
}

let markerFont: string | undefined;

const resolveMarkerFont = (): string => {
	if (markerFont) return markerFont;

	const family = getComputedStyle(document.documentElement)
		.getPropertyValue("--font-ui")
		.trim();

	markerFont = `800 12px ${family ? `${family}, ` : ""}${MARKER_FONT_FALLBACK}`;

	return markerFont;
}

export const supportsWebgl = (): boolean => {
  try {
    const canvas = document.createElement("canvas");

    return Boolean(canvas.getContext("webgl2") ?? canvas.getContext("webgl"));
  } catch {
    return false;
  }
}

export const drawMarkerImage = (style: MarkerStyle): ImageData => {
	const pixelRatio = 2;
	const imageSize = 64;

	const canvas = document.createElement("canvas");

	canvas.width = imageSize * pixelRatio;
	canvas.height = imageSize * pixelRatio;

	const context = canvas.getContext("2d")!;

	context.scale(pixelRatio, pixelRatio);
	context.translate(imageSize / 2, imageSize / 2);

	context.save();

	context.rotate(Math.PI / 4);

	context.beginPath();

	context.roundRect(-20, -20, 40, 40, [20, 20, 20, 6]);

	context.fillStyle = resolveMarkerBackground(style.backgroundColor);
	context.strokeStyle = "#ffffff";
	context.lineWidth = 3;
	context.shadowColor = "rgb(14 67 40 / 0.28)";
	context.shadowBlur = 6;
	context.shadowOffsetY = 2;

	context.fill();

	context.shadowColor = "transparent";

	context.stroke();
	context.restore();

	const markerText = resolveMarkerLabel(style.shortName);

	context.fillStyle = resolveMarkerText(style.textColor);
	context.font = resolveMarkerFont();
	context.textAlign = "left";
	context.textBaseline = "alphabetic";

	const textMetrics = context.measureText(markerText);
  
	const textX = (textMetrics.actualBoundingBoxLeft - textMetrics.actualBoundingBoxRight) / 2;
	const textY = (textMetrics.actualBoundingBoxAscent - textMetrics.actualBoundingBoxDescent) / 2;

	context.fillText(markerText, textX, textY, 30);

	const statusColor = STATUS_PIN_COLORS[style.status];

	context.beginPath();

	context.arc(15, -14, 7, 0, Math.PI * 2);
	context.fillStyle = "#ffffff";

	context.fill();

	context.strokeStyle = "rgb(14 42 27 / 0.16)";
	context.lineWidth = 1;

	context.stroke();

	context.strokeStyle = statusColor;
	context.fillStyle = statusColor;
	context.lineWidth = 1.9;
	context.lineCap = "round";
	context.lineJoin = "round";

	if (style.status === "dostepny") {
		context.beginPath();
		context.moveTo(11.5, -14);
		context.lineTo(14, -11.5);
		context.lineTo(18.5, -17);
		context.stroke();
	} else if (style.status === "brak_informacji") {
		context.beginPath();
		context.moveTo(15, -17.5);
		context.lineTo(15, -13);
		context.stroke();
		context.beginPath();
		context.arc(15, -10.5, 1.1, 0, Math.PI * 2);
		context.fill();
	} else {
		context.beginPath();
		context.moveTo(11.5, -14);
		context.lineTo(18.5, -14);
		context.stroke();
	}

	return context.getImageData(0, 0, canvas.width, canvas.height);
}

export const markerStyleOf = (point: ReturnPoint): MarkerStyle => ({
	backgroundColor: point.bg_color,
	textColor: point.text_color,
	shortName: point.short_name,
	status: point.current_status,
});

export const markerImageKey = (style: MarkerStyle): string =>
	`${style.backgroundColor}|${style.textColor}|${style.shortName}|${style.status}`;
