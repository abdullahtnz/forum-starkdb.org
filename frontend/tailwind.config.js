/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        stark: {
          bg: '#0a0a0f',
          surface: '#111118',
          'surface-2': '#18181f',
          border: '#2a2a38',
          accent: '#e8ff47',
          'accent-2': '#47ffe8',
          'accent-3': '#ff6b47',
          'accent-4': '#c47fff',
          'accent-5': '#ff9f47',
          text: '#e8e8f0',
          muted: '#7070a0',
        },
      },
      fontFamily: {
        display: ['Syne', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['Space Mono', 'monospace'],
      },
      animation: {
        'fade-up': 'fadeUp 0.6s ease-out both',
        'fade-down': 'fadeDown 0.6s ease-out both',
        pulse: 'pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        fadeUp: {
          from: { opacity: '0', transform: 'translateY(20px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        fadeDown: {
          from: { opacity: '0', transform: 'translateY(-20px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
