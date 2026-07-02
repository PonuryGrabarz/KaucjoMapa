export const INPUT_CHANGE_EVENT = 'input-change';

export type InputChangeEvent = CustomEvent<{ value: string }>;

export class Input extends HTMLElement {
  private readonly _boundHandleInput = this.handleInput.bind(this);

  private readonly _inputElement: HTMLInputElement;

  private _isDisabled = false;

  constructor() {
    super();

    const inputElement = this.querySelector('input');

    if (!(inputElement instanceof HTMLInputElement)) {
      throw new Error('km-input requires an inner input element.');
    }

    this._inputElement = inputElement;
    this._isDisabled = this.hasAttribute('disabled') || this._inputElement.disabled;

    this._inputElement.disabled = this._isDisabled;

    if (this._isDisabled) {
      this.setAttribute('disabled', '');
    } else {
      this.removeAttribute('disabled');
    }
  }

  connectedCallback() {
    this._inputElement.addEventListener('input', this._boundHandleInput);
  }

  disconnectedCallback() {
    this._inputElement.removeEventListener('input', this._boundHandleInput);
  }

  public get value(): string {
    return this._inputElement.value;
  }

  public set value(value: string) {
    this._inputElement.value = value;
  }

  public get disabled(): boolean {
    return this._isDisabled;
  }

  public set disabled(value: boolean) {
    this._isDisabled = value;
    this._inputElement.disabled = value;

    if (value) {
      this.setAttribute('disabled', '');
    } else {
      this.removeAttribute('disabled');
    }
  }

  private handleInput() {
    const event: InputChangeEvent = new CustomEvent(INPUT_CHANGE_EVENT, {
      detail: { value: this._inputElement.value },
      bubbles: true,
      composed: true,
    });

    this.dispatchEvent(event);
  }
}

if (!customElements.get('km-input')) customElements.define('km-input', Input);
