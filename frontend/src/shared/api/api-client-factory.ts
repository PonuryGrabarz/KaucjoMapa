import { API_URL, USE_DUMMY_API } from 'astro:env/client';

import type ApiClientBase from './abstractions/api-client-base';

let instance: ApiClientBase;

const useDummyApi = !import.meta.env.PROD && USE_DUMMY_API;

if (useDummyApi) {
  instance = await import('./implementations/dummy-api-client').then(
    (module) => new module.DummyApiClient(),
  );
} else {
  const apiUrl = API_URL;

  if (!apiUrl) {
    throw new Error('API_URL must be set when using the real API client.');
  }

  instance = await import('./implementations/api-client').then(
    (module) => new module.ApiClient(apiUrl),
  );
}

export const apiClient = instance;
