import { TEMPLATE_ID } from "./map-pin.constants";

export class MapPin extends HTMLElement {
    connectedCallback() {
      console.log("MapPin connected");
    }

    disconnectedCallback() {}

    public static CloneMapPin(): MapPin {
      const template = document.getElementById(TEMPLATE_ID) as HTMLTemplateElement;

      if (!template) {
        throw new Error(`Template with id ${TEMPLATE_ID} not found`);
      }

      const clone = template.content.cloneNode(true) as DocumentFragment;
      const mapPin = clone.querySelector("map-pin") as MapPin;

      if (!mapPin) {
        throw new Error("MapPin element not found in the template");
      }

      return mapPin;
    }
}

customElements.define("map-pin", MapPin);
