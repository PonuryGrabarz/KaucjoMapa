import { ActionError } from "astro:actions";

import { isApiError } from "@shared/api/errors";

import type { Result } from "@shared/api/result";
import type { ActionErrorCode } from "astro:actions";

const HTTP_STATUS_TO_ACTION_CODE: Record<number, ActionErrorCode> = {
  400: "BAD_REQUEST",
  401: "UNAUTHORIZED",
  403: "FORBIDDEN",
  404: "NOT_FOUND",
  409: "CONFLICT",
  422: "UNPROCESSABLE_CONTENT",
  429: "TOO_MANY_REQUESTS",
};

export const describeErrors = (errors: readonly Error[]): string =>
  errors.map((error) => error.message).join(" ");

export const toActionError = (errors: readonly Error[]): ActionError => {
  const [firstError] = errors;

  const message = describeErrors(errors);

  if (!isApiError(firstError)) return new ActionError({ code: "BAD_GATEWAY", message });

  if (firstError.kind === "timeout") return new ActionError({ code: "GATEWAY_TIMEOUT", message });

  if (firstError.kind === "network")
    return new ActionError({ code: "SERVICE_UNAVAILABLE", message });

  const code = firstError.status ? HTTP_STATUS_TO_ACTION_CODE[firstError.status] : undefined;

  return new ActionError({ code: code ?? "BAD_GATEWAY", message });
};

export const unwrapRequest = <TValue>(result: Result<TValue>): TValue => {
  if (result.isError)
    throw new ActionError({ code: "BAD_REQUEST", message: describeErrors(result.errors) });

  return result.value;
};

export const unwrapApiResult = <TValue>(result: Result<TValue>): TValue => {
  if (result.isError) throw toActionError(result.errors);

  return result.value;
};
