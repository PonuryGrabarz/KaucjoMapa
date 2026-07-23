export interface AlertItem {
    id: string;
    icon: string;
    title: string;
    description: string;
    timestamp: string;
    dateTime: string;
    unread: boolean;
}

const alertTemplate = {
    icon: "mdi:check-circle-outline",
    title: "Butelkomat Dworzec PKP znów działa",
    description: "Punkt jest dostępny 24/7 i ma krótki czas oczekiwania.",
    timestamp: "2 min temu",
    dateTime: "2026-07-22T08:56:00+02:00",
} satisfies Omit<AlertItem, "id" | "unread">;

export const alertItems: AlertItem[] = Array.from({ length: 17 }, (_, index) => ({
    ...alertTemplate,
    id: `return-point-active-${index + 1}`,
    unread: index === 0,
}));
