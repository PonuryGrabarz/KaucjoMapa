import "leaflet/dist/leaflet.css";

import LF from "leaflet";

import { MapPin } from "../map_pin/map-pin.ts";

import type { GetMapPointsResponse } from "@shared/api/responses/get-map-points-response.ts";

const DEFAULT_MAP_ZOOM = 7;
const DEFAULT_MAP_CENTER: [number, number] = [52.191097, 19.355406];

class ReturnPointsMap extends HTMLElement {
  private _isMapInitialized: boolean = false;
  private _map: LF.Map | null = null;

  async connectedCallback() {
    if (this._isMapInitialized) {
      console.debug("Map is already initialized.");
      return;
    }

    const mapContainer = this.querySelector("[data-map-container]") as HTMLElement;

    if (mapContainer) this._map = LF.map(mapContainer, { zoomControl: false });

    await this.initializeMap();
  }

  disconnectedCallback() {
    this._isMapInitialized = false;

    if (this._map) {
      this._map.remove();
      this._map = null;
    }
  }

  private async initializeMap() {
    if (!this._map) {
      console.error("Map container not found.");
      return;
    }

    this._map.setView(DEFAULT_MAP_CENTER, DEFAULT_MAP_ZOOM);

    LF.control.zoom({ position: "bottomleft" }).addTo(this._map);

    LF.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png',
      { attribution: '&copy; OpenStreetMap &copy; CARTO', maxZoom: 20, subdomains: 'abcd' }).addTo(this._map);

    const mapPoints = (window as unknown as { initialReturnPoints: GetMapPointsResponse })
      .initialReturnPoints.points;

    for (const point of mapPoints) {
      (await MapPin.CreateMarker(point)).addTo(this._map!);
    }

    this._isMapInitialized = true;
  }
}

// Register the custom element if it hasn't been registered yet
if (!customElements.get("return-points-map")) customElements.define("return-points-map", ReturnPointsMap);
