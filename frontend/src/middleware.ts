import { defineMiddleware } from "astro:middleware";

import { resolveSession } from "@shared/api/auth";

export const onRequest = defineMiddleware(async (context, next) => {
  context.locals.session = await resolveSession({
    cookies: context.cookies,
    url: context.url,
  });

  return next();
});
