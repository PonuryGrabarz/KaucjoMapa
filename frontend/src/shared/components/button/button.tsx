import Styles from "./button.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import { ButtonBase, type ButtonBaseProps } from "@components/base/button/button-base";

import type { ComponentChildren, JSX } from "preact";

export type ButtonVariant = "text" | "contained" | "outlined";
export type ButtonLoadingPosition = "start" | "center" | "end";

export type ButtonProps = ButtonBaseProps & {
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
		? <IconifyIcon class={Styles["btn__icon"]} icon={icon} />
		: icon;
};

export default function Button(props: ButtonProps): JSX.Element {
	const {
		variant = "text",
		loading,
		loadingPosition = "center",
		fullWidth,
		startIcon,
		endIcon,
		children,
	} = props;

	const spinner = <span class={Styles["btn__spinner"]} />;

	const isDisabled = !!props.disabled || !!loading;

	const rootClass = [
		Styles["btn"],
		Styles[`btn--${variant}`],
		fullWidth ? Styles["btn--full-width"] : undefined,
		props.disabled ? Styles["btn--disabled"] : undefined,
		loading ? Styles["btn--loading"] : undefined,
	].filter(Boolean).join(" ");

	const labelClass = [
		Styles["btn__label"],
		loading && loadingPosition === "center" ? Styles["btn--loading-hidden-label"] : undefined,
	].filter(Boolean).join(" ");

	return (
		<ButtonBase {...props} rootClass={rootClass} isDisabled={isDisabled} isBusy={!!loading}>
			{loading && loadingPosition === "start" ? spinner : renderIcon(startIcon)}

			<span class={labelClass}>{children}</span>

			{loading && loadingPosition === "center" && (
				<span class={Styles["btn__loading-center"]}>{spinner}</span>
			)}

			{loading && loadingPosition === "end" ? spinner : renderIcon(endIcon)}
		</ButtonBase>
	);
}
