<script lang="ts">
  import type { Problem } from './types'

  let { problems }: { problems: Problem[] } = $props()

  const title = $derived(
    problems.length === 1 ? 'Configuration error' : `${problems.length} configuration errors`,
  )
</script>

<div class="problems" role="alert">
  <p>{title}. Showing the last valid version.</p>
  <ul>
    {#each problems as problem}
      <li>
        <span class="where">{problem.file}{problem.line ? `:${problem.line}` : ''}</span>
        {problem.message}
      </li>
    {/each}
  </ul>
</div>

<style lang="scss">
  .problems {
    margin-bottom: 24px;
    padding: 10px 14px;
    border-left: 2px solid var(--error);
    border-radius: var(--radius-panel);
    background: var(--error-bg);
  }

  ul {
    margin: 8px 0 0;
    padding: 0;
    list-style: none;
  }

  li {
    padding-left: 1.6em;
    font-family: var(--mono);
    font-size: 12px;
    line-height: 1.7;
    text-indent: -1.6em;
  }

  .where {
    color: var(--text-muted);
  }
</style>
