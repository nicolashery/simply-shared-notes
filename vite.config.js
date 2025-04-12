import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    manifest: true,
    emptyOutDir: false,
    rollupOptions: {
      input: "assets/app.js",
    },
  }
});
