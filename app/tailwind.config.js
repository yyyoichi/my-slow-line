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
    },
  },
  plugins: [],
};
