import { resolve } from "path";
import { defineConfig } from "vite";

import soliid from "vite-plugin-solid";

export default defineConfig({
  plugins: [soliid()],
  resolve: {
    alias: {
      "~": resolve(__dirname, "src"),
    },
  },
});
