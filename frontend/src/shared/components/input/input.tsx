import Styles from "./input.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import { useState } from "preact/hooks";

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
	showPasswordToggle?: boolean;
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
	showPasswordToggle = false,
	class: className,
	className: nativeClassName,
	inputClass,
	style,
	inputStyle,
	inputRef,
	...inputProps
}: InputProps): JSX.Element {
	const [isPasswordVisible, setIsPasswordVisible] = useState(false);

	const hasStartAdornment = hasAdornment(startAdornment);

	const hasError = error || isAriaInvalid(inputProps["aria-invalid"]);

	const isDisabled = inputProps.disabled === true;

	const isReadOnly = inputProps.readOnly === true || inputProps.readonly === true;

	const isRequired = inputProps.required === true;

	const inputType = typeof inputProps.type === "string" && inputProps.type
		? inputProps.type
		: "text";

	const hasPasswordToggle = showPasswordToggle && inputType === "password";

	const hasEndAdornment = hasAdornment(endAdornment) || hasPasswordToggle;

	const renderedInputType = hasPasswordToggle && isPasswordVisible
		? "text"
		: inputProps.type;

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
			data-password-visible={hasPasswordToggle && isPasswordVisible ? "true" : undefined}
		>
			{hasStartAdornment && (
				<span class={Styles["input__start-adornment"]}>
					{renderAdornment(startAdornment)}
				</span>
			)}

			<input
				{...inputProps}
				type={renderedInputType}
				ref={inputRef}
				class={controlClasses}
				style={inputStyle}
				aria-invalid={error ? "true" : inputProps["aria-invalid"]}
			/>

			{hasEndAdornment && (
				<span class={Styles["input__end-adornment"]}>
					{hasAdornment(endAdornment) && renderAdornment(endAdornment)}

					{hasPasswordToggle && (
						<button
							class={Styles["input__password-toggle"]}
							type="button"
							disabled={isDisabled}
							aria-label={isPasswordVisible ? "Ukryj hasło" : "Pokaż hasło"}
							aria-pressed={isPasswordVisible}
							aria-controls={typeof inputProps.id === "string" ? inputProps.id : undefined}
							onClick={() => setIsPasswordVisible((isVisible) => !isVisible)}
						>
							<IconifyIcon
								icon={isPasswordVisible ? "mdi:eye-off-outline" : "mdi:eye-outline"}
								aria-hidden="true"
							/>
						</button>
					)}
				</span>
			)}
		</span>
	);
}
