import Styles from "./icon-button.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import { ButtonBase, type ButtonBaseProps } from "@components/base/button/button-base";

import type { ComponentChildren, JSX } from "preact";

export type IconButtonEdge = "none" | "start" | "end";

export type IconButtonProps = ButtonBaseProps & {
	icon?: string;
	edge?: IconButtonEdge;
	loading?: boolean;
	children?: ComponentChildren;
};

export default function IconButton(props: IconButtonProps): JSX.Element {
	const { icon, edge = "none", loading, children } = props;

	const isDisabled = !!props.disabled || !!loading;

	const rootClass = [
		Styles["icon-btn"],
		edge !== "none" ? Styles[`icon-btn--edge-${edge}`] : undefined,
		props.disabled ? Styles["icon-btn--disabled"] : undefined,
		loading ? Styles["icon-btn--loading"] : undefined,
	].filter(Boolean).join(" ");

	const content = icon ? <IconifyIcon icon={icon} /> : children;

	return (
		<ButtonBase {...props} rootClass={rootClass} isDisabled={isDisabled} isBusy={!!loading}>
			<span class={Styles["icon-btn__content"]}>
				{loading ? <span class={Styles["icon-btn__spinner"]} /> : content}
			</span>
		</ButtonBase>
	);
}