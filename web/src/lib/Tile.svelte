<script lang="ts">
  import TileIcon from './TileIcon.svelte'
  import type { Tile } from './types'

  let { tile, columns }: { tile: Tile; columns: number } = $props()

  const tall = $derived(tile.height >= 2)
</script>

<a
  class="tile"
  class:tall
  href={tile.url}
  target={tile.newTab ? '_blank' : undefined}
  rel={tile.newTab ? 'noopener noreferrer' : undefined}
  style:grid-column="span {Math.min(tile.width, columns)}"
  style:grid-row="span {tile.height}"
>
  <TileIcon src={tile.icon} name={tile.name} large={tall} />
  <span class="text">
    <span class="name">{tile.name}</span>
    {#if tile.description}
      <span class="description">{tile.description}</span>
    {/if}
  </span>
</a>

<style lang="scss">
  @use '../styles/mixins' as *;

  .tile {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
    padding: 0 12px;
    border: 1px solid var(--border);
    border-radius: var(--radius-tile);
    background: var(--surface);
    transition:
      background-color 160ms ease,
      border-color 160ms ease;

    &:hover {
      background: var(--surface-hover);
      border-color: var(--border-hover);
    }

    &.tall {
      flex-direction: column;
      align-items: flex-start;
      justify-content: space-between;
      padding: 12px;
    }

    :global(.has-wallpaper) & {
      background: var(--surface-glass);
      backdrop-filter: var(--glass);
    }
  }

  .text {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .name,
  .description {
    @include ellipsis;
  }

  .name {
    font-size: 14px;
    font-weight: 500;
    line-height: 1.3;
  }

  .description {
    font-size: 12px;
    line-height: 1.35;
    color: var(--text-muted);
  }

  .tall .name {
    font-size: 15px;
  }
</style>
