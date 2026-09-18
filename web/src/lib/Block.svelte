<script lang="ts">
  import { edit } from './edit.svelte'
  import Tile from './Tile.svelte'
  import type { Block } from './types'

  let { block, index, columns }: { block: Block; index: number; columns: number } = $props()
</script>

<section>
  {#if block.section}
    <h2>{block.section}</h2>
  {/if}
  <div class="grid" class:editing={edit.on} data-block={index} style:--columns={columns}>
    {#each block.tiles as tile (tile.id)}
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

    &.editing {
      min-height: var(--row);
      border-radius: var(--radius-tile);
      outline: 1px dashed var(--border);
      outline-offset: 6px;
    }
  }
</style>
