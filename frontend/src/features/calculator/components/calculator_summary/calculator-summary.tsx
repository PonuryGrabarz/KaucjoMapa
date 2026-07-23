import Styles from "./calculator-summary.module.sass";

import { formatDeposit, getPackageCountLabel } from "../../lib/calculator";

import type { JSX } from "preact";

type CalculatorSummaryProps = {
  packageCount: number;
  depositGrosze: number;
};

export default function CalculatorSummary({
  packageCount,
  depositGrosze,
}: CalculatorSummaryProps): JSX.Element {
  return (
    <section
      class={Styles["calculator-summary"]}
      aria-live="polite"
      aria-atomic="true"
      data-calculator-summary=""
    >
      <h2 class={Styles["calculator-summary__title"]}>
        Szacowany zwrot kaucji
      </h2>

      <output class={Styles["calculator-summary__value"]}>
        {formatDeposit(depositGrosze)}
      </output>

      <p class={Styles["calculator-summary__count"]}>
        {packageCount} {getPackageCountLabel(packageCount)}
      </p>
    </section>
  );
}
