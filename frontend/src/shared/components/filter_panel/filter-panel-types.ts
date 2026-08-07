import type { FilterOption } from "@shared/components/filter_group/filter-group-types";

export interface FilterPanelGroup {
	legend: string;
	name: string;
	options: readonly FilterOption[];
	checked?: readonly string[];
	onToggle?: (value: string, checked: boolean) => void;
	columns?: boolean;
	wide?: boolean;
}
