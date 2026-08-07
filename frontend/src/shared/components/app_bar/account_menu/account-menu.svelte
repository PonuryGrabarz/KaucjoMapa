<template src="./account-menu.svelte.html"></template>

<style lang="sass" src="./account-menu.sass"></style>

<script lang="ts">
	import Menu from "@shared/components/menu/menu.svelte";
	import Button from "@shared/components/button/button.svelte";
	import MenuItem from "@shared/components/menu/menu_item/menu-item.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { actions } from "astro:actions";

	import type { Snippet } from "svelte";

	interface Props {
		avatar?: Snippet;
	}

	let { avatar }: Props = $props();

	const MENU_ID = "account-menu-popover";

	let open = $state(false);

	let triggerElement = $state<HTMLElement | null>(null);

	const handleLogout = async () => {
		const { error } = await actions.auth.logout();

		if (error) return;

		window.location.assign("/logowanie");
	}
</script>
