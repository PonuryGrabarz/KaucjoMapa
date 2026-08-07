import { z } from "astro/zod";
import { isApiError } from "@shared/api/errors";
import { ActionError, defineAction } from "astro:actions";
import { apiClient } from "@shared/api/api-client-factory";
import { unwrapApiResult, unwrapRequest } from "@shared/helpers/action-errors";

import {
	CreateBookmarkRequest,
	GetBookmarkByReturnPointRequest,
	DeleteBookmarkByReturnPointRequest,
} from "@shared/api/requests";

import type { ApiSendOptions } from "@shared/api/abstractions/api-client-base";

const POINT_INPUT_SCHEMA = z.object({
	returnPointId: z.uuid(),
});

const requireSession = (session: App.Session | null): ApiSendOptions => {
	if (!session) {
		throw new ActionError({
			code: "UNAUTHORIZED",
			message: "Wymagane logowanie, aby zarządzać zapisanymi punktami.",
		});
	}

	return { token: session.accessToken };
}

export const bookmarks = {
	isBookmarked: defineAction({
		input: POINT_INPUT_SCHEMA,
		handler: async (input, context) => {
			const options = requireSession(context.locals.session);

			const request = unwrapRequest(
				GetBookmarkByReturnPointRequest.create({ return_point_id: input.returnPointId }),
			);

			const result = await apiClient.send(request, options);

			if (result.isError) {
				const [error] = result.errors;

				if (isApiError(error) && error.status === 404) return { bookmarked: false };

				return unwrapApiResult(result) as never;
			}

			return { bookmarked: true };
		},
	}),

	add: defineAction({
		input: POINT_INPUT_SCHEMA,
		handler: async (input, context) => {
			const options = requireSession(context.locals.session);

			const request = unwrapRequest(
				CreateBookmarkRequest.create({ return_point_id: input.returnPointId }),
			);

			unwrapApiResult(await apiClient.send(request, options));

			return { bookmarked: true };
		},
	}),

	remove: defineAction({
		input: POINT_INPUT_SCHEMA,
		handler: async (input, context) => {
			const options = requireSession(context.locals.session);

			const request = unwrapRequest(
				DeleteBookmarkByReturnPointRequest.create({ return_point_id: input.returnPointId }),
			);

			unwrapApiResult(await apiClient.send(request, options));

			return { bookmarked: false };
		},
	}),
};
