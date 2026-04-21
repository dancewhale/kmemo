<script setup lang="ts">
import AppStatusTag from '@/shared/components/AppStatusTag.vue'
import { useContextMenuStore } from '@/modules/context-actions/stores/context-menu.store'
import type { SearchResultItem } from '../types'

const contextMenu = useContextMenuStore()

const props = defineProps<{
  item: SearchResultItem
  active: boolean
}>()

const emit = defineEmits<{
  open: []
  hover: []
}>()

function formatDate(iso?: string | null): string {
  if (!iso) {
    return ''
  }
  try {
    return new Date(iso).toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return iso
  }
}

function onContextMenu(e: MouseEvent) {
  e.preventDefault()
  contextMenu.open({
    entityType: 'search-result',
    entityId: props.item.id,
    searchResultId: props.item.id,
    x: e.clientX,
    y: e.clientY,
  })
}
</script>

<template>
  <button
    type="button"
    class="search-result-item"
    :class="{ 'search-result-item--active': props.active }"
    @mouseenter="emit('hover')"
    @click="emit('open')"
    @contextmenu="onContextMenu"
  >
    <div class="search-result-item__top">
      <span class="search-result-item__title">{{ props.item.title }}</span>
      <AppStatusTag kind="neutral">{{ props.item.entityType }}</AppStatusTag>
    </div>
    <p class="search-result-item__snippet">{{ props.item.snippet }}</p>
    <div class="search-result-item__meta">
      <span v-if="props.item.sourceTitle">{{ props.item.sourceTitle }}</span>
      <span v-if="props.item.subtitle">{{ props.item.subtitle }}</span>
      <span v-if="props.item.updatedAt">{{ formatDate(props.item.updatedAt) }}</span>
    </div>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.search-result-item {
  width: 100%;
  text-align: left;
  border: 1px solid transparent;
  border-radius: $radius-sm;
  background: transparent;
  color: $color-text;
  padding: $space-sm $space-md;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.search-result-item:hover {
  background: $color-hover;
}

.search-result-item--active {
  border-color: $color-active-border;
  background: $color-active;
}

.search-result-item__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
}

.search-result-item__title {
  font-size: $font-size-sm;
  line-height: $line-tight;
  font-weight: 600;
}

.search-result-item__snippet {
  margin: 0;
  font-size: $font-size-xs;
  line-height: $line-normal;
  color: $color-text-secondary;
}

.search-result-item__meta {
  display: flex;
  align-items: center;
  gap: $space-sm;
  font-size: 10px;
  color: $color-text-secondary;
}
</style>
