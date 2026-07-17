// @ts-check
import { defineConfig, envField } from "astro/config";

import preact from "@astrojs/preact";

import node from "@astrojs/node";

// https://astro.build/config
export default defineConfig({
  integrations: [preact()],

  server: {
    host: true,
    port: 4321,
    allowedHosts: ["kaucjo-mapa.pl"],
  },

  env: {
    validateSecrets: true,
    schema: {
      API_URL: envField.string({ context: "client", access: "public", optional: true }),
      USE_DUMMY_API: envField.boolean({ context: "client", access: "public", default: false}),
      SITE_URL: envField.string({ context: "server", access: "public", optional: true }),
    }
  },

  output: "server",

  adapter: node({
    mode: "standalone",
  }),
});
