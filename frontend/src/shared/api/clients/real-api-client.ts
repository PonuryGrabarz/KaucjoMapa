import * as requests from "../requests";
import * as responses from "../responses";

import ky from "ky";

import { z } from "zod";
import { Result } from "../result";
import { ApiError } from "../errors";
import { API_TIMEOUT_MS, API_URL } from "astro:env/server";
import { ApiClientBase } from "../abstractions/api-client-base";
import { isHTTPError, isNetworkError, isTimeoutError } from "ky";

import type { ApiFieldError } from "../errors";
import type { KyInstance, Options as KyOptions } from "ky";
import type { ApiCallContext } from "../abstractions/api-client-base";

const ENDPOINTS = {
	returnPoints: "return_points",
	reports: "reports",
	votes: "votes",
	bookmarks: "bookmarks",
	login: "login",
	users: "users",
	refreshToken: "auth/token",
	revokeToken: "auth/revoke",
	notifications: "notifications",
} as const;

const REFRESH_COOKIE_NAME = "RefreshToken";

const TOKEN_BODY_SCHEMA = z.object({ token: z.string().min(1) }).loose();

interface RefreshCookie {
	value: string;
	expiresAt?: string;
}

const API_FIELD_ERROR_SCHEMA = z.object({
	field: z.string(),
	rule: z.string(),
	param: z.string().optional(),
	message: z.string(),
});

const API_ERROR_BODY_SCHEMA = z
	.object({
		message: z.string().optional(),
		error: z.string().optional(),
		errors: z.array(API_FIELD_ERROR_SCHEMA).optional(),
		details: z.array(API_FIELD_ERROR_SCHEMA).optional(),
	}).loose();

export interface RealApiClientOptions {
	baseUrl?: string;
	timeoutMs?: number;
	retryLimit?: number;
	fetch?: typeof globalThis.fetch;
}

export class RealApiClient extends ApiClientBase {
	private readonly http: KyInstance;

	public constructor(options: RealApiClientOptions = {}) {
		super();

		const kyOptions: KyOptions = {
			prefix: options.baseUrl ?? API_URL,
			timeout: options.timeoutMs ?? API_TIMEOUT_MS,
			retry: { limit: options.retryLimit ?? 2 },
		};

		if (options.fetch) kyOptions.fetch = options.fetch;

		this.http = ky.create(kyOptions);
	}

	protected async handleGetNotificationsRequest(_request: requests.GetNotificationsRequest, context: ApiCallContext):
    Promise<Result<responses.GetNotificationsResponse>>
  {
		return this.getJson<responses.GetNotificationsResponse>(ENDPOINTS.notifications, context);
	}

	protected async handleCreateNotificationRequest(request: requests.CreateNotificationRequest, context: ApiCallContext):
    Promise<Result<responses.CreateNotificationResponse>>
  {
		return this.postJson<responses.CreateNotificationResponse>(ENDPOINTS.notifications, context, {
			json: request.toBody(),
		});
	}

	protected async handleMarkNotificationSeenRequest(request: requests.MarkNotificationSeenRequest, context: ApiCallContext):
    Promise<Result<responses.MarkNotificationSeenResponse>>
  {
		return this.putJson<responses.MarkNotificationSeenResponse>(
			`${ENDPOINTS.notifications}/${request.notificationId}`,
			context,
		);
	}

	protected async handleDeleteNotificationRequest(request: requests.DeleteNotificationRequest, context: ApiCallContext):
    Promise<Result<responses.DeleteNotificationResponse>>
  {
		return this.sendNoContent(`${ENDPOINTS.notifications}/${request.notificationId}`, context);
	}

	protected async handleLoginRequest(request: requests.LoginRequest, context: ApiCallContext):
    Promise<Result<responses.LoginResponse>>
  {
		return this.postForTokens(ENDPOINTS.login, context, request.toBody());
	}

	protected async handleRegisterRequest(request: requests.RegisterRequest, context: ApiCallContext):
    Promise<Result<responses.RegisterResponse>>
  {
		return this.postForTokens(ENDPOINTS.users, context, request.toBody());
	}

	protected async handleRefreshTokenRequest(request: requests.RefreshTokenRequest, context: ApiCallContext):
    Promise<Result<responses.RefreshTokenResponse>>
  {
		try {
			const payload = await this.http
				.post(ENDPOINTS.refreshToken, {
					headers: this.withRefreshCookie(context.headers, request.refreshToken),
					signal: context.signal,
				}).json<unknown>();

			const parsed = TOKEN_BODY_SCHEMA.safeParse(payload);

			if (!parsed.success) return Result.failure(this.invalidTokenBody(ENDPOINTS.refreshToken));

			return Result.success({ accessToken: parsed.data.token });
		} catch (error) {
			return Result.failure<responses.RefreshTokenResponse>(
				this.toApiError(error, ENDPOINTS.refreshToken),
			);
		}
	}

	protected async handleRevokeTokenRequest(request: requests.RevokeTokenRequest, context: ApiCallContext):
    Promise<Result<responses.RevokeTokenResponse>>
  {
		try {
			await this.http.post(ENDPOINTS.revokeToken, {
				headers: this.withRefreshCookie(context.headers, request.refreshToken),
				signal: context.signal,
			});

			return Result.success(undefined);
		} catch (error) {
			return Result.failure<undefined>(this.toApiError(error, ENDPOINTS.revokeToken));
		}
	}

	protected async handleGetReturnPointsRequest(request: requests.GetReturnPointsRequest, context: ApiCallContext):
    Promise<Result<responses.GetReturnPointsResponse>>
  {
		return this.getJson<responses.GetReturnPointsResponse>(ENDPOINTS.returnPoints, context, {
			searchParams: request.toSearchParams(),
		});
	}

	protected async handleCreateCommunityReportRequest(request: requests.CreateCommunityReportRequest, context: ApiCallContext): 
    Promise<Result<responses.CreateCommunityReportResponse>> 
  {
		return this.postJson<responses.CreateCommunityReportResponse>(ENDPOINTS.reports, context, {
			json: request.toBody(),
		});
	}

	protected async handleCreateReturnPointRequest(request: requests.CreateReturnPointRequest, context: ApiCallContext): 
    Promise<Result<responses.ReturnPoint>> 
  {
		return this.postJson<responses.ReturnPoint>(ENDPOINTS.returnPoints, context, {
			json: request.toBody(),
		});
	}

	protected async handleGetCommunityReportsRequest(request: requests.GetCommunityReportsRequest, context: ApiCallContext): 
    Promise<Result<responses.GetCommunityReportsResponse>> 
  {
		return this.getJson<responses.GetCommunityReportsResponse>(
			`${ENDPOINTS.reports}/${request.returnPointId}`,
			context,
		);
	}

	protected async handleCreateCommunityReportVoteRequest(request: requests.CreateCommunityReportVoteRequest, context: ApiCallContext): 
    Promise<Result<responses.CreateCommunityReportVoteResponse>> 
  {
		return this.postJson<responses.CreateCommunityReportVoteResponse>(
			`${ENDPOINTS.votes}/${request.reportId}`,
			context,
			{ json: request.toBody() },
		);
	}

	protected async handleCreateBookmarkRequest(request: requests.CreateBookmarkRequest, context: ApiCallContext): 
    Promise<Result<responses.CreateBookmarkResponse>> 
  {
		return this.postJson<responses.CreateBookmarkResponse>(ENDPOINTS.bookmarks, context, {
			json: request.toBody(),
		});
	}

	protected async handleGetBookmarkByReturnPointRequest(request: requests.GetBookmarkByReturnPointRequest, context: ApiCallContext): 
    Promise<Result<responses.GetBookmarkByReturnPointResponse>> 
  {
		return this.getJson<responses.GetBookmarkByReturnPointResponse>(
			`${ENDPOINTS.bookmarks}/return_point/${request.returnPointId}`,
			context,
		);
	}

	protected async handleDeleteBookmarkByReturnPointRequest(request: requests.DeleteBookmarkByReturnPointRequest, context: ApiCallContext): 
    Promise<Result<responses.DeleteBookmarkByReturnPointResponse>> 
  {
		return this.sendNoContent(
			`${ENDPOINTS.bookmarks}/return_point/${request.returnPointId}`,
			context,
		);
	}

	private async postForTokens(path: string, context: ApiCallContext, body: Record<string, string>): Promise<Result<responses.AuthTokens>> {
		try {
			const response = await this.http.post(path, {
				json: body,
				headers: context.headers,
				signal: context.signal,
			});

			const parsed = TOKEN_BODY_SCHEMA.safeParse(await response.json<unknown>());

			if (!parsed.success) return Result.failure(this.invalidTokenBody(path));

			const refreshCookie = this.readRefreshCookie(response);

			if (!refreshCookie) {
				return Result.failure<responses.AuthTokens>(
					new ApiError({
						message: `Odpowiedź API dla ${path} nie zawierała ciasteczka ${REFRESH_COOKIE_NAME}.`,
						kind: "invalid_response",
					}),
				);
			}

			return Result.success({
				accessToken: parsed.data.token,
				refreshToken: refreshCookie.value,
				refreshTokenExpiresAt: refreshCookie.expiresAt,
			});
		} catch (error) {
			return Result.failure<responses.AuthTokens>(this.toApiError(error, path));
		}
	}

	private withRefreshCookie(headers: Record<string, string>, refreshToken: string): Record<string, string> {
		return { ...headers, Cookie: `${REFRESH_COOKIE_NAME}=${refreshToken}` };
	}

	private readRefreshCookie(response: Response): RefreshCookie | undefined {
		const rawCookies = typeof response.headers.getSetCookie === "function"
			? response.headers.getSetCookie()
			: [response.headers.get("set-cookie") ?? ""];

		for (const rawCookie of rawCookies) {
			const [pair, ...attributes] = rawCookie.split(";");

			const separator = pair?.indexOf("=") ?? -1;

			if (!pair || separator < 0) continue;

			if (pair.slice(0, separator).trim() !== REFRESH_COOKIE_NAME) continue;

			const value = pair.slice(separator + 1).trim();

			if (!value) continue;

			const expires = attributes
				.map((attribute) => attribute.trim())
				.find((attribute) => attribute.toLowerCase().startsWith("expires="));

			return { value, expiresAt: expires?.slice("expires=".length) };
		}

		return undefined;
	}

	private invalidTokenBody(path: string): ApiError {
		return new ApiError({
			message: `Odpowiedź API dla ${path} nie zawierała pola token.`,
			kind: "invalid_response",
		});
	}

	private async sendNoContent(path: string, context: ApiCallContext): Promise<Result<undefined>> {
		try {
			await this.http.delete(path, { headers: context.headers, signal: context.signal });

			return Result.success(undefined);
		} catch (error) {
			return Result.failure<undefined>(this.toApiError(error, path));
		}
	}

	private async postJson<TResponse>(path: string, context: ApiCallContext, options: KyOptions = {}): Promise<Result<TResponse>> {
		try {
			const payload = await this.http
				.post(path, {
					...options,
					headers: context.headers,
					signal: context.signal,
				}).json<TResponse>();

			return Result.success(payload);
		} catch (error) {
			return Result.failure<TResponse>(this.toApiError(error, path));
		}
	}

	private async putJson<TResponse>(path: string, context: ApiCallContext, options: KyOptions = {}): Promise<Result<TResponse>> {
		try {
			const payload = await this.http
				.put(path, {
					...options,
					headers: context.headers,
					signal: context.signal,
				}).json<TResponse>();

			return Result.success(payload);
		} catch (error) {
			return Result.failure<TResponse>(this.toApiError(error, path));
		}
	}

	private async getJson<TResponse>(path: string, context: ApiCallContext, options: KyOptions = {}): Promise<Result<TResponse>> {
		try {
			const payload = await this.http
				.get(path, {
					...options,
					headers: context.headers,
					signal: context.signal,
				}).json<TResponse>();
        
			return Result.success(payload);
		} catch (error) {
			return Result.failure<TResponse>(this.toApiError(error, path));
		}
	}

	private toApiError(error: unknown, path: string): ApiError {
		if (isHTTPError(error)) {
			const body = this.parseErrorBody(error.data);

			return new ApiError({
				message:
					body?.message
					?? body?.error
					?? `Serwer odpowiedział błędem ${error.response.status} dla ${path}.`,
				kind: "http",
				status: error.response.status,
				fieldErrors: body?.errors ?? body?.details ?? [],
				cause: error,
			});
		}

		if (isTimeoutError(error)) {
			return new ApiError({
				message: `Przekroczono czas oczekiwania na odpowiedź dla ${path}.`,
				kind: "timeout",
				cause: error,
			});
		}

		if (isNetworkError(error)) {
			return new ApiError({
				message: `Nie udało się połączyć z API dla ${path}.`,
				kind: "network",
				cause: error,
			});
		}

		return new ApiError({
			message: `Odpowiedź API dla ${path} miała nieoczekiwaną postać.`,
			kind: "invalid_response",
			cause: error,
		});
	}

	private parseErrorBody(data: unknown): { message?: string; error?: string; errors?: ApiFieldError[]; details?: ApiFieldError[] } | undefined {
		const parsed = API_ERROR_BODY_SCHEMA.safeParse(data);

		return parsed.success ? parsed.data : undefined;
	}
}
