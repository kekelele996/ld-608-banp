import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Local dev: proxy /api and /health to the Go backend (container/internal port 3000).
// Production uses nginx.conf for the same reverse proxy.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 20108,
    host: "0.0.0.0",
    proxy: {
      "/api": { target: "http://localhost:21108", changeOrigin: true },
      "/health": { target: "http://localhost:21108", changeOrigin: true }
    }
  }
});
