import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

const api = process.env.CAERN_API ?? 'http://localhost:7676'

export default defineConfig({
  plugins: [svelte()],
  server: {
    host: true,
    proxy: {
      '/api': api,
      '/icons': api,
      '/images': api,
    },
  },
})
