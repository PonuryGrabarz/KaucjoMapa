import { isOpenNow } from "@shared/helpers/opening-hours";

import {
	RETURN_POINT_TYPES,
	RETURN_POINT_STATUSES,
	RETURN_POINT_MATERIALS,
	RETURN_POINT_REFUND_METHODS,
} from "@shared/api/responses";

import type {
	ReturnPoint,
	ReturnPointType,
	ReturnPointStatus,
	ReturnPointMaterial,
	ReturnPointRefundMethod,
} from "@shared/api/responses";

export const DISTANCE_MAX_KM = 100;

export const DEFAULT_PAGE_SIZE = 20;

const MAX_PAGE_SIZE = 100;

export const PARAM = {
	query: "query",
	status: "status",
	type: "type",
	material: "material",
	refundMethod: "refund_method",
	openNow: "open",
	minDistance: "min_distance",
	maxDistance: "max_distance",
	lat: "lat",
	lon: "lon",
	page: "page",
	pageSize: "page_size",
} as const;

export interface ListFilters {
	query: string;
	statuses: ReturnPointStatus[];
	types: ReturnPointType[];
	materials: ReturnPointMaterial[];
	refundMethods: ReturnPointRefundMethod[];
	openNow: boolean;
	minDistanceKm: number;
	maxDistanceKm: number;
	lat: number | undefined;
	lon: number | undefined;
	page: number;
	pageSize: number;
}

const readEnum = <TValue extends string>(
	params: URLSearchParams,
	key: string,
	allowed: readonly TValue[],
): TValue[] => {
	if (!params.has(key)) return [...allowed];

	const wanted = new Set(params.getAll(key));

	return allowed.filter((value) => wanted.has(value));
}

const readNumber = (
	params: URLSearchParams,
	key: string,
	fallback: number,
	min: number,
	max: number,
): number => {
	const raw = params.get(key);

	if (raw === null) return fallback;

	const value = Number(raw);

	if (!Number.isFinite(value)) return fallback;

	return Math.min(Math.max(value, min), max);
}

const readCoordinate = (
	params: URLSearchParams,
	key: string,
	limit: number,
): number | undefined => {
	const raw = params.get(key);

	if (raw === null) return undefined;

	const value = Number(raw);

	return Number.isFinite(value) && Math.abs(value) <= limit ? value : undefined;
}

export const parseListFilters = (params: URLSearchParams): ListFilters => {
	const minDistanceKm = readNumber(params, PARAM.minDistance, 0, 0, DISTANCE_MAX_KM);
	const maxDistanceKm = readNumber(params, PARAM.maxDistance, DISTANCE_MAX_KM, 0, DISTANCE_MAX_KM);

	const lat = readCoordinate(params, PARAM.lat, 90);
	const lon = readCoordinate(params, PARAM.lon, 180);

	return {
		query: (params.get(PARAM.query) ?? "").trim(),
		statuses: readEnum(params, PARAM.status, RETURN_POINT_STATUSES),
		types: readEnum(params, PARAM.type, RETURN_POINT_TYPES),
		materials: readEnum(params, PARAM.material, RETURN_POINT_MATERIALS),
		refundMethods: readEnum(params, PARAM.refundMethod, RETURN_POINT_REFUND_METHODS),
		openNow: params.get(PARAM.openNow) === "1",
		minDistanceKm: Math.min(minDistanceKm, maxDistanceKm),
		maxDistanceKm: Math.max(minDistanceKm, maxDistanceKm),
		lat: lat !== undefined && lon !== undefined ? lat : undefined,
		lon: lat !== undefined && lon !== undefined ? lon : undefined,
		page: Math.trunc(readNumber(params, PARAM.page, 1, 1, Number.MAX_SAFE_INTEGER)),
		pageSize: Math.trunc(readNumber(params, PARAM.pageSize, DEFAULT_PAGE_SIZE, 1, MAX_PAGE_SIZE)),
	};
}

export const activeFilterGroupCount = (filters: ListFilters): number =>
	[
		filters.statuses.length < RETURN_POINT_STATUSES.length,
		filters.types.length < RETURN_POINT_TYPES.length,
		filters.materials.length < RETURN_POINT_MATERIALS.length,
		filters.refundMethods.length < RETURN_POINT_REFUND_METHODS.length,
		filters.openNow,
		filters.lat !== undefined
			&& (filters.minDistanceKm > 0 || filters.maxDistanceKm < DISTANCE_MAX_KM),
	].filter(Boolean).length;

export const hasActiveFilters = (filters: ListFilters): boolean =>
	filters.query.length > 0 || activeFilterGroupCount(filters) > 0;

const acceptsMaterial = (point: ReturnPoint, material: ReturnPointMaterial): boolean =>
	material === "pet" ? point.pet : material === "can" ? point.can : point.glass;

export const matchesFilters = (point: ReturnPoint, filters: ListFilters): boolean => {
	if (!filters.statuses.includes(point.current_status)) return false;

	if (!filters.types.includes(point.type)) return false;

	if (!filters.materials.some((material) => acceptsMaterial(point, material))) return false;

	if (
		filters.refundMethods.length < RETURN_POINT_REFUND_METHODS.length
		&& !point.refund_method.some((method) => filters.refundMethods.includes(method))
	) {
		return false;
	}

	if (filters.openNow && !isOpenNow(point.opening_hours)) return false;

	return true;
}
