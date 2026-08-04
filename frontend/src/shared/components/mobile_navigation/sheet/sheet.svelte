<template src="./sheet.svelte.html"></template>

<style lang="sass" src="./sheet.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";
	import Dialog from "@shared/components/dialog/dialog.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { actions } from "astro:actions";

	import { unreadCountText } from "@shared/helpers/format";
	import { accountItems } from "@shared/helpers/navigation-links";
	import { notificationsStore } from "@shared/stores/notifications.svelte";

	interface Props {
		active: boolean;
		pathname: string;
		unreadCount?: number;
	}

	const { active, pathname, unreadCount = 0 }: Props = $props();

	let open = $state(false);

	const resolvedUnreadCount = $derived(
		notificationsStore.hydrated ? notificationsStore.unreadCount : unreadCount,
	);

	const handleLogout = async () => {
		const { error } = await actions.auth.logout();

		if (error) return;

		open = false;

		window.location.assign("/logowanie");
	}
</script>
