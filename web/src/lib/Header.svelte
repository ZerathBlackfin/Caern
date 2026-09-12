<script lang="ts">
  import Logo from './Logo.svelte'

  let { title }: { title: string } = $props()
</script>

<header>
  <div class="bar container">
    <Logo />
    <h1>{title}</h1>
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
</style>
