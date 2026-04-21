<script setup lang="ts">
import SearchEmptyState from './SearchEmptyState.vue'
import SearchResultItem from './SearchResultItem.vue'
import type { SearchResultItem as ResultItem } from '../types'

const props = defineProps<{
  items: ResultItem[]
  selectedId: string | null
  hasQuery: boolean
}>()

const emit = defineEmits<{
  select: [id: string]
  open: [item: ResultItem]
}>()
</script>

<template>
  <div class="search-result-list">
    <SearchEmptyState v-if="!props.items.length" :has-query="props.hasQuery" />
    <div v-else class="search-result-list__scroll">
      <SearchResultItem
        v-for="item in props.items"
        :key="item.id"
        :item="item"
        :active="props.selectedId === item.id"
        @hover="emit('select', item.id)"
        @open="emit('open', item)"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.search-result-list {
  flex: 1 1 auto;
  min-height: 0;
  border-top: 1px solid $color-border-subtle;
}

.search-result-list__scroll {
  height: 100%;
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: $space-xs;
}
</style>
