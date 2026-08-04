/** Skrót 3/4/6/8-cyfrowy. Wąska lista zamiast dowolnego koloru CSS, bo wartości
 *  trafiają wprost do atrybutu `style` i do canvasa z danych użytkowników. */
export const HEX_COLOR_PATTERN = /^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$/i;

const SHORT_HEX_PATTERN = /^#([0-9a-f]{3})$/i;
const FULL_HEX_PATTERN = /^#([0-9a-f]{6})$/i;

export const isHexColor = (value: string | null | undefined): boolean =>
	Boolean(value && HEX_COLOR_PATTERN.test(value.trim()));

/** Sprowadza `#abc` i `#AABBCC` do `#aabbcc`. Zwraca `undefined` dla wszystkiego,
 *  czego nie przyjmuje `<input type="color">` — czyli także dla zapisów z alfą. */
export const normalizeHexColor = (value: string | null | undefined): string | undefined => {
	const color = value?.trim().toLowerCase();

	if (!color) return undefined;

	const short = color.match(SHORT_HEX_PATTERN);

	if (short?.[1]) return `#${[...short[1]].map((digit) => digit + digit).join("")}`;

	return FULL_HEX_PATTERN.test(color) ? color : undefined;
}
