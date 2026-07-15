import Styles from "./button.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import { ButtonBase, type ButtonBaseProps } from "@components/base/button/button-base";

import type { ComponentChildren } from "preact";

export type ButtonVariant = "text" | "contained" | "outlined";
export type ButtonLoadingPosition = "start" | "center" | "end";

type DistributiveOmit<Props, Keys extends PropertyKey> = Props extends unknown
	? Omit<Props, Keys>
	: never;

type ButtonNativeProps = DistributiveOmit<
	ButtonBaseProps,
	"children" | "dangerouslySetInnerHTML" | "isBusy" | "rootClass"
>;

export type ButtonProps = ButtonNativeProps & {
	variant?: ButtonVariant;
	loading?: boolean;
	loadingPosition?: ButtonLoadingPosition;
	fullWidth?: boolean;
	startIcon?: string | ComponentChildren;
	endIcon?: string | ComponentChildren;
	children?: ComponentChildren;
};

const renderIcon = (icon: string | ComponentChildren): ComponentChildren => {
	return typeof icon === "string"
		? <IconifyIcon class={Styles["btn__icon"]} icon={icon} aria-hidden="true" />
		: icon;
};

export default function Button({
	variant = "text",
	loading = false,
	loadingPosition = "center",
	fullWidth = false,
	startIcon,
	endIcon,
	children,
	...buttonBaseProps
}: ButtonProps) {
	const spinner = <span class={Styles["btn__spinner"]} aria-hidden="true" />;
	const isDisabled = Boolean(buttonBaseProps.disabled || loading);

	return (
		<ButtonBase
			{...buttonBaseProps}
			rootClass={Styles["btn"]}
			disabled={isDisabled}
			isBusy={loading}
			data-button=""
			data-variant={variant}
			data-loading={loading ? "true" : undefined}
			data-loading-position={loading ? loadingPosition : undefined}
			data-full-width={fullWidth ? "true" : undefined}
			data-disabled={isDisabled ? "true" : undefined}
		>
			{loading && loadingPosition === "start" ? spinner : renderIcon(startIcon)}

			<span class={Styles["btn__label"]}>{children}</span>

			{loading && loadingPosition === "center" && (
				<span class={Styles["btn__loading-center"]}>{spinner}</span>
			)}

			{loading && loadingPosition === "end" ? spinner : renderIcon(endIcon)}
		</ButtonBase>
	);
}
