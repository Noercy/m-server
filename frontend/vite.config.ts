import { defineConfig } from 'vite'
import preact from '@preact/preset-vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [preact()],
  resolve: {
    alias: {
      react: 'preact/compat',
      'react-dom': "preact/compat"
    }
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true, 
        secure: false
      },
      "/thumbnails" : "http://localhost:8080",
      "/pages" : "http://localhost:8080", 
    }
  }
})
