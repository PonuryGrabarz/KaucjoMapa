export interface HistorySummaryItem {
    value: string;
    label: string;
}

export interface HistoryStatItem {
    icon: string;
    label: string;
    value: string;
}

export interface HistoryActivityItem {
    icon: string;
    title: string;
    timestamp: string;
    isoTimestamp: string;
    pills: string[];
    amount: string;
}

export const historySummaryItems: HistorySummaryItem[] = [
    { value: "248", label: "opakowań łącznie" },
    { value: "124 zł", label: "odzyskanej kaucji" },
    { value: "97%", label: "ocena zaufania" },
    { value: "5", label: "seria zgodnych zgłoszeń" },
];

export const historyStatsItems: HistoryStatItem[] = [
    {
        icon: "mdi:map-marker-check-outline",
        label: "Najczęściej odwiedzany punkt",
        value: "RecyHub przy Rynku - 123 razy",
    },
    {
        icon: "mdi:bottle-soda-classic-outline",
        label: "Oddanych butelek PET",
        value: "123 butelek",
    },
    {
        icon: "mdi:cup-outline",
        label: "Oddanych puszek",
        value: "123 puszek",
    },
    {
        icon: "mdi:bottle-wine-outline",
        label: "Oddanych butelek szklanych",
        value: "123 butelek",
    },
    {
        icon: "mdi:shield-check-outline",
        label: "Wskaźnik zaufania",
        value: "97% (9/10)",
    },
    {
        icon: "mdi:cash-multiple",
        label: "Średnia kwota za zwrot",
        value: "10 zł",
    },
];

export const historyActivityItems: HistoryActivityItem[] = Array.from(
    { length: 9 },
    () => ({
        icon: "mdi:check-decagram-outline",
        title: "RecyHub przy Rynku",
        timestamp: "Dziś, 9:24",
        isoTimestamp: "2026-07-21T09:24:00+02:00",
        pills: ["PET: 10", "Puszki: 7"],
        amount: "+10 zł",
    }),
);
