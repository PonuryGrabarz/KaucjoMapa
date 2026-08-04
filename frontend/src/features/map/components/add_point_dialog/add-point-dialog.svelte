<template src="./add-point-dialog.svelte.html"></template>

<style lang="sass" src="./add-point-dialog.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";
	import Dialog from "@shared/components/dialog/dialog.svelte";
	import Checkbox from "@shared/components/checkbox/checkbox.svelte";
	import TextField from "@shared/components/text_field/text-field.svelte";
	import ColorField from "@shared/components/color_field/color-field.svelte";
	import PointMarker from "@shared/components/point_marker/point-marker.svelte";
	import SelectField from "@shared/components/select_field/select-field.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { actions } from "astro:actions";
	import { events } from "@shared/events/app-events";
	import { returnPoints } from "../../stores/return-points.svelte";
	import { SubmitState } from "@shared/helpers/submit-state.svelte";
	import { RETURN_POINT_TYPES, RETURN_POINT_REFUND_METHODS } from "@shared/api/responses";

	import {
		POINT_STATUS_LABELS,
		RETURN_POINT_TYPE_LABELS,
		RETURN_POINT_REFUND_METHOD_LABELS,
	} from "@shared/helpers/return-point-labels";

	import type {
		ReturnPointType,
		ReturnPointStatus,
		ReturnPointRefundMethod,
	} from "@shared/api/responses";

	interface PickedLocation {
		latitude: number;
		longitude: number;
	}

	const DESCRIPTION_MAX_LENGTH = 500;

	const DEFAULT_BACKGROUND = "#193f2a";
	const DEFAULT_TEXT_COLOR = "#ffffff";

	const CREATABLE_STATUSES: readonly ReturnPointStatus[] = ["dostepny", "brak_informacji"];

	const TYPE_OPTIONS = RETURN_POINT_TYPES.map((value) => ({
		value,
		label: RETURN_POINT_TYPE_LABELS[value],
	}));

	const STATUS_OPTIONS = CREATABLE_STATUSES.map((value) => ({
		value,
		label: POINT_STATUS_LABELS[value],
	}));

	const componentId = $props.id();

	const formId = `${componentId}-form`;

	let open = $state(false);

	let name = $state("");
	let type = $state<ReturnPointType>("butelkomat");
	let currentStatus = $state<ReturnPointStatus>("dostepny");
	let address = $state("");
	let pet = $state(true);
	let can = $state(true);
	let glass = $state(false);
	let hoursWeekday = $state("");
	let hoursWeekend = $state("");
	let description = $state("");

	let shortName = $state("");
	let backgroundColor = $state(DEFAULT_BACKGROUND);
	let textColor = $state(DEFAULT_TEXT_COLOR);

	let pickedLocation = $state<PickedLocation | undefined>(undefined);

	const refundMethods = $state<Record<ReturnPointRefundMethod, boolean>>({
		gotowka: false,
		bon: false,
		karta: false,
	});

	const submitState = new SubmitState();

	let keepDraftOnClose = false;

	const previewShortName = $derived(shortName.trim().slice(0, 3).toUpperCase() || "??");

	const locationLabel = $derived(
		pickedLocation
			? `${pickedLocation.latitude.toFixed(5)}, ${pickedLocation.longitude.toFixed(5)}`
			: "Kliknij, aby wybrać miejsce",
	);

	const reset = () => {
		name = "";
		type = "butelkomat";
		currentStatus = "dostepny";
		address = "";
		pet = true;
		can = true;
		glass = false;
		hoursWeekday = "";
		hoursWeekend = "";
		description = "";
		shortName = "";
		backgroundColor = DEFAULT_BACKGROUND;
		textColor = DEFAULT_TEXT_COLOR;
		pickedLocation = undefined;

		for (const method of RETURN_POINT_REFUND_METHODS) refundMethods[method] = false;
	}

	const handleClose = () => {
		if (keepDraftOnClose) {
			keepDraftOnClose = false;

			return;
		}

		submitState.cancel();

		reset();
	}

	const startLocationPick = () => {
		keepDraftOnClose = true;
		open = false;

		void events.emit("map:begin_location_pick");
	}

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();

		const location = pickedLocation;

		if (!location) {
			submitState.fail("Wskaż lokalizację punktu na mapie.");

			return;
		}

		const openingHours: Record<string, string> = {};

		if (hoursWeekday.trim()) openingHours["mon-fri"] = hoursWeekday.trim();

		if (hoursWeekend.trim()) {
			openingHours["sat"] = hoursWeekend.trim();
			openingHours["sun"] = hoursWeekend.trim();
		}

		const created = await submitState.run(
			() =>
				actions.returnPoints.create({
					name,
					type,
					address,
					latitude: location.latitude,
					longitude: location.longitude,
					pet,
					can,
					glass,
					refund_method: RETURN_POINT_REFUND_METHODS.filter((method) => refundMethods[method]),
					opening_hours: openingHours,
					description,
					current_status: currentStatus,
				}),
			"Nie udało się dodać punktu.",
		);

		if (!created) return;

		returnPoints.add(created);

		open = false;
	}

	$effect(() => {
		const unsubscribeOpen = events.on("addLocation:click", () => {
			open = true;
		});

		const unsubscribePicked = events.on("map:location_picked", ({ data }) => {
			pickedLocation = { latitude: data.latitude, longitude: data.longitude };
			open = true;
		});

		return () => {
			unsubscribeOpen();
			unsubscribePicked();
		};
	});
</script>
