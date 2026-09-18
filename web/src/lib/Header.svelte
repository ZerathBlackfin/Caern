<script lang="ts">
  import { edit, save, start, stop } from './edit.svelte'
  import Logo from './Logo.svelte'

  let { title }: { title: string } = $props()
</script>

<header>
  <div class="bar container">
    <Logo />
    <h1>{title}</h1>
    <nav>
      {#if edit.on}
        <button onclick={stop}>Cancel</button>
        <button class="primary" onclick={save}>Save</button>
      {:else}
        <button onclick={start}>Edit</button>
      {/if}
    </nav>
  </div>
</header>

<style lang="scss">
  @use '../styles/mixins' as *;

  header {
    position: sticky;
    top: 0;
    z-index: 10;
    border-bottom: 1px solid var(--header-border);
    background: var(--header-bg);

    :global(.has-wallpaper) & {
      backdrop-filter: var(--glass);
    }
  }

  .bar {
    display: flex;
    align-items: center;
    gap: 10px;
    height: var(--header-height);
  }

  h1 {
    @include ellipsis;

    font-size: 15px;
    font-weight: 600;
    line-height: 1.3;
  }

  nav {
    display: flex;
    gap: 8px;
    margin-left: auto;
  }

  button {
    padding: 5px 12px;
    border: 1px solid var(--border);
    border-radius: 999px;
    background: var(--surface);
    font: inherit;
    font-size: 13px;
    color: var(--text-muted);
    cursor: pointer;

    &:hover {
      background: var(--surface-hover);
      color: var(--text);
    }

    &.primary {
      border-color: transparent;
      background: var(--accent);
      font-weight: 600;
      color: #1b1403;
    }
  }
</style>
