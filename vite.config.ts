import { defineConfig } from 'vite';
import { nodePolyfills } from "vite-plugin-node-polyfills";


export default defineConfig({
  define: {
    global: 'globalThis',
    process: 'process',
  },
  resolve: {
    alias: {
      crypto: 'crypto-browserify',
      stream: 'stream-browserify',
      buffer: 'buffer',
      process: 'process/browser',
    },
  },
  optimizeDeps: {
    include: [
      'crypto-browserify',
      'stream-browserify', 
      'buffer',
      'process',
    ],
    force: true,
  },
  build: {
    rollupOptions: {
      output: {
        globals: {},
      },
    },
  },
  plugins: [
    nodePolyfills({
      protocolImports: true, // enables `node:crypto` style too
    }),
  ],
});
