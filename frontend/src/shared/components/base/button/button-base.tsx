import Ripple from "@components/ripple/ripple";

import Styles from "./button-base.module.sass";

import { createElement } from "preact";

import type {
	AnchorHTMLAttributes,
	ButtonHTMLAttributes,
	ComponentChildren,
} from "preact";

export type ButtonType = "button" | "submit" | "reset";

type ButtonBaseCommonProps = {
	rootClass: string;
	disableRipple?: boolean;
	isBusy?: boolean;
	class?: string;
	className?: string;
	children?: ComponentChildren;
};

type LinkOnlyAttribute =
	| "download"
	| "hrefLang"
	| "hreflang"
	| "media"
	| "ping"
	| "referrerPolicy"
	| "referrerpolicy"
	| "rel"
	| "target";

type ButtonOnlyAttribute =
	| "command"
	| "commandFor"
	| "commandfor"
	| "form"
	| "formAction"
	| "formaction"
	| "formEncType"
	| "formenctype"
	| "formMethod"
	| "formmethod"
	| "formNoValidate"
	| "formnovalidate"
	| "formTarget"
	| "formtarget"
	| "name"
	| "popoverTarget"
	| "popovertarget"
	| "popoverTargetAction"
	| "popovertargetaction"
	| "value";

type NeverAttributes<Attribute extends PropertyKey> = {
	[Key in Attribute]?: never;
};

export type ButtonBaseLinkProps = ButtonBaseCommonProps
	& Omit<
		AnchorHTMLAttributes<HTMLAnchorElement>,
		keyof ButtonBaseCommonProps | ButtonOnlyAttribute | "href" | "ref"
	>
	& NeverAttributes<ButtonOnlyAttribute>
	& {
		href: string;
		disabled?: boolean;
	};

export type ButtonBaseButtonProps = ButtonBaseCommonProps
	& Omit<
		ButtonHTMLAttributes<HTMLButtonElement>,
		keyof ButtonBaseCommonProps | LinkOnlyAttribute | "disabled" | "href" | "ref" | "type"
	>
	& NeverAttributes<LinkOnlyAttribute>
	& {
		href?: never;
		disabled?: boolean;
		type?: ButtonType;
	};

export type ButtonBaseProps = ButtonBaseLinkProps | ButtonBaseButtonProps;

function getClasses(rootClass: string, className?: string, nativeClassName?: string): string {
	return [Styles["button-base"], rootClass, className, nativeClassName].filter(Boolean).join(" ");
}

function getAriaBusy(
	isBusy: boolean,
	ariaBusy: AnchorHTMLAttributes<HTMLAnchorElement>["aria-busy"],
): AnchorHTMLAttributes<HTMLAnchorElement>["aria-busy"] {
	return isBusy ? "true" : ariaBusy;
}

function getLinkClickHandler(
	disabled: boolean,
	onClick: AnchorHTMLAttributes<HTMLAnchorElement>["onClick"],
): AnchorHTMLAttributes<HTMLAnchorElement>["onClick"] {
	if (!disabled && !onClick) return undefined;

	return (e): void => {
		if (disabled) {
			e.preventDefault();
			return;
		}

		onClick?.(e);
	};
}

function getRipple(disableRipple: boolean, disabled: boolean) {
	return disableRipple ? null : createElement(Ripple, { disabled });
}

export function ButtonBase({
	rootClass,
	disableRipple = false,
	isBusy = false,
	class: className,
	className: nativeClassName,
	children,
	...elementProps
}: ButtonBaseProps) {
	const classes = getClasses(rootClass, className, nativeClassName);

	if (elementProps.href !== undefined) {
		const {
			href,
			disabled = false,
			onClick,
			"aria-busy": ariaBusy,
			"aria-disabled": ariaDisabled,
			...anchorProps
		} = elementProps;

		return (
			<a
				{...anchorProps}
				class={classes}
				href={href}
				aria-busy={getAriaBusy(isBusy, ariaBusy)}
				aria-disabled={disabled ? "true" : ariaDisabled}
				onClick={getLinkClickHandler(disabled, onClick)}
			>
				{children}{getRipple(disableRipple, disabled)}
			</a>
		);
	}

	const {
		disabled = false,
		type = "button",
		"aria-busy": ariaBusy,
		...buttonProps
	} = elementProps;

	return (
		<button
			{...buttonProps}
			class={classes}
			type={type}
			disabled={disabled}
			aria-busy={getAriaBusy(isBusy, ariaBusy)}
		>
			{children}{getRipple(disableRipple, disabled)}
		</button>
	);
}
