<template src="./point-detail-dialog.svelte.html"></template>

<style lang="sass" src="./point-detail-dialog.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";
	import Dialog from "@shared/components/dialog/dialog.svelte";
	import ReportDialog from "./report_dialog/report-dialog.svelte";
	import StatusPill from "@shared/components/status_pill/status-pill.svelte";
	import PointMarker from "@shared/components/point_marker/point-marker.svelte";

	import { actions } from "astro:actions";
	import { dateTimeLabel } from "@shared/helpers/format";
	import { mapState } from "../../stores/map-state.svelte";
	import { returnPoints } from "../../stores/return-points.svelte";
	import { formatHourRange, isOpenNow, openingHoursEntries } from "@shared/helpers/opening-hours";

	import {
		OPENING_HOURS_LABELS,
		RETURN_POINT_TYPE_LABELS,
		RETURN_POINT_REFUND_METHOD_LABELS,
	} from "@shared/helpers/return-point-labels";

	import type {
		OpeningHours,
		CommunityReportVoteValue,
		CreateCommunityReportResponse,
	} from "@shared/api/responses";

	interface Props {
		isLoggedIn: boolean;
	}

	const { isLoggedIn }: Props = $props();

	const TRUST_STARS = [0, 1, 2, 3, 4];

	let reportDialogOpen = $state(false);

	let reports = $state<CreateCommunityReportResponse[]>([]);
	let reportsLoading = $state(false);
	let bookmarked = $state(false);
	let bookmarkPending = $state(false);
	let votedReportIds = $state<string[]>([]);
	let pendingVoteReportId = $state<string | undefined>(undefined);

	let loadToken = 0;

	const point = $derived(returnPoints.find(mapState.selectedReturnPointId));

	let open = $derived(point !== undefined && mapState.selectionSettled);

	const openNow = $derived(isOpenNow(point?.opening_hours));

	const materials = $derived(
		point
			? [
					point.pet && "Butelki PET",
					point.can && "Puszki",
					point.glass && "Butelki szklane",
				].filter((entry): entry is string => typeof entry === "string")
			: [],
	);

	const hoursEntries = $derived(openingHoursEntries(point?.opening_hours));

	const navigateHref = $derived(
		point
			? `https://www.google.com/maps/dir/?api=1&destination=${point.latitude},${point.longitude}`
			: "#",
	);

	const trustStars = (reputation: number): number =>
		Number.isFinite(reputation) ? Math.min(5, Math.max(0, Math.round(reputation))) : 0;

	const handleClose = () => {
		mapState.selectedReturnPointId = undefined;
	}

	const handleVote = async (reportId: string, vote: CommunityReportVoteValue) => {
		pendingVoteReportId = reportId;

		const { error } = await actions.communityReports.vote({ reportId, vote });

		pendingVoteReportId = undefined;

		if (!error) votedReportIds = [...votedReportIds, reportId];
	}

	const toggleBookmark = async () => {
		if (!point || bookmarkPending) return;

		bookmarkPending = true;

		const { data } = bookmarked
			? await actions.bookmarks.remove({ returnPointId: point.id })
			: await actions.bookmarks.add({ returnPointId: point.id });

		bookmarkPending = false;

		if (data) bookmarked = data.bookmarked;
	}

	const handleReportSubmitted = (report: CreateCommunityReportResponse) => {
		reports = [report, ...reports];
	}

	$effect(() => {
		const selected = point;

		const token = (loadToken += 1);

		reports = [];
		votedReportIds = [];
		bookmarked = false;

		if (!selected) return;

		reportsLoading = true;

		void actions.communityReports.list({ returnPointId: selected.id }).then(({ data }) => {
			if (token !== loadToken) return;

			reportsLoading = false;

			if (data) reports = data;
		});

		if (!isLoggedIn) return;

		void actions.bookmarks.isBookmarked({ returnPointId: selected.id }).then(({ data }) => {
			if (token !== loadToken) return;

			bookmarked = data?.bookmarked ?? false;
		});
	});
</script>
