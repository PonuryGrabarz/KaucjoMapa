import Styles from "./calculator.module.sass";
import CalculatorItem, {
  type CalculatorItemVariant,
} from "../calculator_item/calculator-item";
import CalculatorSummary from "../calculator_summary/calculator-summary";

import { useState } from "preact/hooks";

import type { JSX } from "preact";

type PackageType = "pet" | "can" | "glass";

type PackageDefinition = {
  id: PackageType;
  label: string;
  description: string;
  icon: string;
  variant: CalculatorItemVariant;
  depositGrosze: number;
};

type Quantities = Record<PackageType, number>;

const PACKAGE_DEFINITIONS: readonly PackageDefinition[] = [
  {
    id: "pet",
    label: "Butelki PET",
    description: "do 3 l · 0,50 zł",
    icon: "mdi:bottle-soda-classic-outline",
    variant: "pet",
    depositGrosze: 50,
  },
  {
    id: "can",
    label: "Puszki",
    description: "do 1 l · 0,50 zł",
    icon: "mdi:cup-outline",
    variant: "can",
    depositGrosze: 50,
  },
  {
    id: "glass",
    label: "Szkło wielorazowe",
    description: "do 1,5 l · 1,00 zł",
    icon: "mdi:bottle-wine-outline",
    variant: "glass",
    depositGrosze: 100,
  },
];

const INITIAL_QUANTITIES: Quantities = {
  pet: 0,
  can: 0,
  glass: 0,
};

export default function Calculator(): JSX.Element {
  const [quantities, setQuantities] = useState<Quantities>(INITIAL_QUANTITIES);

  const totalCount = PACKAGE_DEFINITIONS.reduce(
    (total, packageDefinition) => total + quantities[packageDefinition.id],
    0,
  );
  const totalDepositGrosze = PACKAGE_DEFINITIONS.reduce(
    (total, packageDefinition) => (
      total + quantities[packageDefinition.id] * packageDefinition.depositGrosze
    ),
    0,
  );

  const updateQuantity = (packageType: PackageType, quantity: number): void => {
    setQuantities((currentQuantities) => {
      if (currentQuantities[packageType] === quantity) return currentQuantities;

      return {
        ...currentQuantities,
        [packageType]: quantity,
      };
    });
  };

  return (
    <section
      class={Styles["calculator-page"]}
      aria-labelledby="calculator-page-title"
      data-calculator=""
    >
      <h1 id="calculator-page-title" class={Styles["calculator-page__title"]}>
        Kalkulator zwrotów opakowań kaucyjnych
      </h1>

      <div class={Styles["calculator-page__content"]}>
        <div class={Styles["calculator-page__planner"]}>
          {PACKAGE_DEFINITIONS.map((packageDefinition) => (
            <CalculatorItem
              key={packageDefinition.id}
              id={`calculator-${packageDefinition.id}`}
              label={packageDefinition.label}
              description={packageDefinition.description}
              icon={packageDefinition.icon}
              variant={packageDefinition.variant}
              quantity={quantities[packageDefinition.id]}
              onQuantityChange={(quantity) => (
                updateQuantity(packageDefinition.id, quantity)
              )}
            />
          ))}
        </div>

        <CalculatorSummary
          packageCount={totalCount}
          depositGrosze={totalDepositGrosze}
        />
      </div>
    </section>
  );
}
