import { batch, computed, signal } from "@preact/signals";

export const DISTANCE_OPTIONS_KM = [1, 2, 5, 10, 25, 50] as const;
export const DEFAULT_DISTANCE_KM = 5;

export const RETURN_POINT_QUERY_PARAMS = {
	distance: "distance",
	pointType: "type",
	pointStatus: "status",
} as const;

export const POINT_TYPES = [
	"butelkomat",
	"sklep",
	"punkt_skupu",
	"inny",
] as const;

export const POINT_STATUSES = [
	"dostepny",
	"zapelniony",
	"nieczynny",
	"awaria",
	"brak_informacji",
	"do_weryfikacji",
] as const;

export type DistanceKm = (typeof DISTANCE_OPTIONS_KM)[number];
export type PointType = (typeof POINT_TYPES)[number];
export type PointStatus = (typeof POINT_STATUSES)[number];
export type GeolocationStatus =
	| "idle"
	| "requesting"
	| "enabled"
	| "denied"
	| "unavailable"
	| "timeout"
	| "error";
export type GeolocationFailureStatus = Exclude<
	GeolocationStatus,
	"idle" | "requesting" | "enabled"
>;

export type UserLocation = Readonly<{
	latitude: number;
	longitude: number;
	accuracy: number;
}>;

export const searchQuery = signal("");
export const selectedDistanceKm = signal<DistanceKm>(DEFAULT_DISTANCE_KM);
export const selectedPointTypes = signal<readonly PointType[]>([]);
export const selectedPointStatuses = signal<readonly PointStatus[]>([]);
export const isLocalizationEnabled = signal(false);
export const geolocationStatus = signal<GeolocationStatus>("idle");
export const userLocation = signal<UserLocation | null>(null);

export const activePointFiltersCount = computed(() => {
	return selectedPointTypes.value.length + selectedPointStatuses.value.length;
});

const updateSelection = <Value extends string>(
	selection: readonly Value[],
	value: Value,
	selected: boolean,
): readonly Value[] => {
	if (selected) {
		return selection.includes(value) ? selection : [...selection, value];
	}

	return selection.filter((selectedValue) => selectedValue !== value);
};

const isDistanceKm = (value: number): value is DistanceKm => {
	return DISTANCE_OPTIONS_KM.some((distance) => distance === value);
};

export const setPointTypeSelected = (pointType: PointType, selected: boolean): void => {
	selectedPointTypes.value = updateSelection(selectedPointTypes.value, pointType, selected);
};

export const setPointStatusSelected = (pointStatus: PointStatus, selected: boolean): void => {
	selectedPointStatuses.value = updateSelection(
		selectedPointStatuses.value,
		pointStatus,
		selected,
	);
};

export const clearPointFilters = (): void => {
	batch(() => {
		selectedPointTypes.value = [];
		selectedPointStatuses.value = [];
	});
};

export const applyReturnPointFilterQueryParams = (searchParams: URLSearchParams): void => {
	const parsedDistance = Number(searchParams.get(RETURN_POINT_QUERY_PARAMS.distance));
	const queryPointTypes = searchParams.getAll(RETURN_POINT_QUERY_PARAMS.pointType);
	const queryPointStatuses = searchParams.getAll(RETURN_POINT_QUERY_PARAMS.pointStatus);

	batch(() => {
		selectedDistanceKm.value = isDistanceKm(parsedDistance)
			? parsedDistance
			: DEFAULT_DISTANCE_KM;
		selectedPointTypes.value = POINT_TYPES.filter((pointType) => {
			return queryPointTypes.includes(pointType);
		});
		selectedPointStatuses.value = POINT_STATUSES.filter((pointStatus) => {
			return queryPointStatuses.includes(pointStatus);
		});
	});
};

export const updateReturnPointFilterQueryParams = (
	searchParams: URLSearchParams,
): URLSearchParams => {
	searchParams.set(
		RETURN_POINT_QUERY_PARAMS.distance,
		String(selectedDistanceKm.value),
	);
	searchParams.delete(RETURN_POINT_QUERY_PARAMS.pointType);
	searchParams.delete(RETURN_POINT_QUERY_PARAMS.pointStatus);

	POINT_TYPES.filter((pointType) => selectedPointTypes.value.includes(pointType)).forEach(
		(pointType) => searchParams.append(RETURN_POINT_QUERY_PARAMS.pointType, pointType),
	);
	POINT_STATUSES.filter((pointStatus) => {
		return selectedPointStatuses.value.includes(pointStatus);
	}).forEach((pointStatus) => {
		searchParams.append(RETURN_POINT_QUERY_PARAMS.pointStatus, pointStatus);
	});

	return searchParams;
};

export const beginLocalizationRequest = (): void => {
	batch(() => {
		isLocalizationEnabled.value = false;
		geolocationStatus.value = "requesting";
	});
};

export const enableLocalization = (location: UserLocation): void => {
	batch(() => {
		userLocation.value = location;
		isLocalizationEnabled.value = true;
		geolocationStatus.value = "enabled";
	});
};

export const disableLocalization = (): void => {
	batch(() => {
		userLocation.value = null;
		isLocalizationEnabled.value = false;
		geolocationStatus.value = "idle";
	});
};

export const failLocalizationRequest = (status: GeolocationFailureStatus): void => {
	batch(() => {
		userLocation.value = null;
		isLocalizationEnabled.value = false;
		geolocationStatus.value = status;
	});
};