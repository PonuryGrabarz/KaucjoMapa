import Styles from "./form-field.module.sass";
import Input, { type InputProps } from "@components/input/input";

type FormFieldProps = Omit<InputProps, "error"> & {
  id: string;
  label: string;
  error?: string;
};

const getDefaultAdornment = (type: InputProps["type"]): string => {
  if (type === "email") return "ic:outline-email";
  if (type === "password") return "material-symbols:lock-outline";

  return "material-symbols:person-outline";
};

export default function FormField({
  id,
  label,
  error,
  fullWidth = true,
  startAdornment,
  "aria-describedby": ariaDescribedBy,
  ...inputProps
}: FormFieldProps) {
  const errorId = `${id}-error`;
  const describedBy =
    [
      typeof ariaDescribedBy === "string" ? ariaDescribedBy : undefined,
      error ? errorId : undefined,
    ]
      .filter(Boolean)
      .join(" ") || undefined;

  return (
    <div class={Styles["form-field"]} data-error={error ? "true" : undefined}>
      <label for={id}>{label}</label>

      <Input
        {...inputProps}
        id={id}
        class={Styles["form-field__input"]}
        fullWidth={fullWidth}
        startAdornment={startAdornment ?? getDefaultAdornment(inputProps.type)}
        error={Boolean(error)}
        aria-describedby={describedBy}
      />

      {error && (
        <span id={errorId} class={Styles["form-field__error"]} role="alert">
          {error}
        </span>
      )}
    </div>
  );
}
