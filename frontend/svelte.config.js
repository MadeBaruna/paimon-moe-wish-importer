const sveltePreprocess = require('svelte-preprocess');
const tailwindcss = require('tailwindcss');
const autoprefixer = require('autoprefixer');
const tailwindConfig = require('./tailwind.config');

module.exports = {
  preprocess: sveltePreprocess({
    postcss: {
      plugins: [
        tailwindcss(tailwindConfig),
        autoprefixer,
      ]
    }
  })
};
