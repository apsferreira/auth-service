/** @type {import('tailwindcss').Config} */
export default {
  presets: [
    // IIT Design System base config (tokens, fonts, colors)
    // We inline the preset here since tailwind.config.ts from @iit/ui is not a CJS module
    {
      darkMode: 'class',
      theme: {
        extend: {
          colors: {
            iit: {
              brand: {
                primary: '#0097D6',
                'primary-light': '#33ADE0',
                'primary-dark': '#006FA3',
                secondary: '#00D6A0',
                'secondary-light': '#33DEB3',
                'secondary-dark': '#00A87E',
              },
              surface: {
                base: '#FFFFFF',
                subtle: '#F5F8FA',
                elevated: '#EBF4FB',
              },
              'surface-dark': {
                base: '#0A0F14',
                subtle: '#111820',
                elevated: '#1A2530',
              },
              text: {
                primary: '#0D1B26',
                secondary: '#4A6070',
                muted: '#8A9FAF',
                'on-brand': '#FFFFFF',
              },
              border: {
                default: '#D0DDE6',
                subtle: '#EBF2F7',
                brand: '#0097D6',
              },
              status: {
                success: '#00D6A0',
                'success-bg': '#E6FAF5',
                warning: '#F59E0B',
                'warning-bg': '#FEF3C7',
                error: '#EF4444',
                'error-bg': '#FEE2E2',
                info: '#0097D6',
                'info-bg': '#EBF4FB',
              },
            },
          },
          fontFamily: {
            sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
            mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
          },
          borderRadius: {
            iit: '0.5rem',
            'iit-sm': '0.375rem',
            'iit-lg': '0.75rem',
          },
        },
      },
    },
  ],
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@iit/ui/src/**/*.{ts,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eef9fd',
          100: '#d8f0f9',
          200: '#b1e2f4',
          300: '#7acde9',
          400: '#40b4da',
          500: '#0097D6',
          600: '#0097D6',
          700: '#006FA3',
          800: '#005A8A',
          900: '#004066',
          950: '#002a45',
        },
        secondary: {
          50: '#f8fafc',
          100: '#f1f5f9',
          200: '#e2e8f0',
          300: '#cbd5e1',
          400: '#94a3b8',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          800: '#1e293b',
          900: '#0f172a',
          950: '#020617',
        },
        success: {
          50: '#E6FAF5',
          100: '#ccf5eb',
          500: '#00D6A0',
          600: '#00A87E',
          700: '#008A68',
        },
        danger: {
          50: '#FEE2E2',
          100: '#fee2e2',
          400: '#f87171',
          500: '#EF4444',
          600: '#dc2626',
          700: '#b91c1c',
        },
        warning: {
          50: '#FEF3C7',
          100: '#fef3c7',
          500: '#F59E0B',
          600: '#d97706',
        },
      },
    },
  },
  plugins: [],
}
