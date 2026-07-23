import Styles from "./calculator-item.module.sass";
import IconButton from "@components/icon_button/icon-button";
import IconifyIcon from "@components/iconify_icon/iconify-icon";
import Input from "@components/input/input";

import {
  MAX_QUANTITY,
  MIN_QUANTITY,
  clampQuantity,
  normalizeQuantity,
} from "../../lib/calculator";

import type { JSX } from "preact";

export type CalculatorItemVariant = "pet" | "can" | "glass";

type CalculatorItemProps = {
  id: string;
  label: string;
  description: string;
  icon: string;
  variant: CalculatorItemVariant;
  quantity: number;
  onQuantityChange: (quantity: number) => void;
};

export default function CalculatorItem({
  id,
  label,
  description,
  icon,
  variant,
  quantity,
  onQuantityChange,
}: CalculatorItemProps): JSX.Element {
  const descriptionId = `${id}-description`;

  return (
    <article class={Styles["calculator-item"]} data-calculator-item={variant}>
      <div class={Styles["calculator-item__info"]}>
        <span class={Styles["calculator-item__icon"]} data-variant={variant}>
          <IconifyIcon icon={icon} />
        </span>

        <div class={Styles["calculator-item__text"]}>
          <label class={Styles["calculator-item__label"]} for={id}>
            {label}
          </label>
          <span id={descriptionId} class={Styles["calculator-item__description"]}>
            {description}
          </span>
        </div>
      </div>

      <div class={Styles["calculator-item__stepper"]}>
        <IconButton
          icon="ic:baseline-minus"
          type="button"
          aria-label={`Zmniejsz liczbę: ${label}`}
          disabled={quantity <= MIN_QUANTITY}
          onClick={() => onQuantityChange(clampQuantity(quantity - 1))}
        />

        <Input
          class={Styles["calculator-item__input"]}
          id={id}
          name={variant}
          type="number"
          inputMode="numeric"
          min={MIN_QUANTITY}
          max={MAX_QUANTITY}
          step={1}
          value={quantity}
          aria-describedby={descriptionId}
          onInput={(event) => {
            const nextQuantity = normalizeQuantity(event.currentTarget.value);
            event.currentTarget.value = String(nextQuantity);
            onQuantityChange(nextQuantity);
          }}
        />

        <IconButton
          icon="ic:baseline-plus"
          type="button"
          aria-label={`Zwiększ liczbę: ${label}`}
          disabled={quantity >= MAX_QUANTITY}
          onClick={() => onQuantityChange(clampQuantity(quantity + 1))}
        />
      </div>
    </article>
  );
}
