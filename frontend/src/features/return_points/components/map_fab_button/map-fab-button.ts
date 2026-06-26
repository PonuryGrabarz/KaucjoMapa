export const MAP_FAB_CLICK_EVENT = "map-fab-click";

export type MapFabClickEvent = CustomEvent<void>;

export class MapFab extends HTMLElement {
  private readonly _boundHandleClick = this.handleClick.bind(this);

  private readonly _buttonElement: HTMLButtonElement;

  private _isDisabled = false;

  constructor() {
    super();

    this._buttonElement = this.querySelector("button") as HTMLButtonElement;

    this._buttonElement.disabled = this._isDisabled;
  }

  connectedCallback() {
    this._buttonElement.addEventListener("click", this._boundHandleClick);
  }

  disconnectedCallback() {
    this._buttonElement.removeEventListener("click", this._boundHandleClick);
  }

  public get disabled(): boolean {
    return this._isDisabled;
  }

  public set disabled(value: boolean) {
    this._isDisabled = value;
    this._buttonElement.disabled = value;
  }

  private handleClick() {
    if (this._isDisabled) return;

    const event: MapFabClickEvent = new CustomEvent(MAP_FAB_CLICK_EVENT, {
      bubbles: true,
      composed: true,
    });

    this.dispatchEvent(event);
  }
}

// Register the custom element if it hasn't been registered yet
if (!customElements.get("map-fab")) customElements.define("map-fab", MapFab);
