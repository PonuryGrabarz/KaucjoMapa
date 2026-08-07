import { adoptSession } from "./session-service";

import type { SessionContext } from "./session-service";

export const HANDOFF_ACCESS_COOKIE = "OAuthAccess";
export const HANDOFF_REFRESH_COOKIE = "OAuthRefresh";

const readHandoff = (context: SessionContext, name: string): string | undefined => {
  const value = context.cookies.get(name)?.value?.trim();

  return value ? value : undefined;
};

export const completeOAuthHandoff = (
  context: SessionContext,
  currentSession: App.Session | null,
): boolean => {
  const accessToken = readHandoff(context, HANDOFF_ACCESS_COOKIE);
  const refreshToken = readHandoff(context, HANDOFF_REFRESH_COOKIE);

  context.cookies.delete(HANDOFF_ACCESS_COOKIE, { path: "/" });
  context.cookies.delete(HANDOFF_REFRESH_COOKIE, { path: "/" });

  if (currentSession) return true;

  if (!accessToken || !refreshToken) return false;

  return !adoptSession(context, { accessToken, refreshToken }).isError;
};
