<template src="./report-dialog.svelte.html"></template>

<style lang="sass" src="./report-dialog.sass"></style>

<script lang="ts">
  import Button from "@shared/components/button/button.svelte";
  import Dialog from "@shared/components/dialog/dialog.svelte";
  import TextField from "@shared/components/text_field/text-field.svelte";

	import { actions } from "astro:actions";
	import { SubmitState } from "@shared/helpers/submit-state.svelte";

	import type {
		ReturnPoint,
		CommunityReportStatus,
		CreateCommunityReportResponse,
	} from "@shared/api/responses";

  interface Props {
		open?: boolean;
		point: ReturnPoint;
		onsubmitted?: (report: CreateCommunityReportResponse) => void;
	}

	interface StatusOption {
		value: CommunityReportStatus;
		label: string;
		hint: string;
	}

	const DESCRIPTION_MAX_LENGTH = 280;

	const STATUS_OPTIONS: readonly StatusOption[] = [
		{ value: "dostepny", label: "Dostępny", hint: "Punkt działa i przyjmuje opakowania." },
		{ value: "nieczynny", label: "Nieczynny", hint: "Punkt jest zamknięty albo nie działa." },
		{
			value: "brak_informacji",
			label: "Brak informacji",
			hint: "Działa częściowo lub wymaga wyjaśnienia.",
		},
	];

  let { open = $bindable(false), point, onsubmitted }: Props = $props();

	const componentId = $props.id();

	const formId = `${componentId}-form`;

	let status = $state<CommunityReportStatus>("dostepny");
	let description = $state("");

	const submitState = new SubmitState();

	const reset = () => {
		status = "dostepny";
		description = "";
	}

	const handleClose = () => {
		submitState.cancel();

		reset();
	}

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();

		const report = await submitState.run(
			() =>
				actions.communityReports.create({
					returnPointId: point.id,
					status,
					description,
				}),
			"Nie udało się dodać zgłoszenia. Sprawdź połączenie i spróbuj ponownie.",
		);

		if (!report) return;

		open = false;

		onsubmitted?.(report);
	}
</script>
