import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  base: '/web/',
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:6116',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: '../dist',
    emptyOutDir: true
  }
})
