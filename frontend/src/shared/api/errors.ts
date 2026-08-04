export interface ApiFieldError {
	field: string;
	rule: string;
	param?: string;
	message: string;
}

export class ApiError extends Error {
	public readonly status?: number;

	public readonly fieldErrors: readonly ApiFieldError[];

	public readonly kind: "http" | "network" | "timeout" | "invalid_response";

	public constructor(options: {
		message: string;
		kind: ApiError["kind"];
		status?: number;
		fieldErrors?: readonly ApiFieldError[];
		cause?: unknown;
	}) {
		super(options.message, { cause: options.cause });

		this.name = "ApiError";
		this.kind = options.kind;
		this.status = options.status;
		this.fieldErrors = options.fieldErrors ?? [];
	}
}

export const isApiError = (error: unknown): error is ApiError => error instanceof ApiError;

const isAbortLikeError = (error: unknown): boolean =>
	error instanceof Error && (error.name === "AbortError" || error.name === "TimeoutError");

export const toApiError = (error: unknown): ApiError => {
	if (isApiError(error)) return error;

	if (isAbortLikeError(error)) {
		return new ApiError({
			message: "Żądanie zostało przerwane.",
			kind: "timeout",
			cause: error,
		});
	}

	return new ApiError({
		message:
			error instanceof Error
				? error.message
				: "Wystąpił nieznany błąd podczas komunikacji z API.",
		kind: "network",
		cause: error,
	});
};
