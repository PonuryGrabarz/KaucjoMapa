import { mapState, RADIUS_MAX_KM } from "./map-state.svelte";

import {
	RETURN_POINT_TYPES,
	RETURN_POINT_STATUSES,
	RETURN_POINT_MATERIALS,
	RETURN_POINT_REFUND_METHODS,
} from "@shared/api/responses";

import type { SvelteSet } from "svelte/reactivity";

const PARAM_TYPE = "type";
const PARAM_POINT = "point";
const PARAM_OPEN = "open";
const PARAM_STATUS = "status";
const PARAM_REFUND = "refund";
const PARAM_MATERIAL = "material";
const PARAM_MIN_RADIUS = "rmin";
const PARAM_MAX_RADIUS = "rmax";

const UUID_PATTERN = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

const MANAGED_PARAMS = [
	PARAM_TYPE,
	PARAM_POINT,
	PARAM_OPEN,
	PARAM_STATUS,
	PARAM_REFUND,
	PARAM_MATERIAL,
	PARAM_MIN_RADIUS,
	PARAM_MAX_RADIUS,
];

const serializeSet = <TValue extends string>(
	set: SvelteSet<TValue>,
	allowed: readonly TValue[],
): string | undefined => (set.size === allowed.length
	? undefined
	: allowed.filter((value) => set.has(value)).join(","));

const applySet = <TValue extends string>(
	set: SvelteSet<TValue>,
	raw: string | null,
	allowed: readonly TValue[],
): void => {
	if (raw === null) return;

	const wanted = raw
		.split(",")
		.filter((value): value is TValue => (allowed as readonly string[]).includes(value));

	set.clear();

	for (const value of wanted) set.add(value);
}

const applyNumber = (raw: string | null, fallback: number): number => {
	if (raw === null) return fallback;

	const value = Number(raw);

	if (!Number.isFinite(value)) return fallback;

	return Math.min(Math.max(Math.round(value), 0), RADIUS_MAX_KM);
}

const mapStateToSearchParams = (): URLSearchParams => {
	const params = new URLSearchParams();

	const entries: [string, string | undefined][] = [
		[PARAM_STATUS, serializeSet(mapState.statuses, RETURN_POINT_STATUSES)],
		[PARAM_TYPE, serializeSet(mapState.types, RETURN_POINT_TYPES)],
		[PARAM_MATERIAL, serializeSet(mapState.materials, RETURN_POINT_MATERIALS)],
		[PARAM_REFUND, serializeSet(mapState.refundMethods, RETURN_POINT_REFUND_METHODS)],
		[PARAM_POINT, mapState.selectedReturnPointId],
		[PARAM_OPEN, mapState.openNowOnly ? "1" : undefined],
		[PARAM_MIN_RADIUS, mapState.minRadiusKm > 0 ? String(mapState.minRadiusKm) : undefined],
		[
			PARAM_MAX_RADIUS,
			mapState.maxRadiusKm < RADIUS_MAX_KM ? String(mapState.maxRadiusKm) : undefined,
		],
	];

	for (const [key, value] of entries) if (value !== undefined) params.set(key, value);

	return params;
}

export const applySearchParamsToMapState = (search: string): void => {
	const params = new URLSearchParams(search);

	applySet(mapState.statuses, params.get(PARAM_STATUS), RETURN_POINT_STATUSES);
	applySet(mapState.types, params.get(PARAM_TYPE), RETURN_POINT_TYPES);
	applySet(mapState.materials, params.get(PARAM_MATERIAL), RETURN_POINT_MATERIALS);
	applySet(mapState.refundMethods, params.get(PARAM_REFUND), RETURN_POINT_REFUND_METHODS);

	const point = params.get(PARAM_POINT);

	mapState.selectedReturnPointId = point && UUID_PATTERN.test(point) ? point : undefined;

	mapState.openNowOnly = params.get(PARAM_OPEN) === "1";

	const min = applyNumber(params.get(PARAM_MIN_RADIUS), 0);
	const max = applyNumber(params.get(PARAM_MAX_RADIUS), RADIUS_MAX_KM);

	mapState.minRadiusKm = Math.min(min, max);
	mapState.maxRadiusKm = Math.max(min, max);
}

export const writeMapStateToUrl = (): void => {
	const url = new URL(window.location.href);

	for (const key of MANAGED_PARAMS) url.searchParams.delete(key);

	for (const [key, value] of mapStateToSearchParams()) url.searchParams.set(key, value);

	if (url.href === window.location.href) return;

	history.replaceState(history.state, "", url);
}
