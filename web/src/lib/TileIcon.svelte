<script lang="ts">
  let { src, name, large = false }: { src?: string; name: string; large?: boolean } = $props()

  let failedSrc = $state<string>()
  const failed = $derived(failedSrc === src)
</script>

<span class="icon" class:large>
  {#if src && !failed}
    <img {src} alt="" decoding="async" onerror={() => (failedSrc = src)} />
  {:else}
    <span class="initial">{name.trim().charAt(0).toUpperCase()}</span>
  {/if}
</span>

<style lang="scss">
  .icon {
    --size: 28px;

    display: block;
    flex: none;
    width: var(--size);
    height: var(--size);

    &.large {
      --size: 40px;
    }
  }

  img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .initial {
    display: grid;
    place-items: center;
    width: 100%;
    height: 100%;
    border-radius: var(--radius-icon);
    background: var(--icon-backdrop);
    font-size: calc(var(--size) * 0.45);
    font-weight: 600;
    color: var(--accent);
  }
</style>
