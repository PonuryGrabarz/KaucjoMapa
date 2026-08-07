import { mapState } from "../../../stores/map-state.svelte";
import { toFilterOptions } from "@shared/helpers/filter-options";

import {
	RETURN_POINT_TYPES,
	RETURN_POINT_STATUSES,
	RETURN_POINT_MATERIALS,
	RETURN_POINT_REFUND_METHODS,
} from "@shared/api/responses";

import {
	RETURN_POINT_TYPE_LABELS,
	POINT_STATUS_LABELS,
	RETURN_POINT_MATERIAL_LABELS,
	RETURN_POINT_REFUND_METHOD_LABELS,
} from "@shared/helpers/return-point-labels";

import type { SvelteSet } from "svelte/reactivity";
import type { FilterPanelGroup } from "@shared/components/filter_panel/filter-panel-types";

const bindGroup = <TValue extends string>(
	legend: string,
	name: string,
	values: readonly TValue[],
	labels: Record<TValue, string>,
	set: SvelteSet<TValue>,
	layout: { columns?: boolean; wide?: boolean } = {},
): FilterPanelGroup => ({
	legend,
	name,
	options: toFilterOptions(values, labels),
	columns: layout.columns,
	wide: layout.wide,

	get checked(): readonly string[] {
		return values.filter((value) => set.has(value));
	},

	onToggle: (value, checked) => {
		if (checked) set.add(value as TValue);
		else set.delete(value as TValue);
	},
});

export const FILTER_GROUPS: readonly FilterPanelGroup[] = [
	bindGroup("Status", "status", RETURN_POINT_STATUSES, POINT_STATUS_LABELS, mapState.statuses),
	bindGroup("Typ punktu", "type", RETURN_POINT_TYPES, RETURN_POINT_TYPE_LABELS, mapState.types),
	bindGroup(
		"Opakowania",
		"material",
		RETURN_POINT_MATERIALS,
		RETURN_POINT_MATERIAL_LABELS,
		mapState.materials,
	),
	bindGroup(
		"Formy rozliczenia",
		"refund_method",
		RETURN_POINT_REFUND_METHODS,
		RETURN_POINT_REFUND_METHOD_LABELS,
		mapState.refundMethods,
		{ columns: true, wide: true },
	),
];
