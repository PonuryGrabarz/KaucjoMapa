import * as requests from "../requests";
import * as responses from "../responses";

import { Result } from "../result";
import { ApiError, toApiError } from "../errors";

interface ApiResponseMap {
	"login": responses.LoginResponse;
	"get-avatar": responses.GetAvatarResponse;
	"get-notifications": responses.GetNotificationsResponse;
	"create-notification": responses.CreateNotificationResponse;
	"mark-notification-seen": responses.MarkNotificationSeenResponse;
	"delete-notification": responses.DeleteNotificationResponse;
	"register": responses.RegisterResponse;
	"refresh-token": responses.RefreshTokenResponse;
	"revoke-token": responses.RevokeTokenResponse;
	"get-return-points": responses.GetReturnPointsResponse;
	"create-community-report": responses.CreateCommunityReportResponse;
	"create-return-point": responses.ReturnPoint;
	"get-community-reports": responses.GetCommunityReportsResponse;
	"create-community-report-vote": responses.CreateCommunityReportVoteResponse;
	"create-bookmark": responses.CreateBookmarkResponse;
	"get-bookmark-by-return-point": responses.GetBookmarkByReturnPointResponse;
	"delete-bookmark-by-return-point": responses.DeleteBookmarkByReturnPointResponse;
}

type ApiRequest =
	| requests.GetAvatarRequest
	| requests.GetNotificationsRequest
	| requests.CreateNotificationRequest
	| requests.MarkNotificationSeenRequest
	| requests.DeleteNotificationRequest
	| requests.LoginRequest
	| requests.RegisterRequest
	| requests.RefreshTokenRequest
	| requests.RevokeTokenRequest
	| requests.GetReturnPointsRequest
	| requests.CreateCommunityReportRequest
	| requests.CreateReturnPointRequest
	| requests.GetCommunityReportsRequest
	| requests.CreateCommunityReportVoteRequest
	| requests.CreateBookmarkRequest
	| requests.GetBookmarkByReturnPointRequest
	| requests.DeleteBookmarkByReturnPointRequest;

type ApiResponseValue<TRequest extends ApiRequest> = ApiResponseMap[TRequest["kind"]];

type ApiResponse<TRequest extends ApiRequest> = Result<ApiResponseValue<TRequest>>;

export interface ApiSendOptions {
	token?: string;
	signal?: AbortSignal;
}

export interface ApiCallContext {
	headers: Record<string, string>;
	signal?: AbortSignal;
}

export abstract class ApiClientBase {
	protected readonly defaultHeaders: Readonly<Record<string, string>> = {
		Accept: "application/json",
		"Accept-Language": "pl-PL",
	};

	public async send<TRequest extends ApiRequest>(request: TRequest, options: ApiSendOptions = {}): 
    Promise<ApiResponse<TRequest>> 
  {
		try {
			if (options.signal?.aborted) {
				return Result.failure<ApiResponseValue<TRequest>>(
					new ApiError({ message: "Żądanie zostało przerwane.", kind: "timeout" }),
				);
			}

			const context = this.createCallContext(options);

			switch (request.kind) {
				case "get-avatar":
					return (await this.handleGetAvatarRequest(request, context)) as ApiResponse<TRequest>;
				case "get-notifications":
					return (await this.handleGetNotificationsRequest(request, context)) as ApiResponse<TRequest>;
				case "create-notification":
					return (await this.handleCreateNotificationRequest(request, context)) as ApiResponse<TRequest>;
				case "mark-notification-seen":
					return (await this.handleMarkNotificationSeenRequest(request, context)) as ApiResponse<TRequest>;
				case "delete-notification":
					return (await this.handleDeleteNotificationRequest(request, context)) as ApiResponse<TRequest>;
				case "login":
					return (await this.handleLoginRequest(request, context)) as ApiResponse<TRequest>;
				case "register":
					return (await this.handleRegisterRequest(request, context)) as ApiResponse<TRequest>;
				case "refresh-token":
					return (await this.handleRefreshTokenRequest(request, context)) as ApiResponse<TRequest>;
				case "revoke-token":
					return (await this.handleRevokeTokenRequest(request, context)) as ApiResponse<TRequest>;
				case "get-return-points":
					return (await this.handleGetReturnPointsRequest(request, context)) as ApiResponse<TRequest>;
				case "create-community-report":
					return (await this.handleCreateCommunityReportRequest(request, context)) as ApiResponse<TRequest>;
				case "create-return-point":
					return (await this.handleCreateReturnPointRequest(request, context)) as ApiResponse<TRequest>;
				case "get-community-reports":
					return (await this.handleGetCommunityReportsRequest(request, context)) as ApiResponse<TRequest>;
				case "create-community-report-vote":
					return (await this.handleCreateCommunityReportVoteRequest(request, context)) as ApiResponse<TRequest>;
				case "create-bookmark":
					return (await this.handleCreateBookmarkRequest(request, context)) as ApiResponse<TRequest>;
				case "get-bookmark-by-return-point":
					return (await this.handleGetBookmarkByReturnPointRequest(request, context)) as ApiResponse<TRequest>;
				case "delete-bookmark-by-return-point":
					return (await this.handleDeleteBookmarkByReturnPointRequest(request, context)) as ApiResponse<TRequest>;
				default: {
					const unhandled: never = request;

					return Result.failure<ApiResponseValue<TRequest>>(
						new ApiError({
							message: `Nieobsługiwany typ żądania: ${(unhandled as ApiRequest).kind}`,
							kind: "invalid_response",
						}),
					);
				}
			}
		} catch (error) {
			return Result.failure<ApiResponseValue<TRequest>>(toApiError(error));
		}
	}
  
	protected createCallContext(options: ApiSendOptions): ApiCallContext {
		const headers: Record<string, string> = { ...this.defaultHeaders };

		if (options.token) headers["Authorization"] = `Bearer ${options.token}`;

		return { headers, signal: options.signal };
	}

	protected abstract handleGetAvatarRequest(request: requests.GetAvatarRequest, context: ApiCallContext):
    Promise<Result<responses.GetAvatarResponse>>;

	protected abstract handleGetNotificationsRequest(request: requests.GetNotificationsRequest, context: ApiCallContext):
    Promise<Result<responses.GetNotificationsResponse>>;

	protected abstract handleCreateNotificationRequest(request: requests.CreateNotificationRequest, context: ApiCallContext):
    Promise<Result<responses.CreateNotificationResponse>>;

	protected abstract handleMarkNotificationSeenRequest(request: requests.MarkNotificationSeenRequest, context: ApiCallContext):
    Promise<Result<responses.MarkNotificationSeenResponse>>;

	protected abstract handleDeleteNotificationRequest(request: requests.DeleteNotificationRequest, context: ApiCallContext):
    Promise<Result<responses.DeleteNotificationResponse>>;

	protected abstract handleLoginRequest(request: requests.LoginRequest, context: ApiCallContext):
    Promise<Result<responses.LoginResponse>>;

	protected abstract handleRegisterRequest(request: requests.RegisterRequest, context: ApiCallContext):
    Promise<Result<responses.RegisterResponse>>;

	protected abstract handleRefreshTokenRequest(request: requests.RefreshTokenRequest, context: ApiCallContext):
    Promise<Result<responses.RefreshTokenResponse>>;

	protected abstract handleRevokeTokenRequest(request: requests.RevokeTokenRequest, context: ApiCallContext):
    Promise<Result<responses.RevokeTokenResponse>>;

	protected abstract handleGetReturnPointsRequest(request: requests.GetReturnPointsRequest, context: ApiCallContext):
    Promise<Result<responses.GetReturnPointsResponse>>;

	protected abstract handleCreateCommunityReportRequest(request: requests.CreateCommunityReportRequest, context: ApiCallContext): 
    Promise<Result<responses.CreateCommunityReportResponse>>;

	protected abstract handleCreateReturnPointRequest(request: requests.CreateReturnPointRequest, context: ApiCallContext): 
    Promise<Result<responses.ReturnPoint>>;

	protected abstract handleGetCommunityReportsRequest(request: requests.GetCommunityReportsRequest, context: ApiCallContext): 
    Promise<Result<responses.GetCommunityReportsResponse>>;

	protected abstract handleCreateCommunityReportVoteRequest(request: requests.CreateCommunityReportVoteRequest, context: ApiCallContext): 
    Promise<Result<responses.CreateCommunityReportVoteResponse>>;

	protected abstract handleCreateBookmarkRequest(request: requests.CreateBookmarkRequest, context: ApiCallContext): 
    Promise<Result<responses.CreateBookmarkResponse>>;

	protected abstract handleGetBookmarkByReturnPointRequest(request: requests.GetBookmarkByReturnPointRequest, context: ApiCallContext): 
    Promise<Result<responses.GetBookmarkByReturnPointResponse>>;

	protected abstract handleDeleteBookmarkByReturnPointRequest(request: requests.DeleteBookmarkByReturnPointRequest, context: ApiCallContext): 
    Promise<Result<responses.DeleteBookmarkByReturnPointResponse>>;
}
