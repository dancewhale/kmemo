<script setup lang="ts">
import type { ReviewStats } from '../types'

defineProps<{
  stats: ReviewStats
}>()
</script>

<template>
  <div class="review-stats">
    <div class="review-stats__item">
      <span class="review-stats__label">Completed today</span>
      <strong class="review-stats__value">{{ stats.completedToday }}</strong>
    </div>
    <div class="review-stats__item">
      <span class="review-stats__label">Remaining</span>
      <strong class="review-stats__value">{{ stats.remaining }}</strong>
    </div>
    <div class="review-stats__item">
      <span class="review-stats__label">Again count</span>
      <strong class="review-stats__value">{{ stats.againCount }}</strong>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.review-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: $space-sm;
  padding: $space-sm;
  border-bottom: 1px solid $color-border-subtle;
}

.review-stats__item {
  padding: $space-sm;
  border: 1px solid $color-border-subtle;
  border-radius: $radius-sm;
  background: $color-bg-pane-elevated;
}

.review-stats__label {
  display: block;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.review-stats__value {
  display: block;
  margin-top: 2px;
  font-size: $font-size-lg;
  line-height: 1.1;
}
</style>
