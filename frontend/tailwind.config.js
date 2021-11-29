module.exports = {
  purge: ['./src/**/*.svelte'],
  darkMode: false,
  theme: {
    extend: {
      colors: {
        primary: '#4E7CFF',
        secondary: '#202442',
      }
    },
    fontFamily: {
      'display': ['Catamaran', 'sans-serif'],
      'body': ['Poppins', 'sans-serif'],
    }
  },
  variants: {
    extend: {
      ringColor: ['hover']
    },
  },
  plugins: [],
}
