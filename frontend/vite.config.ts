import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://127.0.0.1:317",
        changeOrigin: true,
        ws: true,
        rewrite: path => path.replace(/^\/api/, "/api"),
      },
      "/vela": {
        target: "http://127.0.0.1:317",
        changeOrigin: true,
        ws: true,
        rewrite: path => path.replace(/^\/api/, "/api"),
      },
      "/z-inject": {
        target: "http://127.0.0.1:317",
        changeOrigin: true,
        ws: true,
        rewrite: path => path.replace(/^\/z-inject/, "/z-inject"),
      },
    },
    
  },
});