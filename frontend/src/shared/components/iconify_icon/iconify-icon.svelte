<template src="./iconify-icon.svelte.html"></template>

<style lang="sass" src="./iconify-icon.sass"></style>

<script lang="ts">
import type { HTMLAttributes } from "svelte/elements";

export type IconifyIconFlip = "both" | "horizontal" | "vertical";
export type IconifyIconRotation = 0 | 90 | 180 | 270;

type IconSize = number | string;

export type IconifyIconProps = Omit<HTMLAttributes<HTMLSpanElement>, "children"> & {
	icon: string;
	size?: IconSize;
	width?: IconSize;
	height?: IconSize;
	label?: string;
	colored?: boolean;
	color?: string;
	rotate?: IconifyIconRotation;
	flip?: IconifyIconFlip;
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

function normalizeText(value: string | null | undefined): string | undefined {
	return value?.trim() || undefined;
}

let {
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
	role,
	"aria-hidden": ariaHidden,
	"aria-label": ariaLabel,
	"aria-labelledby": ariaLabelledBy,
	...rest
}: IconifyIconProps = $props();

const iconUrl = $derived(getIconUrl(icon));

const accessibleLabel = $derived(normalizeText(ariaLabel) ?? normalizeText(label));
const labelledBy = $derived(normalizeText(ariaLabelledBy));
const hasAccessibleName = $derived(Boolean(accessibleLabel || labelledBy));
const explicitlyHidden = $derived(ariaHidden === true || ariaHidden === "true");
</script>
