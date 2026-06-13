import plugin from 'tailwindcss/plugin'

/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        bg: {
          base: 'var(--color-bg-base)',
          surface: 'var(--color-bg-surface)',
          elevated: 'var(--color-bg-elevated)',
          overlay: 'var(--color-bg-overlay)',
        },
        text: {
          primary: 'var(--color-text-primary)',
          secondary: 'var(--color-text-secondary)',
          tertiary: 'var(--color-text-tertiary)',
          muted: 'var(--color-text-muted)',
        },
        border: {
          DEFAULT: 'var(--color-border-default)',
          subtle: 'var(--color-border-subtle)',
        },
        brand: {
          DEFAULT: 'var(--color-brand)',
          hover: 'var(--color-brand-hover)',
          bg: 'var(--color-brand-bg)',
        },
        success: {
          DEFAULT: 'var(--color-success)',
          bg: 'var(--color-success-bg)',
          border: 'var(--color-success-border)',
        },
        danger: {
          DEFAULT: 'var(--color-danger)',
          bg: 'var(--color-danger-bg)',
          border: 'var(--color-danger-border)',
        },
      },
      boxShadow: {
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
      },
    },
  },
  plugins: [
    plugin(({ addUtilities }) => {
      addUtilities({
        '.bg-base': { backgroundColor: 'var(--color-bg-base)' },
        '.bg-surface': { backgroundColor: 'var(--color-bg-surface)' },
        '.bg-elevated': { backgroundColor: 'var(--color-bg-elevated)' },
        '.text-primary': { color: 'var(--color-text-primary)' },
        '.text-secondary': { color: 'var(--color-text-secondary)' },
        '.text-tertiary': { color: 'var(--color-text-tertiary)' },
        '.border-default': { borderColor: 'var(--color-border-default)' },
        '.border-subtle': { borderColor: 'var(--color-border-subtle)' },
      })
    }),
  ],
}
