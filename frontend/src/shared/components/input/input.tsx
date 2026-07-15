import Styles from "./input.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import type { ComponentChildren, CSSProperties, InputHTMLAttributes, JSX, Ref } from "preact";

type NativeInputProps = Omit<
	InputHTMLAttributes<HTMLInputElement>,
	"children" | "class" | "className" | "dangerouslySetInnerHTML" | "ref" | "style"
>;

export type InputProps = NativeInputProps & {
	class?: string;
	className?: string;
	inputClass?: string;
	style?: string | CSSProperties;
	inputStyle?: string | CSSProperties;
	inputRef?: Ref<HTMLInputElement>;
	fullWidth?: boolean;
	error?: boolean;
	startAdornment?: string | ComponentChildren;
	endAdornment?: string | ComponentChildren;
};

const hasAdornment = (adornment: ComponentChildren): boolean => {
	return adornment !== undefined && adornment !== null && adornment !== false;
};

const renderAdornment = (adornment: string | ComponentChildren): ComponentChildren => {
	return typeof adornment === "string"
		? <IconifyIcon icon={adornment} aria-hidden="true" />
		: adornment;
};

const isAriaInvalid = (
	ariaInvalid: InputHTMLAttributes<HTMLInputElement>["aria-invalid"],
): boolean => {
	return ariaInvalid === true
		|| ariaInvalid === "true"
		|| ariaInvalid === "grammar"
		|| ariaInvalid === "spelling";
};

export default function Input({
	fullWidth = false,
	error = false,
	startAdornment,
	endAdornment,
	class: className,
	className: nativeClassName,
	inputClass,
	style,
	inputStyle,
	inputRef,
	...inputProps
}: InputProps): JSX.Element {
	const hasStartAdornment = hasAdornment(startAdornment);

	const hasEndAdornment = hasAdornment(endAdornment);

	const hasError = error || isAriaInvalid(inputProps["aria-invalid"]);

	const isDisabled = inputProps.disabled === true;

	const isReadOnly = inputProps.readOnly === true || inputProps.readonly === true;

	const isRequired = inputProps.required === true;

	const inputType = typeof inputProps.type === "string" && inputProps.type
		? inputProps.type
		: "text";

	const rootClasses = [
		Styles["input"],
		className,
		nativeClassName,
	].filter(Boolean).join(" ");

	const controlClasses = [Styles["input__control"], inputClass].filter(Boolean).join(" ");

	return (
		<span
			class={rootClasses}
			style={style}
			data-input=""
			data-type={inputType}
			data-error={hasError ? "true" : undefined}
			data-disabled={isDisabled ? "true" : undefined}
			data-readonly={isReadOnly ? "true" : undefined}
			data-required={isRequired ? "true" : undefined}
			data-full-width={fullWidth ? "true" : undefined}
			data-adorned-start={hasStartAdornment ? "true" : undefined}
			data-adorned-end={hasEndAdornment ? "true" : undefined}
		>
			{hasStartAdornment && (
				<span class={Styles["input__start-adornment"]}>
					{renderAdornment(startAdornment)}
				</span>
			)}

			<input
				{...inputProps}
				ref={inputRef}
				class={controlClasses}
				style={inputStyle}
				aria-invalid={error ? "true" : inputProps["aria-invalid"]}
			/>

			{hasEndAdornment && (
				<span class={Styles["input__end-adornment"]}>
					{renderAdornment(endAdornment)}
				</span>
			)}
		</span>
	);
}
