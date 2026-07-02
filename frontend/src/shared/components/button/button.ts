export const BUTTON_CLICK_EVENT = 'button-click';

export type ButtonClickEvent = CustomEvent<void>;

export class Button extends HTMLElement {
  private static readonly RIPPLE_CLASS = 'button__ripple';
  private static readonly RIPPLE_DURATION_MS = 340;

  private readonly _boundHandleClick = this.handleClick.bind(this);
  private readonly _boundHandlePointerDown = this.handlePointerDown.bind(this);
  private readonly _boundHandleKeydown = this.handleKeydown.bind(this);

  private readonly _buttonElement: HTMLButtonElement;

  private _isDisabled = false;

  constructor() {
    super();

    const buttonElement = this.querySelector('button');

    if (!(buttonElement instanceof HTMLButtonElement)) {
      throw new Error('km-button requires an inner button element.');
    }

    this._buttonElement = buttonElement;
    this._isDisabled = this.hasAttribute('disabled') || this._buttonElement.disabled;

    this.syncDisabledState();
  }

  connectedCallback() {
    this._buttonElement.addEventListener('click', this._boundHandleClick);
    this._buttonElement.addEventListener('pointerdown', this._boundHandlePointerDown);
    this._buttonElement.addEventListener('keydown', this._boundHandleKeydown);
  }

  disconnectedCallback() {
    this._buttonElement.removeEventListener('click', this._boundHandleClick);
    this._buttonElement.removeEventListener('pointerdown', this._boundHandlePointerDown);
    this._buttonElement.removeEventListener('keydown', this._boundHandleKeydown);
  }

  public get disabled(): boolean {
    return this._isDisabled;
  }

  public set disabled(value: boolean) {
    this._isDisabled = value;
    this.syncDisabledState();
  }

  private handleClick(event: MouseEvent) {
    if (this._isDisabled) {
      event.preventDefault();
      event.stopPropagation();
      return;
    }

    const buttonClickEvent: ButtonClickEvent = new CustomEvent(BUTTON_CLICK_EVENT, {
      bubbles: true,
      composed: true,
    });

    this.dispatchEvent(buttonClickEvent);
  }

  private handlePointerDown(event: PointerEvent) {
    if (this._isDisabled || event.button !== 0) return;

    const { x, y } = this.getRippleOrigin(event);
    this.spawnRipple(x, y);
  }

  private handleKeydown(event: KeyboardEvent) {
    if (this._isDisabled) return;

    if (event.key !== ' ' && event.key !== 'Enter') return;

    const { width, height } = this._buttonElement.getBoundingClientRect();
    this.spawnRipple(width / 2, height / 2);
  }

  private getRippleOrigin(event: PointerEvent): { x: number; y: number } {
    if (event.target === this._buttonElement) {
      return { x: event.offsetX, y: event.offsetY };
    }

    const rect = this._buttonElement.getBoundingClientRect();

    return {
      x: event.clientX - rect.left,
      y: event.clientY - rect.top,
    };
  }

  private spawnRipple(x: number, y: number) {
    const { width, height } = this._buttonElement.getBoundingClientRect();
    const size = Math.max(width, height) * 2;
    const ripple = document.createElement('span');

    ripple.className = Button.RIPPLE_CLASS;
    ripple.style.width = `${size}px`;
    ripple.style.height = `${size}px`;
    ripple.style.left = `${x}px`;
    ripple.style.top = `${y}px`;

    this._buttonElement.appendChild(ripple);

    setTimeout(() => {
      ripple.remove();
    }, Button.RIPPLE_DURATION_MS);
  }

  private syncDisabledState() {
    this._buttonElement.disabled = this._isDisabled;
    this._buttonElement.setAttribute('aria-disabled', String(this._isDisabled));

    if (this._isDisabled) {
      this.setAttribute('disabled', '');
    } else {
      this.removeAttribute('disabled');
    }
  }
}

if (!customElements.get('km-button')) customElements.define('km-button', Button);
