<script setup lang="ts">
import type { CardItem } from '../types'

const props = defineProps<{
  card: CardItem
  hasReview: boolean
}>()
</script>

<template>
  <section class="card-meta">
    <h3 class="card-meta__title">Meta</h3>
    <dl class="card-meta__dl">
      <div><dt>Card id</dt><dd class="kmono">{{ props.card.id }}</dd></div>
      <div><dt>Node id</dt><dd class="kmono">{{ props.card.nodeId }}</dd></div>
      <div><dt>Parent</dt><dd class="kmono">{{ props.card.parentNodeId ?? '—' }}</dd></div>
      <div><dt>Created</dt><dd>{{ props.card.createdAt }}</dd></div>
      <div><dt>Updated</dt><dd>{{ props.card.updatedAt }}</dd></div>
      <div><dt>Review</dt><dd>{{ props.hasReview ? 'Linked' : 'Not linked' }}</dd></div>
    </dl>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.card-meta__title {
  margin: 0 0 $space-xs;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.card-meta__dl {
  margin: 0;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.card-meta__dl > div {
  display: grid;
  grid-template-columns: 72px 1fr;
  gap: $space-sm;
  padding: $space-xs 0;
  border-top: 1px solid $color-border-subtle;
}

.card-meta__dl dt {
  margin: 0;
}

.card-meta__dl dd {
  margin: 0;
  color: $color-text;
  word-break: break-all;
}
</style>
