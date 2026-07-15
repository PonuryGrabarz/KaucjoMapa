import Styles from "./toolbar.module.sass";
import Input from "@components/input/input";
import Modal from "@components/modal/modal";
import Switch from "@components/switch/switch";
import Button from "@components/button/button";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import {
    DISTANCE_OPTIONS_KM,
    POINT_STATUSES,
    POINT_TYPES,
    activePointFiltersCount,
    applyReturnPointFilterQueryParams,
    beginLocalizationRequest,
    clearPointFilters,
    disableLocalization,
    enableLocalization,
    failLocalizationRequest,
    geolocationStatus,
    isLocalizationEnabled,
    searchQuery,
    selectedDistanceKm,
    selectedPointStatuses,
    selectedPointTypes,
    setPointStatusSelected,
    setPointTypeSelected,
    updateReturnPointFilterQueryParams,
    type GeolocationStatus,
    type PointStatus,
    type PointType,
} from "../../store";

import { useEffect, useId, useState } from "preact/hooks";

const POINT_TYPE_LABELS: Readonly<Record<PointType, string>> = {
    butelkomat: "Butelkomat",
    sklep: "Sklep",
    punkt_skupu: "Punkt skupu",
    inny: "Inny typ punktu",
};

const POINT_STATUS_LABELS: Readonly<Record<PointStatus, string>> = {
    dostepny: "Dostępny",
    zapelniony: "Zapełniony",
    nieczynny: "Nieczynny",
    awaria: "Awaria",
    brak_informacji: "Brak informacji",
    do_weryfikacji: "Do weryfikacji",
};

const GEOLOCATION_STATUS_MESSAGES: Partial<Record<GeolocationStatus, string>> = {
    denied: "Dostęp do lokalizacji został odrzucony. Zmień uprawnienia tej strony w przeglądarce.",
    unavailable: "Twoja lokalizacja jest obecnie niedostępna.",
    timeout: "Nie udało się pobrać lokalizacji w wyznaczonym czasie. Spróbuj ponownie.",
    error: "Nie udało się pobrać lokalizacji. Spróbuj ponownie.",
};

function useReturnPointFilterQueryParams(): void {
    useEffect(() => {
        let isApplyingQueryParams = false;
        let isInitializingSubscriptions = true;

        const applyQueryParams = (): void => {
            isApplyingQueryParams = true;
            applyReturnPointFilterQueryParams(new URLSearchParams(window.location.search));
            isApplyingQueryParams = false;
        };

        const updateUrl = (): void => {
            if (isApplyingQueryParams || isInitializingSubscriptions) return;

            const url = new URL(window.location.href);
            updateReturnPointFilterQueryParams(url.searchParams);

            const nextQuery = url.searchParams.toString();
            const nextUrl = `${url.pathname}${nextQuery ? `?${nextQuery}` : ""}${url.hash}`;
            const currentUrl = `${window.location.pathname}${window.location.search}${window.location.hash}`;

            if (nextUrl !== currentUrl) {
                window.history.replaceState(window.history.state, "", nextUrl);
            }
        };

        const handlePopState = (): void => {
            applyQueryParams();
            updateUrl();
        };

        applyQueryParams();

        const unsubscribers = [
            selectedDistanceKm.subscribe(updateUrl),
            selectedPointTypes.subscribe(updateUrl),
            selectedPointStatuses.subscribe(updateUrl),
        ];

        isInitializingSubscriptions = false;
        updateUrl();
        window.addEventListener("popstate", handlePopState);

        return () => {
            unsubscribers.forEach((unsubscribe) => unsubscribe());
            window.removeEventListener("popstate", handlePopState);
        };
    }, []);
}

const requestCurrentLocation = (): void => {
    if (typeof navigator === "undefined" || !navigator.geolocation) {
        failLocalizationRequest("unavailable");
        return;
    }

    beginLocalizationRequest();

    navigator.geolocation.getCurrentPosition(
        ({ coords }) => {
            enableLocalization({
                latitude: coords.latitude,
                longitude: coords.longitude,
                accuracy: coords.accuracy,
            });
        },
        (error) => {
            if (error.code === error.PERMISSION_DENIED) {
                failLocalizationRequest("denied");
                return;
            }

            if (error.code === error.POSITION_UNAVAILABLE) {
                failLocalizationRequest("unavailable");
                return;
            }

            if (error.code === error.TIMEOUT) {
                failLocalizationRequest("timeout");
                return;
            }

            failLocalizationRequest("error");
        },
        {
            enableHighAccuracy: true,
            timeout: 10_000,
            maximumAge: 60_000,
        },
    );
};

type ModalBodyProps = {
    descriptionId: string;
    onDone: () => void;
};

function DistanceFilterBody({ descriptionId, onDone }: ModalBodyProps) {
    const selectedDistance = selectedDistanceKm.value;

    return (
        <div class={Styles["toolbar-modal"]} data-distance-filter-body="">
            <p id={descriptionId} class={Styles["toolbar-modal__description"]}>
                Wybierz maksymalną odległość od swojej lokalizacji.
            </p>

            <fieldset class={Styles["toolbar-modal__fieldset"]}>
                <legend class={Styles["toolbar-modal__legend"]}>Maksymalna odległość</legend>

                <div class={Styles["toolbar-modal__distance-options"]}>
                    {DISTANCE_OPTIONS_KM.map((distance) => (
                        <label
                            key={distance}
                            class={Styles["toolbar-modal__option"]}
                            data-selected={selectedDistance === distance ? "true" : undefined}
                        >
                            <input
                                type="radio"
                                name="distance"
                                value={distance}
                                checked={selectedDistance === distance}
                                onChange={(event) => {
                                    if (event.currentTarget.checked) selectedDistanceKm.value = distance;
                                }}
                            />

                            <span>{distance} km</span>
                        </label>
                    ))}
                </div>
            </fieldset>

            <div class={Styles["toolbar-modal__actions"]}>
                <Button
                    class={Styles["toolbar-modal__done"]}
                    variant="contained"
                    onClick={onDone}
                >
                    Gotowe
                </Button>
            </div>
        </div>
    );
}

function PointFiltersBody({ descriptionId, onDone }: ModalBodyProps) {
    const pointTypes = selectedPointTypes.value;
    const pointStatuses = selectedPointStatuses.value;
    const activeFiltersCount = activePointFiltersCount.value;

    return (
        <div class={Styles["toolbar-modal"]} data-point-filters-body="">
            <p id={descriptionId} class={Styles["toolbar-modal__description"]}>
                Zaznacz typy i statusy punktów, które chcesz zobaczyć. Brak zaznaczenia
                oznacza wszystkie możliwości.
            </p>

            <fieldset class={Styles["toolbar-modal__fieldset"]}>
                <legend class={Styles["toolbar-modal__legend"]}>Typ punktu</legend>

                <div class={Styles["toolbar-modal__filter-options"]}>
                    {POINT_TYPES.map((value) => {
                        const label = POINT_TYPE_LABELS[value];
                        const isSelected = pointTypes.includes(value);

                        return (
                            <label
                                key={value}
                                class={Styles["toolbar-modal__option"]}
                                data-selected={isSelected ? "true" : undefined}
                            >
                                <input
                                    type="checkbox"
                                    name="pointType"
                                    value={value}
                                    checked={isSelected}
                                    onChange={(event) => {
                                        setPointTypeSelected(value, event.currentTarget.checked);
                                    }}
                                />

                                <span>{label}</span>
                            </label>
                        );
                    })}
                </div>
            </fieldset>

            <fieldset class={Styles["toolbar-modal__fieldset"]}>
                <legend class={Styles["toolbar-modal__legend"]}>Status punktu</legend>

                <div class={Styles["toolbar-modal__filter-options"]}>
                    {POINT_STATUSES.map((value) => {
                        const label = POINT_STATUS_LABELS[value];
                        const isSelected = pointStatuses.includes(value);

                        return (
                            <label
                                key={value}
                                class={Styles["toolbar-modal__option"]}
                                data-selected={isSelected ? "true" : undefined}
                            >
                                <input
                                    type="checkbox"
                                    name="pointStatus"
                                    value={value}
                                    checked={isSelected}
                                    onChange={(event) => {
                                        setPointStatusSelected(value, event.currentTarget.checked);
                                    }}
                                />

                                <span>{label}</span>
                            </label>
                        );
                    })}
                </div>
            </fieldset>

            <p class={Styles["toolbar-modal__selection-summary"]} aria-live="polite">
                {activeFiltersCount === 0
                    ? "Wyświetlane są wszystkie typy i statusy punktów."
                    : `Aktywne filtry: ${activeFiltersCount}.`}
            </p>

            <div class={Styles["toolbar-modal__actions"]}>
                <Button
                    class={Styles["toolbar-modal__clear"]}
                    variant="text"
                    disabled={activeFiltersCount === 0}
                    onClick={clearPointFilters}
                >
                    Wyczyść
                </Button>

                <Button
                    class={Styles["toolbar-modal__done"]}
                    variant="contained"
                    onClick={onDone}
                >
                    Gotowe
                </Button>
            </div>
        </div>
    );
}

export default function Toolbar() {
    const [isDistanceFilterOpen, setIsDistanceFilterOpen] = useState(false);
    const [isFiltersOpen, setIsFiltersOpen] = useState(false);

    useReturnPointFilterQueryParams();

    const searchInputId = useId();
    const distanceDialogId = useId();
    const distanceDescriptionId = useId();
    const filtersDialogId = useId();
    const filtersDescriptionId = useId();
    const locationStatusId = useId();

    const selectedDistance = selectedDistanceKm.value;
    const activeFiltersCount = activePointFiltersCount.value;
    const localizationEnabled = isLocalizationEnabled.value;
    const currentGeolocationStatus = geolocationStatus.value;
    const isRequestingLocation = currentGeolocationStatus === "requesting";
    const locationStatusMessage = GEOLOCATION_STATUS_MESSAGES[currentGeolocationStatus];

    const handleLocationToggle = (checked: boolean): void => {
        if (!checked) {
            disableLocalization();
            return;
        }

        requestCurrentLocation();
    };

    return (
        <section
            class={Styles["toolbar"]}
            aria-label="Narzędzia mapy punktów zwrotu"
            data-active-filters={activeFiltersCount || undefined}
            data-geolocation-status={currentGeolocationStatus}
        >
            <form
                class={Styles["toolbar__form"]}
                role="search"
                aria-label="Wyszukiwanie i filtrowanie punktów zwrotu"
                onSubmit={(event) => event.preventDefault()}
            >
                <div class={Styles["toolbar__search"]}>
                    <label class={Styles["toolbar__visually-hidden"]} for={searchInputId}>
                        Wyszukaj punkt zwrotu
                    </label>

                    <Input
                        id={searchInputId}
                        name="query"
                        type="search"
                        value={searchQuery}
                        fullWidth
                        data-return-points-search=""
                        autoComplete="off"
                        enterKeyHint="search"
                        startAdornment="material-symbols:search"
                        placeholder="Nazwa, adres lub kod pocztowy"
                        onInput={(event) => {
                            searchQuery.value = event.currentTarget.value;
                        }}
                    />
                </div>

                <div
                    class={Styles["toolbar__actions"]}
                    role="group"
                    aria-label="Filtry mapy"
                >
                    <Button
                        class={Styles["toolbar__actions__button"]}
                        variant="contained"
                        startIcon="mdi:map-marker-radius-outline"
                        aria-controls={distanceDialogId}
                        aria-expanded={isDistanceFilterOpen}
                        aria-haspopup="dialog"
                        aria-label={`Zakres odległości: ${selectedDistance} km`}
                        title={!localizationEnabled
                            ? "Włącz lokalizację, aby zmienić zakres odległości"
                            : undefined}
                        disabled={!localizationEnabled}
                        data-distance-filter-trigger=""
                        onClick={() => setIsDistanceFilterOpen(true)}
                    >
                        {selectedDistance} km
                    </Button>

                    <Button
                        class={Styles["toolbar__actions__button"]}
                        variant="contained"
                        startIcon="mdi:tune-variant"
                        aria-controls={filtersDialogId}
                        aria-expanded={isFiltersOpen}
                        aria-haspopup="dialog"
                        aria-label={activeFiltersCount > 0
                            ? `Filtry punktów zwrotu, aktywne: ${activeFiltersCount}`
                            : "Filtry punktów zwrotu"}
                        data-point-filters-trigger=""
                        onClick={() => setIsFiltersOpen(true)}
                    >
                        <span>Filtry</span>

                        {activeFiltersCount > 0 && (
                            <span
                                class={Styles["toolbar__filter-count"]}
                                aria-hidden="true"
                                data-filter-count=""
                            >
                                {activeFiltersCount}
                            </span>
                        )}
                    </Button>

                    <Switch
                        class={Styles["toolbar__actions__location"]}
                        checked={localizationEnabled}
                        disabled={isRequestingLocation}
                        aria-describedby={locationStatusMessage ? locationStatusId : undefined}
                        onChange={handleLocationToggle}
                    >
                        <span class={Styles["toolbar__actions__location__content"]}>
                            <IconifyIcon icon="mdi:crosshairs-gps" aria-hidden="true" />

                            <span class={Styles["toolbar__actions__location__text"]}>
                                {isRequestingLocation
                                    ? "Pobieranie lokalizacji…"
                                    : `Lokalizacja ${localizationEnabled ? "włączona" : "wyłączona"}`}
                            </span>
                        </span>
                    </Switch>
                </div>
            </form>

            {locationStatusMessage && (
                <p
                    id={locationStatusId}
                    class={Styles["toolbar__location-status"]}
                    role="status"
                    aria-live="polite"
                    data-status={currentGeolocationStatus}
                >
                    {locationStatusMessage}
                </p>
            )}

            <Modal
                id={distanceDialogId}
                class={Styles["toolbar__modal"]}
                bodyClass={Styles["toolbar__modal-body"]}
                open={isDistanceFilterOpen}
                onClose={() => setIsDistanceFilterOpen(false)}
                title="Zakres odległości"
                closeLabel="Zamknij wybór zakresu odległości"
                aria-describedby={distanceDescriptionId}
            >
                <DistanceFilterBody
                    descriptionId={distanceDescriptionId}
                    onDone={() => setIsDistanceFilterOpen(false)}
                />
            </Modal>

            <Modal
                id={filtersDialogId}
                class={Styles["toolbar__modal"]}
                bodyClass={Styles["toolbar__modal-body"]}
                open={isFiltersOpen}
                onClose={() => setIsFiltersOpen(false)}
                title="Filtry punktów zwrotu"
                closeLabel="Zamknij filtry punktów zwrotu"
                aria-describedby={filtersDescriptionId}
            >
                <PointFiltersBody
                    descriptionId={filtersDescriptionId}
                    onDone={() => setIsFiltersOpen(false)}
                />
            </Modal>
        </section>
    );
}