import Styles from "./input.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import type { ComponentChildren, CSSProperties, InputHTMLAttributes, JSX, Ref } from "preact";

export type InputProps = Omit<InputHTMLAttributes<HTMLInputElement>, "class" | "style" | "ref"> & {
	class?: string;
	inputClass?: string;
	style?: string | CSSProperties;
	inputStyle?: string | CSSProperties;
	inputRef?: Ref<HTMLInputElement>;
	fullWidth?: boolean;
	error?: boolean;
	startAdornment?: string | ComponentChildren;
};

const renderStartAdornment = (adornment: string | ComponentChildren): ComponentChildren => {
	return typeof adornment === "string"
		? <IconifyIcon class={Styles["icon"]} icon={adornment} />
		: adornment;
};

export default function Input(props: InputProps): JSX.Element {
	const {
		fullWidth,
		error,
		startAdornment,
		class: className,
		inputClass,
		style,
		inputStyle,
		inputRef,
		...inputProps
	} = props;

	const hasStartAdornment = startAdornment !== undefined && startAdornment !== null && startAdornment !== false;

	const rootClasses = [
		Styles["input"],
		fullWidth ? Styles["input--full-width"] : undefined,
		hasStartAdornment ? Styles["input--adorned-start"] : undefined,
		props.disabled ? Styles["input--disabled"] : undefined,
		error ? Styles["input--error"] : undefined,
		className ?? "",
	].filter(Boolean).join(" ");

	const controlClasses = [Styles["input__control"], inputClass].filter(Boolean).join(" ");

	return (
		<span class={rootClasses} style={style}>
			{hasStartAdornment && (
				<span class={Styles["input__start-adornment"]} aria-hidden="true">
					{renderStartAdornment(startAdornment)}
				</span>
			)}

			<input
				{...inputProps}
				ref={inputRef}
				class={controlClasses}
				style={inputStyle}
				aria-invalid={error ? "true" : inputProps["aria-invalid"]}
			/>
		</span>
	);
}
