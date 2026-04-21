<script setup lang="ts">
import type { ReviewGrade } from '../types'

const props = defineProps<{
  enabled: boolean
  submitting: boolean
}>()

const emit = defineEmits<{
  grade: [grade: ReviewGrade]
}>()

const grades: Array<{ key: ReviewGrade; label: string }> = [
  { key: 'again', label: 'Again' },
  { key: 'hard', label: 'Hard' },
  { key: 'good', label: 'Good' },
  { key: 'easy', label: 'Easy' },
]

function onGrade(g: ReviewGrade) {
  if (!props.enabled || props.submitting) {
    return
  }
  emit('grade', g)
}
</script>

<template>
  <div class="review-actions">
    <button
      v-for="g in grades"
      :key="g.key"
      type="button"
      class="review-actions__btn"
      :class="`review-actions__btn--${g.key}`"
      :disabled="!enabled || submitting"
      @click="onGrade(g.key)"
    >
      {{ submitting ? 'Submitting…' : g.label }}
    </button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.review-actions {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: $space-sm;
}

.review-actions__btn {
  font-size: $font-size-xs;
  padding: 6px 0;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
}

.review-actions__btn--again {
  border-color: rgba(195, 120, 120, 0.5);
}

.review-actions__btn--hard {
  border-color: rgba(196, 165, 104, 0.5);
}

.review-actions__btn--good {
  border-color: rgba(107, 155, 209, 0.5);
}

.review-actions__btn--easy {
  border-color: rgba(122, 180, 132, 0.5);
}

.review-actions__btn:disabled {
  opacity: 0.45;
}
</style>
