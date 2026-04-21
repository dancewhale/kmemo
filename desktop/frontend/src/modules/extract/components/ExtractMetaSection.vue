<script setup lang="ts">
import type { ExtractItem } from '../types'

defineProps<{
  extract: ExtractItem
}>()
</script>

<template>
  <section class="extract-meta">
    <h3 class="extract-meta__title">Metadata</h3>
    <dl class="extract-meta__list">
      <div><dt>Extract id</dt><dd>{{ extract.id }}</dd></div>
      <div><dt>Tree node</dt><dd>{{ extract.treeNodeId ?? '—' }}</dd></div>
      <div><dt>Parent node</dt><dd>{{ extract.parentNodeId ?? '—' }}</dd></div>
      <div><dt>Source article</dt><dd>{{ extract.sourceArticleId }}</dd></div>
      <div><dt>Created</dt><dd>{{ extract.createdAt }}</dd></div>
      <div><dt>Updated</dt><dd>{{ extract.updatedAt }}</dd></div>
    </dl>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.extract-meta__title {
  margin: 0 0 $space-sm;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.extract-meta__list {
  margin: 0;
  font-size: $font-size-xs;
}

.extract-meta__list > div {
  display: grid;
  grid-template-columns: 88px 1fr;
  gap: $space-sm;
  padding: $space-xs 0;
  border-top: 1px solid $color-border-subtle;
}

.extract-meta__list dt {
  margin: 0;
  color: $color-text-secondary;
}

.extract-meta__list dd {
  margin: 0;
  color: $color-text;
  font-family: $font-mono;
  word-break: break-all;
}
</style>
