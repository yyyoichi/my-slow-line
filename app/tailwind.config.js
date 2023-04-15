/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'my-black': '#1E2028',
        'my-light-black': '#3B4351',
        'my-white': '#AFB3C1',
        'my-red': '#F73A1C',
        'my-green': '#DFEDB6',
        'my-yellow': '#E4B45E',
      },
      keyframes: {
        wiggle: {
          '0%, 100%': { transform: 'rotate(-3deg)' },
          '50%': { transform: 'rotate(3deg)' },
        },
        fadein: {
          '0%': {
            opacity: 0,
            transform: 'translateY(5px)',
            height: 0,
          },
          '100%': {
            opacity: 1,
            transform: 'translateY(0)',
            height: 'inherit',
          },
        },
        fadeout: {
          '0%': {
            opacity: 1,
            transform: 'translateY(-2px)',
            height: 'inherit',
          },
          '100%': {
            opacity: 0,
            transform: 'translateY(0)',
            height: 0,
          },
        },
      },
      animation: {
        wiggle: 'wiggle 1s ease-in-out infinite',
        fadein: 'fadein .8s forwards',
        fadeout: 'fadeout .5s forwards',
      },
    },
  },
  plugins: [],
};
