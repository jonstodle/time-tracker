<script lang="ts">
  import dayjs from "dayjs";
  import { slide } from "svelte/transition";
  import ListItem from "./ListItem.svelte";
  import type { Tracker } from "./types";

  const storageKey = "trackers";

  let tally = 0;
  let trackers: Tracker[] = [];

  $: trackersView = trackers.sort((a, b) => b.start.localeCompare(a.start));

  const stored = localStorage.getItem(storageKey);
  if (stored) {
    const parsed = JSON.parse(stored);
    tally = parsed.tally || 0;
    trackers =
      (parsed.trackers &&
        parsed.trackers.map((t) => {
          if (!t.id) {
            t.id = generateId();
          }
          return t;
        })) ||
      [];
  }

  function addTracker() {
    trackers = [
      ...trackers,
      {
        id: generateId(),
        description: `Tracker #${++tally}`,
        start: dayjs().toISOString(),
      },
    ];
  }

  function removeTracker({ detail: tracker }: { detail: Tracker }) {
    trackers = trackers.filter((t) => t != tracker);
  }

  function generateId() {
    return String.fromCharCode(
      ...[...crypto.getRandomValues(new Uint8Array(5))].map(
        (v) => (v % 26) + 97
      )
    );
  }

  setInterval(() => {
    const now = dayjs();
    trackers = trackers.map((tracker) => {
      tracker.elapsed = dayjs.duration(
        ((tracker.stop && dayjs(tracker.stop)) || now).diff(tracker.start)
      );
      return tracker;
    });
  }, 1000);

  $: {
    localStorage.setItem(
      storageKey,
      JSON.stringify({
        trackers: trackers.map((t) => ({
          id: t.id,
          start: t.start,
          ...(t.stop && { stop: t.stop }),
          description: t.description,
        })),
        tally,
      })
    );
  }
</script>

<nav>
  <button class="button is-link" on:click={addTracker}>Add tracker</button>
</nav>

<main>
  <div>
    {#each trackersView as tracker (tracker.id)}
      <div class="stack" in:slide out:slide={{ duration: 200 }}>
        <ListItem {tracker} on:removeTracker={removeTracker} />
      </div>
    {/each}
  </div>
</main>

<style>
  main {
    padding: 3em 1.5em;
  }

  main > div {
    margin: 0 auto;
    max-width: 1200px;
  }

  nav {
    display: flex;
    justify-content: flex-end;
    position: relative;
  }

  .stack {
    border: lightgray solid 1px;
  }

  .stack:not(:first-child) {
    border-top: 0;
  }
</style>
