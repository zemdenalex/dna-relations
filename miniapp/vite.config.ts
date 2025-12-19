import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  base: '/miniapp/',
  server: {
    port: 5174,
    host: true,
  },
  build: {
    outDir: 'dist',
  },
})
