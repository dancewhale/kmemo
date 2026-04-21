<script setup lang="ts">
import type { SearchFilterType } from '../types'

const props = defineProps<{
  modelValue: SearchFilterType
}>()

const emit = defineEmits<{
  'update:modelValue': [value: SearchFilterType]
}>()

const options: Array<{ id: SearchFilterType; label: string }> = [
  { id: 'all', label: 'All' },
  { id: 'article', label: 'Articles' },
  { id: 'extract', label: 'Extracts' },
  { id: 'card', label: 'Cards' },
  { id: 'node', label: 'Nodes' },
  { id: 'review', label: 'Reviews' },
]
</script>

<template>
  <div class="search-filter-tabs">
    <button
      v-for="opt in options"
      :key="opt.id"
      type="button"
      class="search-filter-tabs__item"
      :class="{ 'search-filter-tabs__item--active': props.modelValue === opt.id }"
      @click="emit('update:modelValue', opt.id)"
    >
      {{ opt.label }}
    </button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.search-filter-tabs {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.search-filter-tabs__item {
  padding: 4px 9px;
  border: 1px solid $color-border;
  border-radius: $radius-sm;
  background: transparent;
  color: $color-text-secondary;
  font-size: $font-size-xs;
  line-height: 1.2;
}

.search-filter-tabs__item:hover {
  border-color: $color-active-border;
  color: $color-text;
}

.search-filter-tabs__item--active {
  color: $color-text;
  border-color: $color-active-border;
  background: $color-active;
}
</style>
