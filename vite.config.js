import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    tailwindcss()
  ],
  build: {
    manifest: true,
    emptyOutDir: false,
    rollupOptions: {
      input: "assets/app.js",
    },
  }
});
