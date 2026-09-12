<script lang="ts">
  let { image }: { image: string } = $props()
</script>

<div class="wallpaper" style:background-image={`url("${encodeURI(image)}")`} aria-hidden="true"></div>

<style lang="scss">
  .wallpaper {
    position: fixed;
    inset: 0;
    z-index: -1;
    background-size: cover;
    background-position: center;

    &::after {
      content: '';
      position: absolute;
      inset: 0;
      background: var(--wallpaper-scrim);
    }
  }
</style>
