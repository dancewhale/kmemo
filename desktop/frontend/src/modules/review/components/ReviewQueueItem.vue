<script setup lang="ts">
import AppStatusTag from '@/shared/components/AppStatusTag.vue'
import { useContextMenuStore } from '@/modules/context-actions/stores/context-menu.store'
import type { ReviewItem } from '../types'

const contextMenu = useContextMenuStore()

defineProps<{
  item: ReviewItem
  selected: boolean
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

function typeKind(type: ReviewItem['type']): 'info' | 'warning' {
  return type === 'extract' ? 'warning' : 'info'
}

function onContextMenu(e: MouseEvent, row: ReviewItem) {
  e.preventDefault()
  contextMenu.open({
    entityType: 'review-item',
    entityId: row.id,
    reviewId: row.id,
    x: e.clientX,
    y: e.clientY,
  })
}
</script>

<template>
  <button
    type="button"
    class="review-queue-item"
    :class="{ 'review-queue-item--active': selected }"
    @click="emit('select', item.id)"
    @contextmenu="onContextMenu($event, item)"
  >
    <div class="review-queue-item__top">
      <span class="review-queue-item__title">{{ item.title }}</span>
      <AppStatusTag :kind="typeKind(item.type)">
        {{ item.type }} · {{ item.status }}
      </AppStatusTag>
    </div>
    <p class="review-queue-item__prompt">{{ item.prompt }}</p>
    <div class="review-queue-item__meta kmono">
      <span>due {{ item.dueAt }}</span>
      <span>· reviews {{ item.reviewCount }}</span>
      <span v-if="item.sourceArticleTitle">· {{ item.sourceArticleTitle }}</span>
    </div>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.review-queue-item {
  width: calc(100% - 4px);
  margin: 0 2px 1px;
  padding: $space-sm $space-md;
  border: 1px solid transparent;
  border-radius: $radius-sm;
  text-align: left;
  color: $color-text;
}

.review-queue-item:hover {
  background: $color-hover;
}

.review-queue-item--active {
  background: $color-active;
  border-color: $color-active-border;
}

.review-queue-item__top {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: $space-sm;
}

.review-queue-item__title {
  font-size: $font-size-sm;
  font-weight: 600;
}

.review-queue-item__prompt {
  margin: 3px 0;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.review-queue-item__meta {
  font-size: 10px;
  color: $color-text-secondary;
}
</style>
