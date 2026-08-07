const fullDateTimeFormatter = new Intl.DateTimeFormat("pl-PL", {
	day: "numeric",
	month: "long",
	year: "numeric",
	hour: "2-digit",
	minute: "2-digit",
});

const shortDateFormatter = new Intl.DateTimeFormat("pl-PL", { day: "numeric", month: "short" });

const parsedDate = (value: string): Date | null => {
	const date = new Date(value);

	return Number.isNaN(date.getTime()) ? null : date;
}

const pad = (value: number): string => String(value).padStart(2, "0");

const dateKey = (date: Date): string => 
  `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;

const preciseDistanceFormatter = new Intl.NumberFormat("pl-PL", { maximumFractionDigits: 1 });
const roundDistanceFormatter = new Intl.NumberFormat("pl-PL", { maximumFractionDigits: 0 });

export const distanceLabel = (km: number): string => {
	if (!Number.isFinite(km) || km < 0) return "";

	if (km < 1) return `${Math.round(km * 1000)} m`;

	return km < 10
		? `${preciseDistanceFormatter.format(km)} km`
		: `${roundDistanceFormatter.format(km)} km`;
}

export const truncateText = (value: string, maxLength: number): string =>
	value.length <= maxLength ? value : `${value.slice(0, maxLength - 1)}…`;

export const pluralPl = (count: number, one: string, few: string, many: string): string => {
	if (count === 1) return one;

	const lastTwoDigits = Math.abs(count) % 100;

	const lastDigit = Math.abs(count) % 10;

	if (lastDigit >= 2 && lastDigit <= 4 && (lastTwoDigits < 12 || lastTwoDigits > 14)) return few;

	return many;
}

export const unreadCountText = (count: number): string => {
	if (count === 0) return "Brak nowych powiadomień";

	const noun = pluralPl(
		count,
		"nieprzeczytane powiadomienie",
		"nieprzeczytane powiadomienia",
		"nieprzeczytanych powiadomień",
	);

	return `${count} ${noun}`;
}

export const relativeTimeLabel = (value: string, now = new Date()): string => {
	const date = parsedDate(value);

	if (!date) return "Nieznany czas";

	const elapsedMinutes = Math.floor((now.getTime() - date.getTime()) / 60_000);

	if (elapsedMinutes >= 0 && elapsedMinutes < 1) return "przed chwilą";

	if (elapsedMinutes >= 1 && elapsedMinutes < 60) return `${elapsedMinutes} min temu`;

	const elapsedHours = Math.floor(elapsedMinutes / 60);

	if (elapsedHours >= 1 && elapsedHours < 24) return `${elapsedHours} godz. temu`;

	const yesterday = new Date(now);

	yesterday.setDate(yesterday.getDate() - 1);

	if (dateKey(date) === dateKey(yesterday)) return "wczoraj";

	return shortDateFormatter.format(date);
}

export const dateTimeLabel = (value: string): string => {
	const date = parsedDate(value);

	return date ? fullDateTimeFormatter.format(date) : "Nieznana data";
}