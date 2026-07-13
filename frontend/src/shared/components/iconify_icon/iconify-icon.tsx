import Styles from "./iconify-icon.module.sass";

import type { JSX } from "preact";

type IconifyIconProps = {
	icon?: string;
	width?: string;
	height?: string;
	label?: string;
	colored?: boolean;
	class?: string;
};

export default function IconifyIcon(props: IconifyIconProps): JSX.Element {
	const [prefix, name] = (props.icon ?? "").split(":");
	
	const iconUrl = prefix && name ? `url("https://api.iconify.design/${prefix}/${name}.svg")` : "";

	const iconifyIconStyle = [
		iconUrl ? `--iconify-icon-url: ${iconUrl}` : undefined,
		props.width ? `--iconify-icon-width: ${props.width}` : undefined,
		props.height ? `--iconify-icon-height: ${props.height}` : undefined,
	].filter(Boolean).join("; ");

	const iconifyIconClasses = [
		Styles["iconify-icon"],
		props.colored ? Styles["colored"] : undefined,
		props.class,
	].filter(Boolean).join(" ");

	return (
		<span
			class={iconifyIconClasses}
			style={iconifyIconStyle}
			role={props.label ? "img" : undefined}
			aria-label={props.label}
			aria-hidden={props.label ? undefined : "true"}
		/>
	);
}
