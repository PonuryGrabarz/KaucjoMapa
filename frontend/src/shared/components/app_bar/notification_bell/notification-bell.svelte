<template src="./notification-bell.svelte.html"></template>

<style lang="sass" src="./notification-bell.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";
	import Popover from "@shared/components/popover/popover.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { untrack } from "svelte";
	import { actions } from "astro:actions";

	import {
		dateTimeLabel,
		relativeTimeLabel,
		unreadCountText,
	} from "@shared/helpers/format";

	import { notificationsStore } from "@shared/stores/notifications.svelte";

	import type { Notification } from "@shared/api/responses/index";

	interface Props {
		initialNotifications: Notification[];
	}

	const { initialNotifications }: Props = $props();

	const POPOVER_ID = "notification-bell-popover";

	let notifications = $state(untrack(() => initialNotifications));
	let open = $state(false);

	let triggerElement = $state<HTMLElement | null>(null);

	const localUnreadCount = $derived(notifications.filter((n) => !n.seen).length);

	const unreadCount = $derived(
		notificationsStore.hydrated ? notificationsStore.unreadCount : localUnreadCount,
	);

	const badgeLabel = $derived(unreadCount > 99 ? "99+" : String(unreadCount));

	const handleMarkSeen = async (notificationId: string) => {
		const formData = new FormData();

		formData.set("notificationId", notificationId);

		const { error } = await actions.notifications.markSeen(formData);

		if (error) return;

		const target = notifications.find((notification) => notification.id === notificationId);

		if (!target || target.seen) return;

		target.seen = true;

		notificationsStore.markOneSeen();
	}
</script>
