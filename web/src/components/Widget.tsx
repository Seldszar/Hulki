import { createMemo, For, Show } from "solid-js";

import { State } from "~/types";

interface WidgetProps {
  state: State;
}

function Widget(props: WidgetProps) {
  const {
    state: { achievements },
  } = props;

  const unlockedAchievements = createMemo(() => achievements.filter((achievement) => achievement.achieved).sort((a, b) => b.unlockedAt - a.unlockedAt));

  const unlockedCount = createMemo(() => unlockedAchievements().length);
  const totalCount = createMemo(() => achievements.length);

  return (
    <Show when={achievements.length > 0}>
      <div class="bg-zinc-700 overflow-hidden shadow-md shadow-black/50">
        <div class="bg-zinc-600 px-8 py-6">
          <div class="mb-2 text-lg">
            Unlocked Achievements: {unlockedCount()} / {totalCount()} <span class="text-white/50">({((unlockedCount() / totalCount()) * 100).toLocaleString(undefined, { maximumFractionDigits: 0 })}%)</span>
          </div>
          <div class="flex rounded bg-black h-3">
            <div class="rounded bg-blue-500" style={{ width: `${(unlockedCount() / totalCount()) * 100}%` }} />
          </div>
        </div>

        <div class="px-2">
          <div class="flex overflow-hidden [mask-image:linear-gradient(to_left,transparent,black_8rem)]">
            <For each={unlockedAchievements().slice(0, 10)}>
              {(achievement) => (
                <div class="flex flex-none items-center gap-3 border-r border-zinc-800 p-6 last:border-none">
                  <img src={achievement.icon} class="flex-none w-12" />

                  <div class="flex-1">
                    <div>{achievement.displayName}</div>
                    <div class="text-sm text-white/50">{new Date(achievement.unlockedAt * 1000).toLocaleString()}</div>
                  </div>
                </div>
              )}
            </For>
          </div>
        </div>
      </div>
    </Show>
  );
}

export default Widget;
