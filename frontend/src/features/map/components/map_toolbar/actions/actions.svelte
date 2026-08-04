<template src="./actions.svelte.html"></template>

<style lang="sass" src="./actions.sass"></style>

<script lang="ts">
  import Button from "@shared/components/button/button.svelte";
  import IconifyIcon from "@shared/components/iconify_icon/iconify-icon.svelte";
  import FilterPanel from "@shared/components/filter_panel/filter-panel.svelte";

  import { FILTER_GROUPS } from "./filter-options";
  import { events } from "@shared/events/app-events";
  import { mapState, RADIUS_MAX_KM } from "../../../stores/map-state.svelte";
  import { geolocationAccess } from "@shared/helpers/geolocation-access.svelte";

  const locationButtonLabel = $derived(
    mapState.locationState === "requesting"
      ? "Ustalam pozycję…"
      : mapState.locationState === "active"
        ? "Wycentruj ponownie"
        : "Moja lokalizacja",
  );

  const requestLocation = () => {
    void events.emit("map:request_location");
  }

  $effect(() => {
    const position = geolocationAccess.position;

    if (position && !mapState.referenceCenter)
      mapState.referenceCenter = [position.lon, position.lat];
  });
</script>
