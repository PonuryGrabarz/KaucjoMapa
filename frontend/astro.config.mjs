// @ts-check
import { defineConfig, envField, fontProviders } from "astro/config";

import node from "@astrojs/node";
import svelte from "@astrojs/svelte";
import sitemap from "@astrojs/sitemap";

export default defineConfig({
  site: process.env.SITE_URL ?? "http://localhost:4321",

	output: "server",
	compressHTML: true,

  integrations: [svelte(), sitemap()],

  server: {
    host: true,
    port: 4321,
    allowedHosts: ["kaucjo-mapa.pl"],
  },

    security: {
      allowedDomains: [
        {
          protocol: "https",
          hostname: "kaucjo-mapa.pl",
        },
      ],
  },

	adapter: node({
		mode: "standalone",
	}),

	env: {
		schema: {
      API_URL: envField.string({
        access: "secret",
        context: "server",
        url: true,
      }),
			API_TIMEOUT_MS: envField.number({
				access: "secret",
				context: "server",
				default: 8000,
			}),
      MOCK_API: envField.boolean({
				access: "secret",
				context: "server",
				default: false,
			}),
		},
	},

  fonts: [{
    provider: fontProviders.google(),
    name: "Inter",
    cssVariable: "--font-ui",
    weights: [400, 600, 700, 800],
    subsets: ["latin", "latin-ext"],
  }],

	vite: {
		optimizeDeps: {
			include: ["emittery"],
			exclude: ["maplibre-gl"],
		},
	},
});
