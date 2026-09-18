<script lang="ts">
  import { edit, startDrag, startResize } from './edit.svelte'
  import TileIcon from './TileIcon.svelte'
  import type { Tile } from './types'

  let { tile, columns }: { tile: Tile; columns: number } = $props()

  const tall = $derived(tile.height >= 2)
  const resize = $derived(edit.resize?.id === tile.id ? edit.resize : null)
  const drag = $derived(edit.drag?.id === tile.id ? edit.drag : null)
</script>

{#if drag}
  <div
    class="placeholder"
    style:grid-column="span {Math.min(tile.width, columns)}"
    style:grid-row="span {tile.height}"
  ></div>
{/if}

<a
  class="tile"
  class:tall
  class:editing={edit.on}
  class:dragging={!!drag}
  class:resizing={!!resize}
  data-tile={tile.id}
  href={edit.on ? undefined : tile.url}
  target={tile.newTab ? '_blank' : undefined}
  rel={tile.newTab ? 'noopener noreferrer' : undefined}
  style:grid-column="span {Math.min(tile.width, columns)}"
  style:grid-row="span {tile.height}"
  style:width={resize ? `${resize.pxWidth}px` : drag ? `${drag.width}px` : undefined}
  style:height={resize ? `${resize.pxHeight}px` : drag ? `${drag.height}px` : undefined}
  style:left={drag ? `${drag.left}px` : undefined}
  style:top={drag ? `${drag.top}px` : undefined}
>
  <TileIcon src={tile.icon} name={tile.name} large={tall} />
  <span class="text">
    <span class="name">{tile.name}</span>
    {#if tile.description}
      <span class="description">{tile.description}</span>
    {/if}
  </span>
  {#if edit.on}
    <span class="grip" aria-hidden="true" onpointerdown={(e) => startDrag(tile.id, e)}>
      <svg viewBox="0 0 12 8" fill="currentColor"><circle cx="2" cy="2" r="1.1" /><circle cx="6" cy="2" r="1.1" /><circle cx="10" cy="2" r="1.1" /><circle cx="2" cy="6" r="1.1" /><circle cx="6" cy="6" r="1.1" /><circle cx="10" cy="6" r="1.1" /></svg>
    </span>
    <span class="resize" aria-hidden="true" onpointerdown={(e) => startResize(tile.id, e)}>
      <svg viewBox="0 0 10 10" stroke="currentColor" stroke-width="1.2" stroke-linecap="round">
        <line x1="9" y1="1" x2="1" y2="9" />
        <line x1="9" y1="5.5" x2="5.5" y2="9" />
      </svg>
      {#if resize}<b>{resize.width}x{resize.height}</b>{/if}
    </span>
  {/if}
</a>

<style lang="scss">
  @use '../styles/mixins' as *;

  .tile {
    position: relative;
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

    &.editing {
      user-select: none;
    }

    &.dragging {
      position: fixed;
      z-index: 30;
      pointer-events: none;
      opacity: 0.9;
    }

    &.resizing {
      z-index: 5;
    }

    :global(.has-wallpaper) & {
      background: var(--surface-glass);
      backdrop-filter: var(--glass);
    }
  }

  .placeholder {
    border: 1px dashed var(--border-hover);
    border-radius: var(--radius-tile);
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

  .grip {
    position: absolute;
    top: 0;
    left: 50%;
    display: grid;
    place-items: center;
    width: 30px;
    height: 14px;
    translate: -50%;
    color: var(--text-muted);
    opacity: 0.4;
    cursor: grab;
    touch-action: none;

    &:hover {
      opacity: 1;
    }

    svg {
      width: 11px;
      height: 7px;
    }
  }

  .resize {
    position: absolute;
    right: 0;
    bottom: 0;
    display: grid;
    place-items: center;
    width: 18px;
    height: 18px;
    color: var(--text-muted);
    opacity: 0.4;
    cursor: nwse-resize;
    touch-action: none;

    &:hover {
      opacity: 1;
    }

    svg {
      width: 10px;
      height: 10px;
    }
  }

  b {
    position: absolute;
    right: 18px;
    bottom: 4px;
    font-family: var(--mono);
    font-size: 11px;
    font-weight: 400;
    color: var(--text-muted);
  }
</style>
