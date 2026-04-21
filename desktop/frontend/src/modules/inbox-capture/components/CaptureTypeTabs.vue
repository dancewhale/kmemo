<script setup lang="ts">
import { useCaptureStore } from '../stores/capture.store'
import type { CaptureType } from '../types'

const capture = useCaptureStore()

const tabs: { type: CaptureType; label: string }[] = [
  { type: 'blank', label: 'Blank' },
  { type: 'paste-text', label: 'Paste Text' },
  { type: 'url', label: 'URL' },
]

function select(type: CaptureType) {
  capture.setActiveType(type)
}
</script>

<template>
  <div class="capture-type-tabs" role="tablist">
    <button
      v-for="t in tabs"
      :key="t.type"
      type="button"
      class="capture-type-tabs__btn"
      :class="{ 'capture-type-tabs__btn--active': capture.activeType === t.type }"
      role="tab"
      :aria-selected="capture.activeType === t.type"
      @click="select(t.type)"
    >
      {{ t.label }}
    </button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.capture-type-tabs {
  display: flex;
  gap: 2px;
  padding: 2px;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  background: $color-bg-pane;
}

.capture-type-tabs__btn {
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

.capture-type-tabs__btn:hover {
  color: $color-text;
  background: $color-hover;
}

.capture-type-tabs__btn--active {
  color: $color-text;
  background: $color-bg-pane-elevated;
  border: 1px solid $color-border-subtle;
}
</style>
