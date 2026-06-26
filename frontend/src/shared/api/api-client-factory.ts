import { API_URL, USE_DUMMY_API } from "astro:env/client";

import type ApiClientBase from "./abstractions/api-client-base";

let instance: ApiClientBase | undefined;

const isProd = import.meta.env.PROD;

if (isProd) {
  instance = await import("./implementations/api-client").then((module) => new module.ApiClient(API_URL as string));
} else {
  if (USE_DUMMY_API) {
    instance = await import("./implementations/dummy-api-client").then((module) => new module.DummyApiClient());
  }
}

export const apiClient = instance!;
