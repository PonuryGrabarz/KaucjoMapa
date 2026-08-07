export const mainItems = [
	{ href: "/", label: "Mapa", icon: "material-symbols:map-outline-rounded" as const, activeIcon: "material-symbols:map-rounded" as const },
	{
		href: "/punkty",
		label: "Punkty",
		icon: "material-symbols:location-on-outline-rounded" as const,
		activeIcon: "material-symbols:location-on-rounded" as const,
	},
	{
		href: "/kalkulator",
		label: "Kalkulator",
		icon: "material-symbols:calculate-outline-rounded" as const,
		activeIcon: "material-symbols:calculate-rounded" as const,
	},
];

export const accountItems = [
	{ href: "/historia", label: "Historia", icon: "material-symbols:history-rounded" as const, activeIcon: "material-symbols:history-rounded" as const },
	{
		href: "/powiadomienia",
		label: "Powiadomienia",
		icon: "material-symbols:notifications-outline-rounded" as const,
		activeIcon: "material-symbols:notifications-rounded" as const,
	},
	{
		href: "/profil",
		label: "Profil",
		icon: "material-symbols:person-outline-rounded" as const,
		activeIcon: "material-symbols:person-rounded" as const,
	},
];

export const guestItems = [
	{
		href: "/logowanie",
		label: "Logowanie",
		icon: "material-symbols:login-rounded" as const,
		activeIcon: "material-symbols:login-rounded" as const,
		activeFor: ["/rejestracja"] as readonly string[],
	},
];

const matchesSection = (pathname: string, href: string): boolean =>
	href === "/" ? pathname === "/" : pathname === href || pathname.startsWith(`${href}/`);

export const isPathActive = (
	pathname: string,
	href: string,
	activeFor: readonly string[] = [],
): boolean =>
	matchesSection(pathname, href) || activeFor.some((alias) => matchesSection(pathname, alias));
