<script lang="ts">
  import Tile from './Tile.svelte'
  import type { Block } from './types'

  let { block, columns }: { block: Block; columns: number } = $props()
</script>

<section>
  {#if block.section}
    <h2>{block.section}</h2>
  {/if}
  <div class="grid" style:--columns={columns}>
    {#each block.tiles as tile}
      <Tile {tile} {columns} />
    {/each}
  </div>
</section>

<style lang="scss">
  h2 {
    margin-bottom: 10px;
    font-size: 13px;
    font-weight: 600;
    line-height: 1.3;
    color: var(--text-muted);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(var(--columns), minmax(0, 1fr));
    grid-auto-rows: var(--row);
    grid-auto-flow: dense;
    gap: var(--gap);
  }
</style>
