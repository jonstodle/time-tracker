<script lang="ts">
  export let icon: string;
  console.log($$slots);
</script>

<button {...$$restProps} type="button" on:click>
  <span class="icon icon-{icon}" class:no-spacing={!$$slots.default} /><slot />
</button>

<style>
  button {
    display: flex;
    align-items: center;
    transition: width 100ms;
  }

  .icon {
    opacity: 0.9;
  }

  span.icon.no-spacing::before {
    margin-right: 0;
  }
</style>
