<script lang="ts">
  import { onMount } from 'svelte'
  import Block from './lib/Block.svelte'
  import EmptyState from './lib/EmptyState.svelte'
  import Header from './lib/Header.svelte'
  import { edit } from './lib/edit.svelte'
  import { connect, live } from './lib/live.svelte'
  import Problems from './lib/Problems.svelte'
  import Wallpaper from './lib/Wallpaper.svelte'

  const MIN_TILE_WIDTH = 150
  const COLUMN_WIDTH = 220

  onMount(connect)

  const page = $derived(live.state?.page)
  const problems = $derived(live.state?.problems ?? [])
  const blocks = $derived(edit.on ? edit.blocks : (page?.blocks ?? []))

  let width = $state(0)
  const fit = $derived(width ? Math.max(1, Math.floor(width / MIN_TILE_WIDTH)) : Infinity)
  const columns = $derived(Math.min(page?.columns ?? 1, fit))

  $effect(() => {
    if (!page) return
    const root = document.documentElement
    document.title = page.title
    root.dataset.theme = page.theme
    root.style.setProperty('--content-width', `${page.columns * COLUMN_WIDTH}px`)
    root.classList.toggle('has-wallpaper', !!page.backgroundImage)
    if (page.backgroundColor) root.style.setProperty('--page-bg', page.backgroundColor)
    else root.style.removeProperty('--page-bg')
  })
</script>

{#if page}
  {#if page.backgroundImage}
    <Wallpaper image={page.backgroundImage} />
  {/if}
  <Header title={page.title} />
  <main class="container">
    {#if edit.error}
      <p class="save-error" role="alert">{edit.error}</p>
    {/if}
    {#if problems.length}
      <Problems {problems} />
    {:else if !blocks.length}
      <EmptyState />
    {/if}
    <div class="blocks" bind:clientWidth={width}>
      {#each blocks as block, index (block.id ?? index)}
        <Block {block} {index} {columns} />
      {/each}
    </div>
  </main>
  {#if !live.connected}
    <p class="reconnecting" role="status">Offline</p>
  {/if}
{:else if !live.connected}
  <main class="container">
    <p class="offline">Can't reach Caern.</p>
  </main>
{/if}

<style lang="scss">
  main {
    padding-block: 28px 48px;
  }

  .blocks {
    display: flex;
    flex-direction: column;
    gap: 28px;
  }

  .offline {
    color: var(--text-muted);
  }

  .save-error {
    margin-bottom: 24px;
    padding: 10px 14px;
    border-left: 2px solid var(--error);
    border-radius: var(--radius-panel);
    background: var(--error-bg);
  }

  .reconnecting {
    position: fixed;
    right: 16px;
    bottom: 16px;
    margin: 0;
    padding: 8px 14px;
    border: 1px solid var(--border);
    border-radius: 999px;
    background: var(--surface-glass);
    backdrop-filter: var(--glass);
    font-size: 12px;
    color: var(--text-muted);
  }
</style>
