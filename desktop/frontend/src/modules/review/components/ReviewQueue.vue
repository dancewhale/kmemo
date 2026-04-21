<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import AppEmpty from '@/shared/components/AppEmpty.vue'
import AppLoading from '@/shared/components/AppLoading.vue'
import { useReviewStore } from '../stores/review.store'
import ReviewStatsBar from './ReviewStatsBar.vue'
import ReviewQueueItem from './ReviewQueueItem.vue'

const review = useReviewStore()
const { queueItems, selectedReviewItemId, loading, stats, isQueueEmpty } = storeToRefs(review)

onMounted(async () => {
  await review.initialize()
})
</script>

<template>
  <div class="review-queue">
    <ReviewStatsBar :stats="stats" />
    <div class="review-queue__list">
      <div v-if="loading" class="review-queue__loading">
        <AppLoading label="Loading review queue…" />
      </div>
      <AppEmpty v-else-if="isQueueEmpty" message="All reviews completed for now." />
      <div v-else>
        <ReviewQueueItem
          v-for="item in queueItems"
          :key="item.id"
          :item="item"
          :selected="selectedReviewItemId === item.id"
          @select="review.setSelectedItem($event)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.review-queue {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: $color-bg-pane;
}

.review-queue__list {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  padding: $space-xs 0 $space-sm;
}

.review-queue__loading {
  display: flex;
  justify-content: center;
  padding: $space-xl 0;
}

.review-queue__list :deep(.app-empty) {
  padding: $space-xl $space-md;
}
</style>
