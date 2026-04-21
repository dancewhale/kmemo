<script setup lang="ts">
import type { ReviewItem } from '../types'

defineProps<{
  item: ReviewItem
}>()
</script>

<template>
  <section class="review-meta">
    <h3 class="review-meta__title">Metadata</h3>
    <dl class="review-meta__list">
      <div><dt>Item id</dt><dd>{{ item.id }}</dd></div>
      <div><dt>Due</dt><dd>{{ item.dueAt }}</dd></div>
      <div><dt>Status</dt><dd>{{ item.status }}</dd></div>
      <div><dt>Review count</dt><dd>{{ item.reviewCount }}</dd></div>
      <div><dt>Last reviewed</dt><dd>{{ item.lastReviewedAt ?? '—' }}</dd></div>
    </dl>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.review-meta__title {
  margin: 0 0 $space-sm;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.review-meta__list {
  margin: 0;
  font-size: $font-size-xs;
}

.review-meta__list > div {
  display: grid;
  grid-template-columns: 96px 1fr;
  gap: $space-sm;
  padding: $space-xs 0;
  border-top: 1px solid $color-border-subtle;
}

.review-meta__list dt {
  margin: 0;
  color: $color-text-secondary;
}

.review-meta__list dd {
  margin: 0;
  color: $color-text;
  font-family: $font-mono;
  word-break: break-all;
}
</style>
