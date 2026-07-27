import { defineMiddleware } from "astro:middleware";

const GUEST_ONLY_PATHS = new Set(["/logowanie", "/rejestracja"]);
const AUTHENTICATED_PATHS = new Set(["/profil"]);

const normalizePathname = (pathname: string): string => {
  if (pathname === "/") return pathname;
  return pathname.replace(/\/+$/u, "");
};

export const onRequest = defineMiddleware(async (context, next) => {
  const pathname = normalizePathname(context.url.pathname);
  const isGuestOnlyPath = GUEST_ONLY_PATHS.has(pathname);
  const isAuthenticatedPath = AUTHENTICATED_PATHS.has(pathname);

  if (!isGuestOnlyPath && !isAuthenticatedPath) return await next();

  const user = await context.session?.get("user");

  if (isGuestOnlyPath && user) return context.redirect("/");
  if (isAuthenticatedPath && !user) return context.redirect("/logowanie");

  return await next();
});
