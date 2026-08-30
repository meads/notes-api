import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react' // or vue, svelte, etc.

export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // Needed for Docker tracking
    port: 5173, 
  }
})