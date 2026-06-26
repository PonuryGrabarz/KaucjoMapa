// @ts-check
import { defineConfig, envField } from 'astro/config';

// https://astro.build/config
export default defineConfig({
  env: {
    validateSecrets: true,
    schema: {
      API_URL: envField.string({ context: 'client', access: 'public', optional: true }),
      USE_DUMMY_API: envField.boolean({ context: 'client', access: 'public', default: false}),
    }
  }
});
