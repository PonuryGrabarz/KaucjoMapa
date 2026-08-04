import { z } from "astro/zod";
import { ActionError, defineAction } from "astro:actions";
import { apiClient } from "@shared/api/api-client-factory";
import { COMMUNITY_REPORT_STATUSES, COMMUNITY_REPORT_VOTE_VALUES } from "@shared/api/responses";
import { unwrapApiResult, unwrapRequest } from "@shared/helpers/action-errors";

import {
	CreateCommunityReportRequest,
	GetCommunityReportsRequest,
	CreateCommunityReportVoteRequest,
} from "@shared/api/requests";


const CREATE_INPUT_SCHEMA = z.object({
	returnPointId: z.uuid(),
	status: z.enum(COMMUNITY_REPORT_STATUSES),
	description: z.string().trim().min(1),
});

const LIST_INPUT_SCHEMA = z.object({
	returnPointId: z.uuid(),
});

const VOTE_INPUT_SCHEMA = z.object({
	reportId: z.uuid(),
	vote: z.enum(COMMUNITY_REPORT_VOTE_VALUES),
});

export const communityReports = {
	list: defineAction({
		input: LIST_INPUT_SCHEMA,
		handler: async (input) => {
			const request = unwrapRequest(
				GetCommunityReportsRequest.create({ return_point_id: input.returnPointId }),
			);

			return unwrapApiResult(await apiClient.send(request));
		},
	}),

	vote: defineAction({
		input: VOTE_INPUT_SCHEMA,
		handler: async (input, context) => {
			const session = context.locals.session;

			if (!session) {
				throw new ActionError({
					code: "UNAUTHORIZED",
					message: "Wymagane logowanie, aby ocenić zgłoszenie.",
				});
			}

			const request = unwrapRequest(
				CreateCommunityReportVoteRequest.create({
					report_id: input.reportId,
					vote: input.vote,
				}),
			);

			return unwrapApiResult(await apiClient.send(request, { token: session.accessToken }));
		},
	}),

	create: defineAction({
		input: CREATE_INPUT_SCHEMA,
		handler: async (input, context) => {
			const session = context.locals.session;

			if (!session) {
				throw new ActionError({
					code: "UNAUTHORIZED",
					message: "Wymagane logowanie, aby zgłosić status punktu.",
				});
			}

			const request = unwrapRequest(
				CreateCommunityReportRequest.create({
					return_point_id: input.returnPointId,
					status: input.status,
					description: input.description,
				}),
			);

			return unwrapApiResult(await apiClient.send(request, { token: session.accessToken }));
		},
	}),
};
