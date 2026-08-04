import { MOCK_API } from "astro:env/server";

import type { ApiClientBase } from "./abstractions/api-client-base";

let instance: ApiClientBase;

if (MOCK_API) {
  // TODO: Uncomment this when the mock API client is implemented
	// instance = await import("./implementations/mock-api-client").then(
	// 	(module) => new module.MockApiClient(),
	// );
} else {
	instance = await import("./clients/real-api-client").then(
		(module) => new module.RealApiClient(),
	);
}

export const apiClient = instance!;
