import { defineMiddleware } from "astro:middleware";

const AUTH_PATHS = new Set(["/logowanie", "/rejestracja"]);

const normalizePathname = (pathname: string): string => {
  if (pathname === "/") return pathname;
  return pathname.replace(/\/+$/u, "");
};

export const onRequest = defineMiddleware(async (context, next) => {
  if (AUTH_PATHS.has(normalizePathname(context.url.pathname))) {
    const user = await context.session?.get("user");

    if (user) return context.redirect("/");
  }

  return await next();
});
