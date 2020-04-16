<script>
    import {onMount} from 'svelte'
    import {blur} from 'svelte/transition'
    import {DateTime, Duration} from 'luxon'

    const storageKey = 'trackers'

    let tally = 0
    let trackers = []

    const stored = localStorage.getItem(storageKey)
    if (stored) {
        const parsed = JSON.parse(stored)
        tally = parsed.tally || 0
        trackers = parsed.trackers && parsed.trackers
                .map(tracker => {
                    tracker.start = DateTime.fromISO(tracker.start)
                    return tracker
                }) || []
    }

    function addTracker() {
        trackers = [...trackers, {
            description: `Tracker #${++tally}`,
            start: DateTime.local(),
        }]
    }

    function removeTracker(tracker) {
      trackers = trackers.filter(t => t != tracker)
    }

    setInterval(() => {
        const now = DateTime.local()
        trackers = trackers
                .map(tracker => {
                    tracker.elapsed = (tracker.stop || now).diff(tracker.start)
                    return tracker
                })
    }, 1000)

    $: {
        localStorage.setItem(storageKey, JSON.stringify({
            trackers: trackers.map(t => ({
                start: t.start,
                stop: t.stop,
                description: t.description,
            })),
            tally,
        }))
    }
</script>

<main class="section">
    <article class="container">
        <nav class="level">
            <div class="level-left"></div>
            <div class="level-right">
                <p class="level-item">
                    <button class="button is-link" on:click={addTracker}>Add tracker</button>
                </p>
            </div>
        </nav>

        {#each trackers as tracker}
            <div class="level box" class:hover-show={!!tracker.stop} in:blur={{duration: 200}} out:blur>
                <div class="level-left">
                    {#if !tracker.editable}
                        <span class="level-item" on:click={() => tracker.editable = true}>{tracker.description}</span>
                    {:else}
                        <input type="text" style="margin-left: calc(-0.75em)" class="input level-item"
                               bind:value={tracker.description} on:keyup={(e) => tracker.editable = e.key != 'Enter'} autofocus>
                        <button class="button is-success level-item" on:click={() => tracker.editable = false}>
                            Done
                        </button>
                    {/if}
                </div>
                <div class="level-right">
                    <span class="level-item">{(tracker.elapsed || Duration.fromMillis(0)).toFormat('hh:mm:ss')}</span>
                    <button class="button is-small is-danger level-item is-hidden" on:click={() => removeTracker(tracker)}>
                        Remove
                    </button>
                    <button class="button is-small is-danger level-item" disabled={!!tracker.stop}
                            on:click={() => tracker.stop = DateTime.local()}>
                        Stop
                    </button>
                </div>
            </div>
        {/each}
    </article>
</main>

<style>
    .hover-show:hover button.is-hidden {
        display: initial !important;
    }
</style>
