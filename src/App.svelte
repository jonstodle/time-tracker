<script>
    import {onMount} from 'svelte'
    import {blur} from 'svelte/transition'
    import dayjs from 'dayjs'
    import duration from 'dayjs/plugin/duration'

    dayjs.extend(duration)

    const storageKey = 'trackers'

    let tally = 0
    let trackers = []

    const stored = localStorage.getItem(storageKey)
    if (stored) {
        const parsed = JSON.parse(stored)
        tally = parsed.tally || 0
        trackers = parsed.trackers && parsed.trackers
                .map(tracker => {
                    tracker.start = dayjs(tracker.start)
                    if (tracker.stop) {
                        tracker.stop = dayjs(tracker.stop)
                    }
                    return tracker
                }) || []
    }

    function addTracker() {
        trackers = [{
            description: `Tracker #${++tally}`,
            start: dayjs(),
        }, ...trackers]
    }

    function removeTracker(tracker) {
        trackers = trackers.filter(t => t != tracker)
    }

    function formatDuration(duration) {
        return `${duration.hours().toString().padStart(2, "0")}:${duration.minutes().toString().padStart(2, "0")}:${duration.seconds().toString().padStart(2, "0")}`
    }

    setInterval(() => {
        const now = dayjs()
        trackers = trackers
                .map(tracker => {
                    tracker.elapsed = dayjs.duration((tracker.stop || now).diff(tracker.start))
                    return tracker
                })
    }, 1000)

    $: {
        localStorage.setItem(storageKey, JSON.stringify({
            trackers: trackers.map(t => ({
                start: t.start.toISOString(),
                ...(t.stop && {stop: t.stop.toISOString()}),
                description: t.description,
            })),
            tally,
        }))
    }
</script>

<nav>
    <button class="button is-link" on:click={addTracker}>Add tracker</button>
</nav>

<main>
    <div>
        {#each trackers as tracker}
            <div class="stack list-item" class:hover-show={!!tracker.stop} in:blur={{duration: 200}} out:blur>
                <div>
                    {#if !tracker.editable}
                        <span on:click={() => tracker.editable = true}>{tracker.description}</span>
                    {:else}
                        <input type="text" style="margin-left: calc(-0.6em)"
                               bind:value={tracker.description}
                               on:keyup={(e) => tracker.editable = e.key != 'Enter'}
                               autofocus>
                        <button class="success" on:click={() => tracker.editable = false}>
                            Done
                        </button>
                    {/if}
                </div>
                <div>
                    <span>{formatDuration(tracker.elapsed || dayjs.duration(0))}</span>
                    {#if !tracker.stop}
                        <button class="button warning" disabled={!!tracker.stop}
                                on:click={() => tracker.stop = dayjs()}>
                            Stop
                        </button>
                    {:else}
                        <button class="button error" on:click={() => removeTracker(tracker)}>
                            Remove
                        </button>
                    {/if}
                </div>
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

    .list-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: .75rem;
        border: lightgray solid 1px;
    }

    .list-item:not(:first-child) {
        border-top: 0;
    }

    .list-item > div {
        display: flex;
        align-items: center;
    }

    .list-item > div > *:not(:last-child) {
        margin-right: .25rem;
    }

    .list-item > div > span:first-child {
        border: white solid 1px;
    }

    .list-item button {
        width: 80px;
        font-size: .75em;
    }
</style>
