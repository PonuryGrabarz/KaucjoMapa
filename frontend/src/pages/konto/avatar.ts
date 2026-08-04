import { apiClient } from "@shared/api/api-client-factory";
import { GetAvatarRequest } from "@shared/api/requests";

import type { APIRoute } from "astro";

export const GET: APIRoute = async (context) => {
	const session = context.locals.session;

	if (!session) return new Response(null, { status: 401 });

	const request = GetAvatarRequest.create();

	if (request.isError) return new Response(null, { status: 500 });

	const result = await apiClient.send(request.value, { token: session.accessToken });

	if (result.isError) return new Response(null, { status: 404 });

	return new Response(result.value.data, {
		headers: {
			"Content-Type": result.value.contentType,
			"Cache-Control": "private, max-age=60",
		},
	});
};
