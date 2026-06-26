import LF from "leaflet";

import { TEMPLATE_ID } from "./map-pin.constants";

import type { MapPoint } from "@shared/api/responses/get-map-points-response";

export const MAP_PIN_CLICK_EVENT = "map-pin-click";

export type MapPinClickEvent = CustomEvent<{ returnPointId: string }>;

export class MapPin extends HTMLElement {
    private static readonly FOCUS_CLASS = "focus";

    private static _focusedMapPin: MapPin | null = null;

    private readonly _contentElement: HTMLDivElement;
    private readonly _shortNameElement: HTMLSpanElement;
    private readonly _leafletMarker: HTMLDivElement;

    private _returnPointId: string | null = null;
    private _shortName: string | null = null;
    private _bgColor: string | null = null;
    private _textColor: string | null = null;

    private _isHovered = false;
    private _isTouchActive = false;
    private _isContentFocused = false;
    private _defaultZIndex: string | null = null;

    private readonly _boundHandleClick = this.handleOnClick.bind(this);
    private readonly _boundHandleContentElementInteraction = this.handleContentElementInteraction.bind(this);
    private readonly _boundHandleDocumentPointerDown = this.handleDocumentPointerDown.bind(this);

    constructor() {
      super();

      this._contentElement = this.firstChild as HTMLDivElement;
      this._shortNameElement = this.querySelector("span") as HTMLSpanElement;
      this._leafletMarker = this.parentElement as HTMLDivElement;

      this._leafletMarker.removeAttribute("tabindex");
      this._leafletMarker.style.border = "none";
      this._leafletMarker.style.background = "none";

      this.updateAttributes();
    }

    connectedCallback() {
      this.addEventListener("click", this._boundHandleClick);

      this._contentElement.addEventListener("pointerenter", this._boundHandleContentElementInteraction);
      this._contentElement.addEventListener("pointerleave", this._boundHandleContentElementInteraction);
      this._contentElement.addEventListener("pointerdown", this._boundHandleContentElementInteraction);

      document.addEventListener("pointerdown", this._boundHandleDocumentPointerDown);
    }

    disconnectedCallback() {
      this.updateAttributes(undefined, null);

      this.removeEventListener("click", this._boundHandleClick);

      this._contentElement.removeEventListener("pointerenter", this._boundHandleContentElementInteraction);
      this._contentElement.removeEventListener("pointerleave", this._boundHandleContentElementInteraction);
      this._contentElement.removeEventListener("pointerdown", this._boundHandleContentElementInteraction);

      document.removeEventListener("pointerdown", this._boundHandleDocumentPointerDown);

      this.blurContentElement();
    }

    attributeChangedCallback(name: string, _oldValue: string | null, newValue: string | null) {
      this.updateAttributes(name, newValue);
    }

    private updateAttributes(name?: string, newValue?: string | null) {
      if (!name) {
        if (!newValue) {
          this._returnPointId = this.getAttribute("data-return-point-id");
          this._shortName = this.getAttribute("data-short-name");
          this._bgColor = this.getAttribute("data-bg-color");
          this._textColor = this.getAttribute("data-text-color");
        } else {
          this._returnPointId = newValue ?? null;
          this._shortName = newValue ?? null;
          this._bgColor = newValue ?? null;
          this._textColor = newValue ?? null;
        }

        this.applyBgColor();
        this.applyTextColor();
        this.applyShortName();

        return;
      }

      if (name === "data-return-point-id") {
        this._returnPointId = newValue ?? null;
      } else if (name === "data-short-name") {
        this._shortName = newValue ?? null;
        this.applyShortName();
      } else if (name === "data-bg-color") {
        this._bgColor = newValue ?? null;
        this.applyBgColor();
      } else if (name === "data-text-color") {
        this._textColor = newValue ?? null;
        this.applyTextColor();
      }
    }

    private applyBgColor() {
      this.style.setProperty("--map-pin-bg-color", this._bgColor ?? "var(--green-500)");
    }

    private applyTextColor() {
      this.style.setProperty("--map-pin-text-color", this._textColor ?? "var(--paper-50)");
    }

    private applyShortName() {
      if (this._shortNameElement) {
        this._shortNameElement.textContent = this._shortName ?? "";
      }
    }

    private handleContentElementInteraction(event: PointerEvent) {
      if (event.type === "pointerenter" && event.pointerType !== "touch") {
        this._isHovered = true;
      } else if (event.type === "pointerleave" && event.pointerType !== "touch") {
        this._isHovered = false;
      } else if (event.type === "pointerdown") {
        this.focusContentElement();

        if (event.pointerType === "touch") this._isTouchActive = true;
      }

      this.updateMarkerZIndex();
    }

    private handleDocumentPointerDown(event: PointerEvent) {
      if (event.pointerType !== "touch") return;

      const target = event.target as Node | null;

      if (target && this.contains(target)) return;

      this._isTouchActive = false;
      this.updateMarkerZIndex();
    }

    private focusContentElement() {
      if (MapPin._focusedMapPin === this) return;

      MapPin._focusedMapPin?.blurContentElement();

      this._isContentFocused = true;
      this._contentElement.classList.add(MapPin.FOCUS_CLASS);

      MapPin._focusedMapPin = this;
      this.updateMarkerZIndex();
    }

    private blurContentElement() {
      this._isContentFocused = false;
      this._contentElement.classList.remove(MapPin.FOCUS_CLASS);

      if (MapPin._focusedMapPin === this) MapPin._focusedMapPin = null;

      this.updateMarkerZIndex();
    }

    private updateMarkerZIndex() {
      const shouldBeElevated = this._isHovered || this._isTouchActive || this._isContentFocused;

      if (shouldBeElevated) {
        if (this._defaultZIndex === null)
          this._defaultZIndex = this._leafletMarker.style.zIndex;

        this._leafletMarker.style.zIndex = "9999";

        return;
      }

      this._leafletMarker.style.zIndex = this._defaultZIndex ?? "";
      this._defaultZIndex = null;
    }

    private handleOnClick() {
      const returnPointId = this._returnPointId;

      if (!returnPointId) return;

      const event: MapPinClickEvent = new CustomEvent("map-pin-click", {
        detail: { returnPointId },
        bubbles: true,
        composed: true,
      });

      this.dispatchEvent(event);
    }

    public static async CreateMarker(data: MapPoint): Promise<LF.Marker> {
      await customElements.whenDefined("map-pin");

      const template = document.getElementById(TEMPLATE_ID) as HTMLTemplateElement;

      if (!template) throw new Error(`Template with id ${TEMPLATE_ID} not found`);

      const clone = template.content.cloneNode(true) as DocumentFragment;

      const mapPin = clone.querySelector("map-pin") as MapPin;

      if (!mapPin) throw new Error("MapPin element not found in the template");

      mapPin.setAttribute("data-return-point-id", data.returnPointId);
      mapPin.setAttribute("data-short-name", data.shortName);
      mapPin.setAttribute("data-bg-color", data.bgColor);
      mapPin.setAttribute("data-text-color", data.textColor);

      const icon = LF.divIcon({ html: mapPin });

      return LF.marker([data.latitude, data.longitude], { icon });
    }
}

// Register the custom element if it hasn't been registered yet
if (!customElements.get("map-pin")) customElements.define("map-pin", MapPin);
