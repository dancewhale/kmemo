<script setup lang="ts">
import { computed } from 'vue'
import { useToast } from '@/shared/composables/useToast'
import { useCardStore } from '../stores/card.store'

const props = defineProps<{
  hasReview: boolean
}>()

const card = useCardStore()
const toast = useToast()

const busy = computed(() => card.saving)

async function onAddToReview() {
  await card.addOrNavigateReviewFromSelected({ openReview: false })
}

async function onOpenReview() {
  const ok = await card.openLinkedReviewFromSelected()
  if (!ok) {
    toast.warning('No review item linked')
  }
}
</script>

<template>
  <section class="card-review-section">
    <h3 class="card-review-section__title">Review</h3>
    <p class="card-review-section__hint">
      {{ props.hasReview ? 'This card is linked to a review queue item.' : 'Not yet in the review queue.' }}
    </p>
    <div class="card-review-section__actions">
      <button type="button" class="card-review-section__btn" :disabled="busy || props.hasReview" @click="onAddToReview">
        Add to Review
      </button>
      <button type="button" class="card-review-section__btn" :disabled="!props.hasReview || busy" @click="onOpenReview">
        Open Review Item
      </button>
    </div>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.card-review-section {
  margin-top: $space-lg;
  padding-top: $space-md;
  border-top: 1px solid $color-border-subtle;
}

.card-review-section__title {
  margin: 0 0 $space-xs;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.card-review-section__hint {
  margin: 0 0 $space-sm;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.card-review-section__actions {
  display: flex;
  flex-wrap: wrap;
  gap: $space-xs;
}

.card-review-section__btn {
  font-size: $font-size-xs;
  padding: 4px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text;
  background: $color-bg-pane-elevated;
}

.card-review-section__btn:disabled {
  opacity: 0.45;
}

.card-review-section__btn:hover:not(:disabled) {
  border-color: $color-active-border;
}
</style>
