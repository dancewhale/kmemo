import { defineConfig } from "vite";

// Minimal Vite setup for Wails dev/build. Port is a sane default; Wails can override via devServerUrl.
export default defineConfig({
  server: {
    port: 9245,
    strictPort: true,
  },
});
