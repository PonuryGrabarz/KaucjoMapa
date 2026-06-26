import type { ReturnPointsMap } from "../../components/return_points_map/return-points-map.ts";

export class HomePageView extends HTMLElement {
  private readonly _returnPointsMap: ReturnPointsMap;

  constructor() {
    super();

    this._returnPointsMap = this.querySelector<ReturnPointsMap>("return-points-map")!;
  }

  connectedCallback() {
    this._returnPointsMap.mapFab.disabled = true;
  }

  disconnectedCallback() {}
}

// Register the custom element if it hasn't been registered yet
if (!customElements.get("home-page-view")) customElements.define("home-page-view", HomePageView);
