import { Result } from "../result";
import { ApiError, isApiError } from "../errors";
import { apiClient } from "../api-client-factory";
import { LoginRequest, RegisterRequest, RefreshTokenRequest, RevokeTokenRequest } from "../requests";
import { decodeAccessToken, isAccessTokenExpired, toSessionUser } from "./access-token";

import {
  clearSessionCookies,
  readAccessToken,
  readRefreshToken,
  writeAccessToken,
  writeRefreshToken,
} from "./session-cookies";

import type { AuthTokens } from "../responses";
import type { LoginParams, RegisterParams } from "../requests";
import type { SessionCookieContext } from "./session-cookies";

export type SessionContext = SessionCookieContext;

const isCredentialsRejected = (errors: readonly Error[]): boolean => {
  const [error] = errors;

  if (!isApiError(error)) return false;

  return error.kind === "http" && (error.status === 401 || error.status === 403);
};

const buildSession = (accessToken: string): App.Session | null => {
  const claims = decodeAccessToken(accessToken);

  return claims ? { user: toSessionUser(claims), accessToken } : null;
};

const persistTokens = (context: SessionContext, tokens: AuthTokens): App.Session | null => {
  writeRefreshToken(context, tokens.refreshToken, tokens.refreshTokenExpiresAt);
  writeAccessToken(context, tokens.accessToken, tokens.refreshTokenExpiresAt);

  return buildSession(tokens.accessToken);
};

const establishSession = (
  context: SessionContext,
  tokens: AuthTokens,
): Result<App.Session> => {
  const session = persistTokens(context, tokens);

  if (session) return Result.success(session);

  clearSessionCookies(context);

  return Result.failure<App.Session>(
    new ApiError({
      message: "Serwer zwrócił token, którego nie udało się odczytać.",
      kind: "invalid_response",
    }),
  );
};

const refreshAccessToken = async (
  context: SessionContext,
  accessToken: string,
  refreshToken: string,
): Promise<Result<string>> => {
  const request = RefreshTokenRequest.create({ refresh_token: refreshToken });

  if (request.isError) return Result.failure<string>(request.errors);

  const result = await apiClient.send(request.value, { token: accessToken });

  if (result.isError) {
    if (isCredentialsRejected(result.errors)) clearSessionCookies(context);

    return Result.failure<string>(result.errors);
  }

  writeAccessToken(context, result.value.accessToken);

  return Result.success(result.value.accessToken);
};

export const resolveSession = async (context: SessionContext): Promise<App.Session | null> => {
  const refreshToken = readRefreshToken(context.cookies);
  const accessToken = readAccessToken(context.cookies);

  if (!refreshToken || !accessToken) {
    if (refreshToken || accessToken) clearSessionCookies(context);

    return null;
  }

  const claims = decodeAccessToken(accessToken);

  if (!claims) {
    clearSessionCookies(context);

    return null;
  }

  if (!isAccessTokenExpired(claims)) return { user: toSessionUser(claims), accessToken };

  const refreshed = await refreshAccessToken(context, accessToken, refreshToken);

  return refreshed.isError ? null : buildSession(refreshed.value);
};

export const adoptSession = (context: SessionContext, tokens: AuthTokens): Result<App.Session> =>
  establishSession(context, tokens);

export const login = async (
  context: SessionContext,
  params: LoginParams,
): Promise<Result<App.Session>> => {
  const request = LoginRequest.create(params);

  if (request.isError) return Result.failure<App.Session>(request.errors);

  const result = await apiClient.send(request.value);

  if (result.isError) return Result.failure<App.Session>(result.errors);

  return establishSession(context, result.value);
};

export const register = async (
  context: SessionContext,
  params: RegisterParams,
): Promise<Result<App.Session>> => {
  const request = RegisterRequest.create(params);

  if (request.isError) return Result.failure<App.Session>(request.errors);

  const result = await apiClient.send(request.value);

  if (result.isError) return Result.failure<App.Session>(result.errors);

  return establishSession(context, result.value);
};

export const logout = async (context: SessionContext): Promise<Result<undefined>> => {
  const refreshToken = readRefreshToken(context.cookies);
  const accessToken = readAccessToken(context.cookies);

  if (!refreshToken || !accessToken) {
    clearSessionCookies(context);

    return Result.success(undefined);
  }

  const claims = decodeAccessToken(accessToken);

  let usableToken = accessToken;

  if (!claims || isAccessTokenExpired(claims)) {
    const refreshed = await refreshAccessToken(context, accessToken, refreshToken);

    if (refreshed.isError) {
      clearSessionCookies(context);

      return Result.success(undefined);
    }

    usableToken = refreshed.value;
  }

  const request = RevokeTokenRequest.create({ refresh_token: refreshToken });

  clearSessionCookies(context);

  if (request.isError) return Result.failure<undefined>(request.errors);

  const result = await apiClient.send(request.value, { token: usableToken });

  return result.isError ? Result.failure<undefined>(result.errors) : Result.success(undefined);
};
