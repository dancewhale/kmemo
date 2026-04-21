<script setup lang="ts">
import type { ThemeMode } from '../types'

const props = defineProps<{
  modelValue: ThemeMode
}>()

const emit = defineEmits<{
  'update:modelValue': [value: ThemeMode]
}>()

const modes: Array<{ id: ThemeMode; label: string }> = [
  { id: 'light', label: 'Light' },
  { id: 'dark', label: 'Dark' },
  { id: 'system', label: 'System' },
]
</script>

<template>
  <div class="theme-mode-select">
    <button
      v-for="mode in modes"
      :key="mode.id"
      type="button"
      class="theme-mode-select__btn"
      :class="{ 'theme-mode-select__btn--active': props.modelValue === mode.id }"
      @click="emit('update:modelValue', mode.id)"
    >
      {{ mode.label }}
    </button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.theme-mode-select {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.theme-mode-select__btn {
  border: 1px solid $color-border;
  border-radius: $radius-sm;
  background: transparent;
  color: $color-text-secondary;
  font-size: $font-size-xs;
  padding: 4px 10px;
}

.theme-mode-select__btn:hover {
  border-color: $color-active-border;
  color: $color-text;
}

.theme-mode-select__btn--active {
  border-color: $color-active-border;
  background: $color-active;
  color: $color-text;
}
</style>
