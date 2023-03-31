import { createResource, For, Show } from "solid-js";

import { State } from "~/types";

import Widget from "./Widget";

function App() {
  const [state, { refetch }] = createResource(async () => {
    const response = await fetch("/api/state");
    const json = await response.json();

    return json as State;
  });

  setInterval(refetch, 10_000);

  return (
    <div class="p-4">
      <Show keyed when={state()}>
        {(state) => <Widget state={state} />}
      </Show>
    </div>
  );
}

export default App;
