export const MIN_QUANTITY = 0;
export const MAX_QUANTITY = 9999;

export function clampQuantity(value: number): number {
  if (!Number.isFinite(value)) return value > 0 ? MAX_QUANTITY : MIN_QUANTITY;

  return Math.min(MAX_QUANTITY, Math.max(MIN_QUANTITY, Math.trunc(value)));
}

export function normalizeQuantity(rawValue: string): number {
  if (!/^\d+$/.test(rawValue)) return MIN_QUANTITY;

  const digits = rawValue.replace(/^0+(?=\d)/, "");
  if (digits.length > String(MAX_QUANTITY).length) return MAX_QUANTITY;

  return clampQuantity(Number(digits));
}

export function formatDeposit(depositGrosze: number): string {
  const normalizedDeposit = Math.max(0, Math.trunc(depositGrosze));
  const zloty = Math.floor(normalizedDeposit / 100)
    .toString()
    .replace(/\B(?=(\d{3})+(?!\d))/g, " ");
  const grosze = String(normalizedDeposit % 100).padStart(2, "0");

  return `${zloty},${grosze} zł`;
}

export function getPackageCountLabel(count: number): string {
  if (count === 1) return "opakowanie";

  const lastTwoDigits = count % 100;
  const lastDigit = count % 10;

  return lastDigit >= 2 && lastDigit <= 4 && (lastTwoDigits < 12 || lastTwoDigits > 14)
    ? "opakowania"
    : "opakowań";
}
