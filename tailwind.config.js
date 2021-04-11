const colors = require('tailwindcss/colors')

module.exports = {
  mode: "jit",
  purge: [
    './server/templates/*.html',
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        gray: colors.blueGray
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
