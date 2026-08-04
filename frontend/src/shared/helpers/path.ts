export const normalizePathname = (pathname: string): string => pathname.replace(/\/$/, "") || "/";
