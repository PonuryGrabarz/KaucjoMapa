import Styles from "./icon-button.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import { ButtonBase, type ButtonBaseProps } from "@components/base/button/button-base";

import type { ComponentChildren, JSX } from "preact";

export type IconButtonEdge = "none" | "start" | "end";

type DistributiveOmit<Props, Keys extends PropertyKey> = Props extends unknown
	? Omit<Props, Keys>
	: never;

type IconButtonNativeProps = DistributiveOmit<
	ButtonBaseProps,
	| "aria-label"
	| "aria-labelledby"
	| "children"
	| "dangerouslySetInnerHTML"
	| "isBusy"
	| "rootClass"
>;

type IconButtonAccessibleName =
	| {
		"aria-label": string;
		"aria-labelledby"?: string;
	}
	| {
		"aria-label"?: never;
		"aria-labelledby": string;
	};

export type IconButtonProps = IconButtonNativeProps & IconButtonAccessibleName & {
	icon?: string;
	edge?: IconButtonEdge;
	loading?: boolean;
	children?: ComponentChildren;
};

export default function IconButton({
	icon,
	edge = "none",
	loading = false,
	children,
	...buttonBaseProps
}: IconButtonProps): JSX.Element {
	const isDisabled = Boolean(buttonBaseProps.disabled || loading);

	const content = icon ? <IconifyIcon icon={icon} aria-hidden="true" /> : children;

	return (
		<ButtonBase
			{...buttonBaseProps}
			rootClass={Styles["icon-btn"]}
			disabled={isDisabled}
			isBusy={loading}
			data-icon-button=""
			data-edge={edge}
			data-loading={loading ? "true" : undefined}
			data-disabled={isDisabled ? "true" : undefined}
		>
			<span class={Styles["icon-btn__content"]}>
				{loading
					? <span class={Styles["icon-btn__spinner"]} aria-hidden="true" />
					: content}
			</span>
		</ButtonBase>
	);
}