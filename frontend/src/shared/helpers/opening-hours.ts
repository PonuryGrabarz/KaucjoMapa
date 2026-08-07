import type { OpeningHours } from "../api/responses";

const WEEKDAYS = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"] as const;

const TIME_RANGE_PATTERN =
  /^(\d{1,2}):([0-5]\d)(?::([0-5]\d))?\s*-\s*(\d{1,2}):([0-5]\d)(?::([0-5]\d))?$/;

const SECONDS_PER_DAY = 24 * 60 * 60;

type Weekday = (typeof WEEKDAYS)[number];

type TimeRange = { start: number; end: number };

const OPENING_HOURS_RULES: readonly {
  readonly key: keyof OpeningHours;
  readonly days: readonly Weekday[];
}[] = [
  { key: "sat", days: ["sat"] },
  { key: "sun", days: ["sun"] },
  { key: "mon-fri", days: ["mon", "tue", "wed", "thu", "fri"] },
  { key: "mon-sat", days: ["mon", "tue", "wed", "thu", "fri", "sat"] },
  { key: "mon-sun", days: WEEKDAYS },
];

const WARSAW_MOMENT_FORMATTER = new Intl.DateTimeFormat("en-GB", {
  timeZone: "Europe/Warsaw",
  weekday: "short",
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
  hourCycle: "h23",
});

const isWeekday = (value: string): value is Weekday => (WEEKDAYS as readonly string[]).includes(value);

const previousDay = (day: Weekday): Weekday => {
  const index = WEEKDAYS.indexOf(day);

  return WEEKDAYS[(index + WEEKDAYS.length - 1) % WEEKDAYS.length]!;
};

const readHourValue = (hours: OpeningHours | undefined, key: keyof OpeningHours): string | undefined => {
  const value: unknown = hours?.[key];

  if (typeof value !== "string") return undefined;

  const trimmed = value.trim();

  return trimmed.length > 0 ? trimmed : undefined;
};

const parseTimeRange = (value: unknown): TimeRange | undefined => {
  const match = typeof value === "string" ? value.trim().match(TIME_RANGE_PATTERN) : null;

  if (!match) return undefined;

  const toSeconds = (
    hours: string | undefined,
    minutes: string | undefined,
    seconds: string | undefined = "0",
  ) => Number(hours ?? "0") * 3600 + Number(minutes ?? "0") * 60 + Number(seconds ?? "0");

  const start = toSeconds(match[1], match[2], match[3]);
  const end = toSeconds(match[4], match[5], match[6]);

  return start > SECONDS_PER_DAY || end > SECONDS_PER_DAY ? undefined : { start, end };
};

const rangeForDay = (openingHours: OpeningHours, day: Weekday): TimeRange | undefined => {
  for (const rule of OPENING_HOURS_RULES) {
    if (!rule.days.includes(day)) continue;

    const range = parseTimeRange(readHourValue(openingHours, rule.key));

    if (range) return range;
  }

  return undefined;
};

const warsawMoment = (now: Date): { day: Weekday; secondsOfDay: number } | undefined => {
  const parts = new Map(
    WARSAW_MOMENT_FORMATTER.formatToParts(now).map((part) => [part.type, part.value]),
  );

  const weekday = parts.get("weekday")?.toLowerCase();
  const hour = Number(parts.get("hour"));
  const minute = Number(parts.get("minute"));
  const second = Number(parts.get("second"));

  if (!weekday || !isWeekday(weekday)) return undefined;

  if (!Number.isFinite(hour) || !Number.isFinite(minute) || !Number.isFinite(second))
    return undefined;

  return { day: weekday, secondsOfDay: hour * 3600 + minute * 60 + second };
};

const DISPLAY_ORDER: readonly (keyof OpeningHours)[] = [
  "mon-fri",
  "mon-sat",
  "mon-sun",
  "sat",
  "sun",
];

export const openingHoursEntries = (
  hours: OpeningHours | undefined,
): [keyof OpeningHours, string][] =>
  DISPLAY_ORDER.flatMap((key) => {
    const value = readHourValue(hours, key);

    return value ? [[key, value] as [keyof OpeningHours, string]] : [];
  });

const pad = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600) % 24;
  const minutes = Math.floor((seconds % 3600) / 60);

  return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}`;
};

export const hourRangeBounds = (value: string): { opens: string; closes: string } | undefined => {
  const range = parseTimeRange(value);

  return range ? { opens: pad(range.start), closes: pad(range.end) } : undefined;
};

export const openingHoursDays = (key: keyof OpeningHours): readonly string[] =>
  OPENING_HOURS_RULES.find((rule) => rule.key === key)?.days ?? [];

export const formatHourRange = (value: string): string => {
  const range = parseTimeRange(value);

  if (!range) return value.trim().replace("-", "–");

  if (range.start === range.end) return "24/7";

  return `${pad(range.start)}–${pad(range.end)}`;
};

export const isOpenNow = (
  openingHours: OpeningHours | undefined,
  now: Date = new Date(),
): boolean | undefined => {
  if (!openingHours || openingHoursEntries(openingHours).length === 0) return undefined;

  const moment = warsawMoment(now);

  if (!moment) return undefined;

  const yesterday = rangeForDay(openingHours, previousDay(moment.day));

  if (yesterday && yesterday.end < yesterday.start && moment.secondsOfDay < yesterday.end) 
    return true;

  const today = rangeForDay(openingHours, moment.day);

  if (!today) return false;

  if (today.start === today.end) return true;

  if (today.end < today.start) return moment.secondsOfDay >= today.start;

  return moment.secondsOfDay >= today.start && moment.secondsOfDay < today.end;
};
