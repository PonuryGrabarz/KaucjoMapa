<template src="./avatar.svelte.html"></template>

<style lang="sass" src="./avatar.sass"></style>

<script lang="ts">
	interface Props {
		username: string;
		avatarUrl?: string | null;
		size?: number;
		alt?: string;
	}

	let { username, avatarUrl, size = 30, alt = "" }: Props = $props();

	let failedUrl = $state<string | null>(null);

	const showImage = $derived(Boolean(avatarUrl) && avatarUrl !== failedUrl);

	const initials = $derived(
		username
			.trim()
			.split(/\s+/)
			.slice(0, 2)
			.map((part) => Array.from(part)[0]?.toUpperCase() ?? "")
			.join("") || "?",
	);
</script>
