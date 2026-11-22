import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Проксируем запросы к API бэкенда (Go, порт 8080) в dev-окружении.
      '/api': 'http://localhost:8080',
    },
  },
})
