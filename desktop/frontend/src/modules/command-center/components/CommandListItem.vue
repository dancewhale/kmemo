<script setup lang="ts">
import { computed } from 'vue'
import { useSettingsStore } from '@/modules/settings/stores/settings.store'
import ShortcutHint from './ShortcutHint.vue'
import type { CommandMatchResult } from '../types'

const props = defineProps<{
  command: CommandMatchResult
  active: boolean
}>()

const emit = defineEmits<{
  select: []
  execute: []
}>()

const settings = useSettingsStore()
const showGroup = computed(() => settings.workspacePreferences.commandPaletteShowGroup)
const showShortcutHints = computed(() => settings.workspacePreferences.showShortcutHints)
</script>

<template>
  <button
    type="button"
    class="command-list-item"
    :disabled="props.command.enabled === false"
    :class="{
      'command-list-item--active': props.active,
      'command-list-item--disabled': props.command.enabled === false,
    }"
    @mouseenter="emit('select')"
    @click="emit('execute')"
  >
    <span class="command-list-item__main">
      <span class="command-list-item__title">{{ props.command.title }}</span>
      <span v-if="props.command.subtitle" class="command-list-item__subtitle">
        {{ props.command.subtitle }}
      </span>
    </span>
    <span class="command-list-item__meta">
      <span v-if="showGroup" class="command-list-item__group">{{ props.command.group }}</span>
      <ShortcutHint
        v-if="showShortcutHints && props.command.shortcut?.length"
        :keys="props.command.shortcut"
      />
    </span>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.command-list-item {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-sm;
  border: 1px solid transparent;
  border-radius: $radius-sm;
  background: transparent;
  color: $color-text;
  text-align: left;
  padding: $space-sm $space-md;
}

.command-list-item:hover {
  background: $color-hover;
}

.command-list-item--active {
  border-color: $color-active-border;
  background: $color-active;
}

.command-list-item--disabled {
  opacity: 0.62;
  cursor: not-allowed;
}

.command-list-item__main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.command-list-item__title {
  font-size: $font-size-sm;
  line-height: $line-tight;
  font-weight: 600;
}

.command-list-item__subtitle {
  font-size: $font-size-xs;
  line-height: $line-normal;
  color: $color-text-secondary;
}

.command-list-item__meta {
  display: inline-flex;
  align-items: center;
  gap: $space-sm;
  flex-shrink: 0;
}

.command-list-item__group {
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-size: 10px;
  color: $color-text-secondary;
}
</style>
