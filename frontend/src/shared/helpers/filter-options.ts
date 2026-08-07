import type { FilterOption } from "@shared/components/filter_group/filter-group-types";

export const toFilterOptions = <TValue extends string>(
	values: readonly TValue[],
	labels: Record<TValue, string>,
): FilterOption[] => values.map((value) => ({ value, label: labels[value] }));
