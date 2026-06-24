// @ts-check
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({
  vite: {
    server: {
      allowedHosts: ["f10c-2a00-f41-587a-e764-58a2-99ff-fef6-14aa.ngrok-free.app"]
    }
  }
});
