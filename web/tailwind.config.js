/* eslint-env node */

const defaultTheme = require("tailwindcss/defaultTheme");

/**
 * @type {import("tailwindcss").Config}
 */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Noto Sans Display", ...defaultTheme.fontFamily.sans],
      },
    },
  },
};
