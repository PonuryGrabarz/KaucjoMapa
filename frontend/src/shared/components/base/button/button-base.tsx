import Ripple from "@components/ripple/ripple";

import type { AriaAttributes, ComponentChildren, CSSProperties, JSX } from "preact";

export type ButtonType = "button" | "submit" | "reset";

export type ButtonBaseProps = Pick<AriaAttributes, "aria-current" | "aria-label"> & {
	disabled?: boolean;
	disableRipple?: boolean;
	href?: string;
	type?: ButtonType;
	class?: string;
	style?: string | CSSProperties;
	onClick?: (e: MouseEvent) => void;
};

type ButtonBaseComponentProps = ButtonBaseProps & {
	rootClass: string;
	isDisabled?: boolean;
	isBusy?: boolean;
	children?: ComponentChildren;
};

export function ButtonBase({
	rootClass,
	isDisabled = false,
	isBusy = false,
	href,
	type = "button",
	disableRipple,
	class: className,
	style,
	onClick,
	"aria-current": ariaCurrent,
	"aria-label": ariaLabel,
	children,
}: ButtonBaseComponentProps): JSX.Element {
	const isLink = (href?.length ?? 0) > 0;
  
	const classes = [rootClass, className ?? ""].filter(Boolean).join(" ");

	const handleClick = (e: MouseEvent): void => {
		if (isDisabled) {
			e.preventDefault();
			e.stopImmediatePropagation();
			return;
		}

		onClick?.(e);
	};

	const ripple = <Ripple disabled={isDisabled || disableRipple} />;

	return isLink ? (
		<a
			class={classes}
			style={style}
			href={isDisabled ? undefined : href}
			role="button"
			aria-disabled={isDisabled ? "true" : "false"}
			aria-busy={isBusy ? "true" : undefined}
			aria-current={ariaCurrent}
			aria-label={ariaLabel}
			tabindex={isDisabled ? -1 : 0}
			onClick={handleClick}
		>
			{children}{ripple}
		</a>
	) : (
		<button
			class={classes}
			style={style}
			type={type}
			disabled={isDisabled}
			aria-busy={isBusy ? "true" : undefined}
			aria-current={ariaCurrent}
			aria-label={ariaLabel}
			onClick={handleClick}
		>
			{children}{ripple}
		</button>
	);
}
