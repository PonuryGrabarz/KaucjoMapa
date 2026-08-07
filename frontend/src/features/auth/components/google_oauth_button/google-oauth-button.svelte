<template src="./google-oauth-button.svelte.html"></template>

<style lang="sass" src="./google-oauth-button.sass"></style>

<script lang="ts">
	import type { HTMLAnchorAttributes } from "svelte/elements";

	export type GoogleOAuthButtonProps = Omit<HTMLAnchorAttributes, "children"> & {
		label?: string;
		disabled?: boolean;
	};

	const {
		label = "Kontynuuj z Google",
		disabled = false,
		href,
		rel = "nofollow",
		...rest
	}: GoogleOAuthButtonProps = $props();
</script>
