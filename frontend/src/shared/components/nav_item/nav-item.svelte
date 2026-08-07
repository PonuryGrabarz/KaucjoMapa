<template src="./nav-item.svelte.html"></template>

<style lang="sass" src="./nav-item.sass"></style>

<script lang="ts">
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { untrack } from "svelte";

	import { unreadCountText } from "@shared/helpers/format";
	import { notificationsStore } from "@shared/stores/notifications.svelte";

	import type { HTMLAnchorAttributes } from "svelte/elements";

	export type NavItemProps = Omit<HTMLAnchorAttributes, "children" | "aria-current"> & {
		href: string;
		label: string;
		icon: string;
		activeIcon?: string;
		active?: boolean;
		variant?: "sidebar" | "bottomnav";
		badgeCount?: number;
		badgeLabel?: string;
		liveBadge?: boolean;
	};

	const {
		href,
		label,
		icon,
		activeIcon,
		active = false,
		variant = "sidebar",
		badgeCount,
		badgeLabel,
		liveBadge = false,
		class: className,
		...rest
	}: NavItemProps = $props();

	if (untrack(() => liveBadge) && typeof window !== "undefined") {
		notificationsStore.seed(untrack(() => badgeCount) ?? 0);
	}

	const resolvedBadgeCount = $derived(
		liveBadge && notificationsStore.hydrated ? notificationsStore.unreadCount : badgeCount,
	);

	const classes = $derived([
		"nav-item",
		`nav-item--${variant}`,
		active && "nav-item--active",
		className,
	]);

	const resolvedIcon = $derived(active ? (activeIcon ?? icon) : icon);

	const badgeText = $derived(
		resolvedBadgeCount === undefined
			? ""
			: resolvedBadgeCount > 99
				? "99+"
				: String(resolvedBadgeCount),
	);

	const badgeDescription = $derived(
		badgeLabel ?? (resolvedBadgeCount === undefined ? "" : unreadCountText(resolvedBadgeCount)),
	);
</script>
