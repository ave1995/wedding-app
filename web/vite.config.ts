import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig(({ mode }) => { // Destructure 'mode' from the config object
  // Load env file based on `mode` in the current working directory.
  // Set the third parameter to '' to load all env regardless of the `VITE_` prefix.
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [react(), tailwindcss()],
    // Access the environment variable from the 'env' object loaded by loadEnv
    base: env.VITE_GCS_BASE_URL || '/', // Fallback to '/' for local dev if not set
  };
});
