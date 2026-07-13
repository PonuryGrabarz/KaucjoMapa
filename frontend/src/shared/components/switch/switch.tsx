import Styles from './switch.module.sass';

import { useRef, useState } from 'preact/hooks';

import type { ButtonHTMLAttributes, CSSProperties, JSX, TargetedMouseEvent } from 'preact';

export const SWITCH_CHANGE_EVENT = 'switch-change';

export type SwitchChangeEvent = CustomEvent<{ checked: boolean }>;
export type SwitchChangeHandler = (
  checked: boolean,
  event: TargetedMouseEvent<HTMLButtonElement>,
) => void;

export type SwitchProps = Omit<
  ButtonHTMLAttributes<HTMLButtonElement>,
  | 'aria-checked'
  | 'checked'
  | 'children'
  | 'class'
  | 'defaultChecked'
  | 'disabled'
  | 'onChange'
  | 'role'
  | 'style'
> & {
  checked?: boolean;
  defaultChecked?: boolean;
  disabled?: boolean;
  onChange?: SwitchChangeHandler;
  class?: string;
  trackClass?: string;
  style?: string | CSSProperties;
  trackStyle?: string | CSSProperties;
};

type SwitchRipple = {
  id: number;
  turnsOn: boolean;
};

export default function Switch(props: SwitchProps): JSX.Element {
  const {
    checked,
    defaultChecked = false,
    disabled,
    onChange,
    onClick,
    class: className,
    trackClass,
    style,
    trackStyle,
    type = 'button',
    ...buttonProps
  } = props;

  const [internalChecked, setInternalChecked] = useState(defaultChecked);
  const [ripples, setRipples] = useState<SwitchRipple[]>([]);
  const rippleIdRef = useRef(0);

  const isControlled = checked !== undefined;
  const isChecked = checked ?? internalChecked;

  const rootClasses = [Styles['switch'], className].filter(Boolean).join(' ');
  const trackClasses = [Styles['switch__track'], trackClass].filter(Boolean).join(' ');

  const removeRipple = (id: number): void => {
    setRipples((current) => current.filter((ripple) => ripple.id !== id));
  };

  const handleClick = (event: TargetedMouseEvent<HTMLButtonElement>): void => {
    onClick?.(event);

    if (event.defaultPrevented || disabled) return;

    const nextChecked = !isChecked;

    setRipples((current) => [...current, { id: ++rippleIdRef.current, turnsOn: nextChecked }]);

    if (!isControlled) setInternalChecked(nextChecked);

    onChange?.(nextChecked, event);

    const switchChangeEvent: SwitchChangeEvent = new CustomEvent(SWITCH_CHANGE_EVENT, {
      detail: { checked: nextChecked },
      bubbles: true,
      composed: true,
    });

    event.currentTarget.dispatchEvent(switchChangeEvent);
  };

  return (
    <span class={rootClasses} style={style}>
      <button
        {...buttonProps}
        class={trackClasses}
        style={trackStyle}
        type={type}
        role="switch"
        aria-checked={isChecked}
        disabled={disabled}
        onClick={handleClick}
      >
        <span class={Styles['switch__thumb']}>
          {ripples.map((ripple) => (
            <span
              key={ripple.id}
              class={[
                Styles['switch__ripple'],
                ripple.turnsOn ? Styles['switch__ripple--on'] : undefined,
              ]
                .filter(Boolean)
                .join(' ')}
              onAnimationEnd={() => removeRipple(ripple.id)}
            />
          ))}
        </span>
      </button>
    </span>
  );
}
