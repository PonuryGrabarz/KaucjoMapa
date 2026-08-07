import { SvelteSet } from "svelte/reactivity";
import { distanceKm } from "@shared/helpers/map";
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

export const RADIUS_MAX_KM = 100;

export type MapLocationState = "idle" | "requesting" | "active" | "denied" | "unavailable";

class MapState {
	public readonly types = new SvelteSet<ReturnPointType>(RETURN_POINT_TYPES);
	public readonly statuses = new SvelteSet<ReturnPointStatus>(RETURN_POINT_STATUSES);
	public readonly materials = new SvelteSet<ReturnPointMaterial>(RETURN_POINT_MATERIALS);
	public readonly refundMethods = new SvelteSet<ReturnPointRefundMethod>(RETURN_POINT_REFUND_METHODS);

	public openNowOnly = $state(false);
	public minRadiusKm = $state(0);
	public maxRadiusKm = $state(RADIUS_MAX_KM);

	public searchNoResults = $state(false);
	public locationState = $state<MapLocationState>("idle");
	public referenceCenter = $state<[number, number] | undefined>(undefined);

	private selectedId = $state<string | undefined>(undefined);

	public selectionSettled = $state(false);

	public get selectedReturnPointId(): string | undefined {
		return this.selectedId;
	}

	public set selectedReturnPointId(value: string | undefined) {
		if (this.selectedId === value) return;

		this.selectedId = value;
		this.selectionSettled = value === undefined;
	}

	public get radiusEnabled(): boolean {
		return this.referenceCenter !== undefined;
	}

	public get activeGroupCount(): number {
		return [
			this.statuses.size < RETURN_POINT_STATUSES.length,
			this.types.size < RETURN_POINT_TYPES.length,
			this.materials.size < RETURN_POINT_MATERIALS.length,
			this.refundMethods.size < RETURN_POINT_REFUND_METHODS.length,
			this.openNowOnly,
			this.radiusEnabled && (this.minRadiusKm > 0 || this.maxRadiusKm < RADIUS_MAX_KM),
		].filter(Boolean).length;
	}

	public setMinRadiusKm(value: number): void {
		this.minRadiusKm = Math.min(value, this.maxRadiusKm);
	}

	public setMaxRadiusKm(value: number): void {
		this.maxRadiusKm = Math.max(value, this.minRadiusKm);
	}

	public matches(point: ReturnPoint): boolean {
		if (!this.statuses.has(point.current_status) || !this.types.has(point.type)) return false;

		const acceptsAnyMaterial =
			(point.pet && this.materials.has("pet")) ||
			(point.can && this.materials.has("can")) ||
			(point.glass && this.materials.has("glass"));

		if (!acceptsAnyMaterial) return false;

		if (
			this.refundMethods.size < RETURN_POINT_REFUND_METHODS.length &&
			!point.refund_method.some((method) => this.refundMethods.has(method))
		) {
			return false;
		}

		if (this.openNowOnly && !isOpenNow(point.opening_hours)) return false;

		if (this.referenceCenter) {
			const distance = distanceKm(this.referenceCenter, [point.longitude, point.latitude]);

			if (distance < this.minRadiusKm || distance > this.maxRadiusKm) return false;
		}

		return true;
	}

	public apply(points: readonly ReturnPoint[]): ReturnPoint[] {
		return points.filter((point) => this.matches(point));
	}

	public reset(): void {
		this.types.clear();
		this.statuses.clear();
		this.materials.clear();
		this.refundMethods.clear();

		for (const type of RETURN_POINT_TYPES) this.types.add(type);
		for (const status of RETURN_POINT_STATUSES) this.statuses.add(status);
		for (const material of RETURN_POINT_MATERIALS) this.materials.add(material);
		for (const method of RETURN_POINT_REFUND_METHODS) this.refundMethods.add(method);

		this.openNowOnly = false;
		this.minRadiusKm = 0;
		this.maxRadiusKm = RADIUS_MAX_KM;
	}
}

export const mapState = new MapState();
