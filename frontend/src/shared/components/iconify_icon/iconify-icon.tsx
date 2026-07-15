import Styles from "./iconify-icon.module.sass";

import type { CSSProperties, HTMLAttributes } from "preact";

export type IconifyIconFlip = "both" | "horizontal" | "vertical";
export type IconifyIconRotation = 0 | 90 | 180 | 270;

type IconSize = number | string;
type AriaHidden = boolean | "false" | "true";

type NativeSpanProps = Omit<
	HTMLAttributes<HTMLSpanElement>,
	| "aria-hidden"
	| "aria-label"
	| "aria-labelledby"
	| "children"
	| "class"
	| "className"
	| "role"
	| "style"
>;

export type IconifyIconProps = NativeSpanProps & {
	icon: string;
	size?: IconSize;
	width?: IconSize;
	height?: IconSize;
	label?: string;
	colored?: boolean;
	color?: string;
	rotate?: IconifyIconRotation;
	flip?: IconifyIconFlip;
	class?: string;
	className?: string;
	style?: CSSProperties | string;
	role?: HTMLAttributes<HTMLSpanElement>["role"];
	"aria-hidden"?: AriaHidden;
	"aria-label"?: string;
	"aria-labelledby"?: string;
};

const ICON_NAME_PATTERN = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

function getIconUrl(icon: string): string | undefined {
	const [prefix, name, ...extraParts] = icon.split(":");

	if (
		extraParts.length > 0
		|| !prefix
		|| !name
		|| !ICON_NAME_PATTERN.test(prefix)
		|| !ICON_NAME_PATTERN.test(name)
	) return undefined;

	return `url("https://api.iconify.design/${encodeURIComponent(prefix)}/${encodeURIComponent(name)}.svg")`;
}

function normalizeSize(size: IconSize | undefined): string | undefined {
	if (typeof size === "number")
		return Number.isFinite(size) && size >= 0 ? `${size}px` : undefined;

	return size?.trim() || undefined;
}

function normalizeText(value: string | undefined): string | undefined {
	return value?.trim() || undefined;
}

function getIconStyle(
	style: CSSProperties | string | undefined,
	iconUrl: string,
	width: IconSize | undefined,
	height: IconSize | undefined,
	size: IconSize | undefined,
	color: string | undefined,
	rotate: IconifyIconRotation | undefined,
	flip: IconifyIconFlip | undefined,
): CSSProperties | string {
	const iconStyle: CSSProperties = {
		"--iconify-icon-url": iconUrl,
		"--iconify-icon-width": normalizeSize(width ?? size),
		"--iconify-icon-height": normalizeSize(height ?? size),
		"--iconify-icon-rotation": rotate === undefined ? undefined : `${rotate}deg`,
		"--iconify-icon-scale-x": flip === "horizontal" || flip === "both" ? -1 : undefined,
		"--iconify-icon-scale-y": flip === "vertical" || flip === "both" ? -1 : undefined,
		color: normalizeText(color),
	};

	if (typeof style !== "string") return { ...style, ...iconStyle };

	const generatedStyle = Object.entries(iconStyle)
		.filter((entry): entry is [string, number | string] => entry[1] !== undefined)
		.map(([property, value]) => `${property}: ${value}`)
		.join("; ");

	return [style.trim().replace(/;+$/, ""), generatedStyle].filter(Boolean).join("; ");
}

export default function IconifyIcon({
	icon,
	size,
	width,
	height,
	label,
	colored = false,
	color,
	rotate,
	flip,
	class: className,
	className: nativeClassName,
	style,
	role,
	"aria-hidden": ariaHidden,
	"aria-label": ariaLabel,
	"aria-labelledby": ariaLabelledBy,
	...spanProps
}: IconifyIconProps) {
	const iconUrl = getIconUrl(icon);

	if (!iconUrl) return null;

	const normalizedAriaLabel = normalizeText(ariaLabel) ?? normalizeText(label);

	const normalizedAriaLabelledBy = normalizeText(ariaLabelledBy);

	const hasAccessibleName = Boolean(normalizedAriaLabel || normalizedAriaLabelledBy);

	const isExplicitlyHidden = ariaHidden === true || ariaHidden === "true";

	const iconifyIconClasses = [
		Styles["iconify-icon"],
		colored ? Styles["iconify-icon--colored"] : undefined,
		className,
		nativeClassName,
	].filter(Boolean).join(" ");

	return (
		<span
			{...spanProps}
			class={iconifyIconClasses}
			style={getIconStyle(style, iconUrl, width, height, size, color, rotate, flip)}
			role={role ?? (hasAccessibleName && !isExplicitlyHidden ? "img" : undefined)}
			aria-label={normalizedAriaLabel}
			aria-labelledby={normalizedAriaLabelledBy}
			aria-hidden={ariaHidden ?? (hasAccessibleName ? undefined : "true")}
		/>
	);
}
