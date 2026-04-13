import { defineConfig } from 'vite';

export default defineConfig({
  root: '.',
  publicDir: 'public',
  build: {
    outDir: 'dist',
    assetsDir: '.',
    emptyOutDir: true
  },
  server: {
    port: 3000
  }
});
