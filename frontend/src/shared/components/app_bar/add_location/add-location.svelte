<template src="./add-location.svelte.html"></template>

<style lang="sass" src="./add-location.sass"></style>

<script lang="ts">
	import Button from "@shared/components/button/button.svelte";
	import IconButton from "@shared/components/icon_button/icon-button.svelte";
	import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";

	import { events } from "@shared/events/app-events";

	const ADD_LOCATION_ICON = "material-symbols:add-location-alt-outline-rounded";

	const handleAddLocation = () => {
		void events.emit("addLocation:click");
	};
</script>
