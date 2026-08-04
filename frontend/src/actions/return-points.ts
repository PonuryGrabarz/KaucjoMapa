import { z } from "astro/zod";
import { ActionError, defineAction } from "astro:actions";
import { apiClient } from "@shared/api/api-client-factory";
import { CreateReturnPointRequest, GetReturnPointsRequest } from "@shared/api/requests";
import { unwrapApiResult, unwrapRequest } from "@shared/helpers/action-errors";

import type { ApiSendOptions } from "@shared/api/abstractions/api-client-base";

import {
	RETURN_POINT_REFUND_METHODS,
	RETURN_POINT_STATUSES,
	RETURN_POINT_TYPES,
} from "@shared/api/responses";

const LIST_INPUT_SCHEMA = z.object({
	type: z.enum(RETURN_POINT_TYPES).optional(),
	status: z.enum(RETURN_POINT_STATUSES).optional(),
	query: z.string().optional(),
	min_distance: z.number().optional(),
	max_distance: z.number().optional(),
	lat: z.number().optional(),
	lon: z.number().optional(),
	pet: z.boolean().optional(),
	can: z.boolean().optional(),
	glass: z.boolean().optional(),
	refund_method: z.array(z.enum(RETURN_POINT_REFUND_METHODS)).optional(),
	page: z.number().optional(),
	page_size: z.number().optional(),

	onlyMine: z.boolean().optional(),
});

const FETCH_ALL_PAGE_SIZE = 20_000;

type ListFilters = Omit<z.infer<typeof LIST_INPUT_SCHEMA>, "onlyMine">;

const fetchPage = async (
	filters: ListFilters,
	options: ApiSendOptions,
	page: number,
	pageSize: number,
) => {
	const request = unwrapRequest(
		GetReturnPointsRequest.create({ ...filters, page, page_size: pageSize }),
	);

	return unwrapApiResult(await apiClient.send(request, options));
};

const CREATE_INPUT_SCHEMA = z.object({
	name: z.string(),
	type: z.enum(RETURN_POINT_TYPES),
	address: z.string().optional(),
	latitude: z.number(),
	longitude: z.number(),
	pet: z.boolean().optional(),
	can: z.boolean().optional(),
	glass: z.boolean().optional(),
	refund_method: z.array(z.enum(RETURN_POINT_REFUND_METHODS)),
	opening_hours: z.record(z.string(), z.string()).optional(),
	description: z.string().optional(),
	current_status: z.enum(RETURN_POINT_STATUSES),
});

export const returnPoints = {
	create: defineAction({
		input: CREATE_INPUT_SCHEMA,
		handler: async (input, context) => {
			const session = context.locals.session;

			if (!session) {
				throw new ActionError({
					code: "UNAUTHORIZED",
					message: "Wymagane logowanie, aby dodać punkt.",
				});
			}

			const request = unwrapRequest(CreateReturnPointRequest.create(input));

			return unwrapApiResult(await apiClient.send(request, { token: session.accessToken }));
		},
	}),

	list: defineAction({
		input: LIST_INPUT_SCHEMA,
		handler: async (input, context) => {
			const { onlyMine = false, ...filters } = input;

			const session = context.locals.session;

			if (onlyMine && !session) {
				throw new ActionError({
					code: "UNAUTHORIZED",
					message: "Wymagane logowanie, aby zobaczyć własne punkty.",
				});
			}

			const options: ApiSendOptions =
				onlyMine && session ? { token: session.accessToken } : {};

			const wantsEverything = filters.page === undefined && filters.page_size === undefined;

			if (!wantsEverything) {
				const request = unwrapRequest(GetReturnPointsRequest.create(filters));

				return unwrapApiResult(await apiClient.send(request, options));
			}

			const firstPage = await fetchPage(filters, options, 1, FETCH_ALL_PAGE_SIZE);

			const returnPoints = [...firstPage.return_points];

			for (let page = 2; page <= firstPage.pagination.total_pages; page += 1) {
				const nextPage = await fetchPage(filters, options, page, FETCH_ALL_PAGE_SIZE);

				returnPoints.push(...nextPage.return_points);
			}

			return { pagination: firstPage.pagination, return_points: returnPoints };
		},
	}),
};
