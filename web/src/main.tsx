import { render } from "solid-js/web";

import App from "./components/App";

const container = document.getElementById("app-root");

if (container) {
  render(() => <App />, container);
}
