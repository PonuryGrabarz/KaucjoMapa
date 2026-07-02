export const SWITCH_CHANGE_EVENT = "switch-change";

export type SwitchChangeEvent = CustomEvent<{ checked: boolean }>;

export class Switch extends HTMLElement {
  private static readonly RIPPLE_CLASS = "switch__ripple";
  private static readonly RIPPLE_ON_MODIFIER_CLASS = "switch__ripple--on";
  private static readonly RIPPLE_DURATION_MS = 320;

  private readonly _boundHandleClick = this.handleClick.bind(this);
  private readonly _boundHandleKeydown = this.handleKeydown.bind(this);

  private readonly _trackElement: HTMLButtonElement;
  private readonly _thumbElement: HTMLSpanElement;

  private _isChecked = false;
  private _isDisabled = false;

  constructor() {
    super();

    this._trackElement = this.querySelector("button") as HTMLButtonElement;
    this._thumbElement = this._trackElement.querySelector("span") as HTMLSpanElement;

    this._isChecked = this.hasAttribute("checked");
    this._isDisabled = this.hasAttribute("disabled");

    this._trackElement.disabled = this._isDisabled;
    this._trackElement.setAttribute("aria-checked", String(this._isChecked));
  }

  connectedCallback() {
    this._trackElement.addEventListener("click", this._boundHandleClick);
    this._trackElement.addEventListener("keydown", this._boundHandleKeydown);
  }

  disconnectedCallback() {
    this._trackElement.removeEventListener("click", this._boundHandleClick);
    this._trackElement.removeEventListener("keydown", this._boundHandleKeydown);
  }

  public get checked(): boolean {
    return this._isChecked;
  }

  public set checked(value: boolean) {
    this._isChecked = value;
    this._trackElement.setAttribute("aria-checked", String(value));

    if (value) {
      this.setAttribute("checked", "");
    } else {
      this.removeAttribute("checked");
    }
  }

  public get disabled(): boolean {
    return this._isDisabled;
  }

  public set disabled(value: boolean) {
    this._isDisabled = value;
    this._trackElement.disabled = value;

    if (value) {
      this.setAttribute("disabled", "");
    } else {
      this.removeAttribute("disabled");
    }
  }

  private handleClick() {
    if (this._isDisabled) return;

    this.spawnRipple();
    this.toggle();
  }

  private handleKeydown(event: KeyboardEvent) {
    if (this._isDisabled) return;

    if (event.key !== " " && event.key !== "Enter") return;

    event.preventDefault();

    this.spawnRipple();
    this.toggle();
  }

  private toggle() {
    this.checked = !this._isChecked;

    const event: SwitchChangeEvent = new CustomEvent(SWITCH_CHANGE_EVENT, {
      detail: { checked: this._isChecked },
      bubbles: true,
      composed: true,
    });

    this.dispatchEvent(event);
  }

  private spawnRipple() {
    const ripple = document.createElement("span");

    ripple.className = this._isChecked ? Switch.RIPPLE_CLASS : `${Switch.RIPPLE_CLASS} ${Switch.RIPPLE_ON_MODIFIER_CLASS}`;

    this._thumbElement.appendChild(ripple);

    setTimeout(() => {
      ripple.remove();
    }, Switch.RIPPLE_DURATION_MS);
  }
}

if (!customElements.get("km-switch")) customElements.define("km-switch", Switch);
