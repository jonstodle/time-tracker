<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { Tracker } from "./types";
  import dayjs from "dayjs";

  export let tracker: Tracker;

  const removeTracker = createEventDispatcher().bind(null, "removeTracker");

  let isEditing = false;
  let inputElement: HTMLInputElement;

  function startEditing() {
    isEditing = true;
    inputElement.select();
  }

  function stopEditing() {
    inputElement.blur();
    isEditing = false;
  }

  function formatDuration(duration) {
    return `${duration.hours().toString().padStart(2, "0")}:${duration
      .minutes()
      .toString()
      .padStart(2, "0")}:${duration.seconds().toString().padStart(2, "0")}`;
  }
</script>

<article>
  <div>
    <input
      type="text"
      bind:this={inputElement}
      class:is-editing={isEditing}
      bind:value={tracker.description}
      on:click={startEditing}
      on:keyup={(e) => isEditing && e.key == "Enter" && stopEditing()}
      on:focusout={stopEditing}
    />
  </div>
  <div>
    <span>{formatDuration(tracker.elapsed || dayjs.duration(0))}</span>
    {#if !tracker.stop}
      <button class="button warning" on:click={() => (tracker.stop = dayjs())}>
        Stop
      </button>
    {:else}
      <button class="button error" on:click={() => removeTracker(tracker)}>
        Remove
      </button>
    {/if}
  </div>
</article>

<style>
  input {
    border-color: transparent;
    margin-right: 0.25rem;
  }

  input.is-editing {
    border-color: inherit;
  }

  article {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    border: lightgray solid 1px;
  }

  article:not(:first-child) {
    border-top: 0;
  }

  article > div {
    display: flex;
    align-items: center;
  }

  article > div > *:not(:last-child) {
    margin-right: 0.25rem;
  }

  article > div > span:first-child {
    border: white solid 1px;
  }

  article button {
    width: 80px;
    font-size: 0.75em;
  }
</style>
