import { z } from "astro/zod";
import { ActionError, defineAction } from "astro:actions";
import { apiClient } from "@shared/api/api-client-factory";
import { unwrapApiResult, unwrapRequest } from "@shared/helpers/action-errors";

import {
	GetNotificationsRequest,
	MarkNotificationSeenRequest,
	DeleteNotificationRequest,
} from "@shared/api/requests";

import type { ApiSendOptions } from "@shared/api/abstractions/api-client-base";

const NOTIFICATION_INPUT_SCHEMA = z.object({
	notificationId: z.uuid(),
});

const requireSession = (session: App.Session | null): ApiSendOptions => {
	if (!session) {
		throw new ActionError({
			code: "UNAUTHORIZED",
			message: "Wymagane logowanie, aby zarządzać powiadomieniami.",
		});
	}

	return { token: session.accessToken };
};

export const notifications = {
	list: defineAction({
		handler: async (_input, context) => {
			const options = requireSession(context.locals.session);

			const request = unwrapRequest(GetNotificationsRequest.create());

			const notificationList = unwrapApiResult(await apiClient.send(request, options));

			return {
				notifications: notificationList,
				unreadCount: notificationList.filter((notification) => !notification.seen).length,
			};
		},
	}),

	markSeen: defineAction({
		accept: "form",
		input: NOTIFICATION_INPUT_SCHEMA,
		handler: async (input, context) => {
			const options = requireSession(context.locals.session);

			const request = unwrapRequest(
				MarkNotificationSeenRequest.create({ notification_id: input.notificationId }),
			);

			return unwrapApiResult(await apiClient.send(request, options));
		},
	}),

	remove: defineAction({
		accept: "form",
		input: NOTIFICATION_INPUT_SCHEMA,
		handler: async (input, context) => {
			const options = requireSession(context.locals.session);

			const request = unwrapRequest(
				DeleteNotificationRequest.create({ notification_id: input.notificationId }),
			);

			unwrapApiResult(await apiClient.send(request, options));

			return { deleted: true };
		},
	}),
};
