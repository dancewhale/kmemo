<script setup lang="ts">
import { useObjectCreationStore } from '../stores/object-creation.store'
import type { CreateObjectType } from '../types'

const creation = useObjectCreationStore()

const tabs: { type: CreateObjectType; label: string }[] = [
  { type: 'node', label: 'Node' },
  { type: 'card', label: 'Card' },
  { type: 'manual-extract', label: 'Manual Extract' },
]

function select(type: CreateObjectType) {
  creation.setActiveType(type)
}
</script>

<template>
  <div class="create-object-type-tabs" role="tablist">
    <button
      v-for="t in tabs"
      :key="t.type"
      type="button"
      class="create-object-type-tabs__btn"
      :class="{ 'create-object-type-tabs__btn--active': creation.activeType === t.type }"
      role="tab"
      :aria-selected="creation.activeType === t.type"
      @click="select(t.type)"
    >
      {{ t.label }}
    </button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.create-object-type-tabs {
  display: flex;
  gap: 2px;
  padding: 2px;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  background: $color-bg-pane;
}

.create-object-type-tabs__btn {
  flex: 1 1 0;
  padding: 5px $space-sm;
  border: none;
  border-radius: calc(#{$radius-sm} - 2px);
  font-size: $font-size-xs;
  font-weight: 500;
  color: $color-text-secondary;
  background: transparent;
  cursor: default;
}

.create-object-type-tabs__btn:hover {
  color: $color-text;
  background: $color-hover;
}

.create-object-type-tabs__btn--active {
  color: $color-text;
  background: $color-bg-pane-elevated;
  border: 1px solid $color-border-subtle;
}
</style>
