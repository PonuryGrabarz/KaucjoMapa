export interface SearchSuggestion {
	key: string;
	label: string;
	secondary?: string;
	icon?: string;
}

export interface SearchSuggestionGroup {
	keyPrefix: string;
	values: readonly string[];
	secondary?: string;
	icon?: string;
}
