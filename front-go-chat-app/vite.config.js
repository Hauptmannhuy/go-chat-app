import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8090',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
      '/socket': {
        target: 'ws://localhost:8090',
        ws: true,
        rewriteWsOrigin: true,
        rewrite: (path) => path.replace(/^\/socket/, ''),
      },
    }
  }
})
