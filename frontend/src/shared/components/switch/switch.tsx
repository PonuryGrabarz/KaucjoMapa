import Styles from "./switch.module.sass";

import { useState } from "preact/hooks";

import type {
	ButtonHTMLAttributes,
	ComponentChildren,
	CSSProperties,
	JSX,
	Ref,
	TargetedMouseEvent,
} from "preact";

export const SWITCH_CHANGE_EVENT = "switch-change";

export type SwitchChangeEvent = CustomEvent<{ checked: boolean }>;

export type SwitchChangeHandler = (
	checked: boolean,
	event: TargetedMouseEvent<HTMLButtonElement>,
) => void;

export type SwitchProps = Omit<
	ButtonHTMLAttributes<HTMLButtonElement>,
	| "aria-checked"
	| "aria-pressed"
	| "checked"
	| "children"
	| "class"
	| "className"
	| "dangerouslySetInnerHTML"
	| "defaultChecked"
	| "disabled"
	| "onChange"
	| "ref"
	| "role"
	| "style"
> & {
	checked?: boolean;
	defaultChecked?: boolean;
	disabled?: boolean;
	onChange?: SwitchChangeHandler;
	children?: ComponentChildren;
	class?: string;
	className?: string;
	trackClass?: string;
	style?: string | CSSProperties;
	trackStyle?: string | CSSProperties;
	buttonRef?: Ref<HTMLButtonElement>;
};

export default function Switch(props: SwitchProps): JSX.Element {
	const {
		checked,
		defaultChecked = false,
		disabled = false,
		onChange,
		onClick,
		children,
		class: className,
		className: nativeClassName,
		trackClass,
		style,
		trackStyle,
		buttonRef,
		type = "button",
		...buttonProps
	} = props;

	const [internalChecked, setInternalChecked] = useState(defaultChecked);

	const isControlled = checked !== undefined;

	const isChecked = checked ?? internalChecked;

	const hasLabel = children !== undefined && children !== null && children !== false;

	const rootClasses = [
		Styles["switch"],
		className,
		nativeClassName,
	].filter(Boolean).join(" ");

	const trackClasses = [Styles["switch__track"], trackClass].filter(Boolean).join(" ");

	const handleClick = (event: TargetedMouseEvent<HTMLButtonElement>): void => {
		onClick?.(event);

		if (event.defaultPrevented || disabled) return;

		const nextChecked = !isChecked;

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
		<button
			{...buttonProps}
			ref={buttonRef}
			class={rootClasses}
			style={style}
			type={type}
			role="switch"
			aria-checked={isChecked}
			disabled={disabled}
			onClick={handleClick}
			data-switch=""
			data-state={isChecked ? "on" : "off"}
			data-on={isChecked ? "true" : undefined}
			data-off={isChecked ? undefined : "true"}
			data-disabled={disabled ? "true" : undefined}
			data-has-label={hasLabel ? "true" : undefined}
		>
			<span
				class={trackClasses}
				style={trackStyle}
				aria-hidden="true"
				data-switch-track=""
			>
				<span class={Styles["switch__thumb"]} data-switch-thumb="" />
			</span>

			{hasLabel && (
				<span class={Styles["switch__label"]} data-switch-label="">
					{children}
				</span>
			)}
		</button>
	);
}
