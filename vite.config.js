import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  publicDir: 'static',
  resolve: {
    alias: {
      $lib: path.resolve('./src/lib')
    }
  },
  build: {
    outDir: 'public',
    emptyOutDir: true
  },
  server: {
    proxy: {
      // Proxy API endpoints to backend
      '^/(config|login|download|upload|api)': {
        target: 'http://localhost:8088',
        changeOrigin: true
      },
      '/ws': {
        target: 'ws://localhost:8088',
        ws: true
      },
      // Proxy POST/PUT to root for file uploads (keeps CLI simple)
      '/': {
        target: 'http://localhost:8088',
        changeOrigin: true,
        bypass: (req) => {
          // Only proxy POST/PUT requests, let Vite handle GET (static files)
          if (req.method === 'POST' || req.method === 'PUT') {
            return null; // null means use the proxy
          }
          return req.url; // return URL to bypass proxy and serve static files
        }
      }
    }
  }
})
