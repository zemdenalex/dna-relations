/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        dna: {
          primary: '#6366f1',
          secondary: '#ec4899',
          bg: '#0f0f0f',
          surface: '#1a1a1a',
          border: '#2a2a2a',
        },
      },
    },
  },
  plugins: [],
}
