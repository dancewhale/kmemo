<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useToast } from '@/shared/composables/useToast'
import { useReviewStore } from '@/modules/review/stores/review.store'
import { useExtractStore } from '../stores/extract.store'

const extract = useExtractStore()
const review = useReviewStore()
const toast = useToast()
const { selectedExtract } = storeToRefs(extract)

const hasReview = computed(() => {
  const ex = selectedExtract.value
  if (!ex) {
    return false
  }
  if (ex.reviewItemId) {
    return true
  }
  return Boolean(review.getReviewItemByExtractId(ex.id))
})

async function onCreate() {
  await extract.addOrOpenReviewFromSelectedExtract({ openReview: false })
}

async function onOpen() {
  const ok = await extract.openLinkedReviewFromSelected()
  if (!ok) {
    toast.warning('No review item linked')
  }
}
</script>

<template>
  <section v-if="selectedExtract" class="extract-review-section">
    <h3 class="extract-review-section__title">Review</h3>
    <p class="extract-review-section__hint">
      {{ hasReview ? 'Linked to a review queue item.' : 'Not yet in the review queue.' }}
    </p>
    <div class="extract-review-section__actions">
      <button type="button" class="extract-review-section__btn" :disabled="hasReview" @click="onCreate">
        Create Review Item
      </button>
      <button type="button" class="extract-review-section__btn" :disabled="!hasReview" @click="onOpen">
        Open Review Item
      </button>
    </div>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.extract-review-section {
  margin-top: $space-lg;
  padding-top: $space-md;
  border-top: 1px solid $color-border-subtle;
}

.extract-review-section__title {
  margin: 0 0 $space-xs;
  font-size: $font-size-xs;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: $color-text-secondary;
}

.extract-review-section__hint {
  margin: 0 0 $space-sm;
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.extract-review-section__actions {
  display: flex;
  flex-wrap: wrap;
  gap: $space-xs;
}

.extract-review-section__btn {
  font-size: $font-size-xs;
  padding: 4px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text;
  background: $color-bg-pane-elevated;
}

.extract-review-section__btn:disabled {
  opacity: 0.45;
}

.extract-review-section__btn:hover:not(:disabled) {
  border-color: $color-active-border;
}
</style>
