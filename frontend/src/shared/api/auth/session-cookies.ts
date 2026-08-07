import type { AstroCookies, AstroCookieSetOptions } from "astro";

export const ACCESS_COOKIE_NAME = "kaucjomapa.access";
export const REFRESH_COOKIE_NAME = "kaucjomapa.refresh";

const REFRESH_MAX_AGE_SECONDS = 60 * 60 * 24 * 60;

export interface SessionCookieContext {
  cookies: AstroCookies;
  url: URL;
}

const baseOptions = (url: URL): AstroCookieSetOptions => ({
  path: "/",
  httpOnly: true,
  sameSite: "lax",
  secure: url.protocol === "https:",
});

const parseExpiry = (value: string | undefined): Date | undefined => {
  if (!value) return undefined;

  const parsed = new Date(value);

  return Number.isNaN(parsed.getTime()) ? undefined : parsed;
};

const readCookie = (cookies: AstroCookies, name: string): string | undefined => {
  const value = cookies.get(name)?.value?.trim();

  return value ? value : undefined;
};

export const readAccessToken = (cookies: AstroCookies): string | undefined =>
  readCookie(cookies, ACCESS_COOKIE_NAME);

export const readRefreshToken = (cookies: AstroCookies): string | undefined =>
  readCookie(cookies, REFRESH_COOKIE_NAME);

export const writeAccessToken = (
  context: SessionCookieContext,
  accessToken: string,
  refreshTokenExpiresAt?: string,
): void => {
  const expires = parseExpiry(refreshTokenExpiresAt);

  context.cookies.set(ACCESS_COOKIE_NAME, accessToken, {
    ...baseOptions(context.url),
    ...(expires ? { expires } : { maxAge: REFRESH_MAX_AGE_SECONDS }),
  });
};

export const writeRefreshToken = (
  context: SessionCookieContext,
  refreshToken: string,
  refreshTokenExpiresAt?: string,
): void => {
  const expires = parseExpiry(refreshTokenExpiresAt);

  context.cookies.set(REFRESH_COOKIE_NAME, refreshToken, {
    ...baseOptions(context.url),
    ...(expires ? { expires } : { maxAge: REFRESH_MAX_AGE_SECONDS }),
  });
};

export const clearSessionCookies = (context: SessionCookieContext): void => {
  const options = baseOptions(context.url);

  context.cookies.delete(ACCESS_COOKIE_NAME, options);
  context.cookies.delete(REFRESH_COOKIE_NAME, options);
};
