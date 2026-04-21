<script setup lang="ts">
import AppEmpty from '@/shared/components/AppEmpty.vue'
import CommandListItem from './CommandListItem.vue'
import type { CommandMatchResult } from '../types'

const props = defineProps<{
  items: CommandMatchResult[]
  activeIndex: number
}>()

const emit = defineEmits<{
  select: [index: number]
  execute: [index: number]
}>()
</script>

<template>
  <div class="command-list">
    <AppEmpty v-if="!props.items.length" message="No commands found" />
    <div v-else class="command-list__scroll">
      <CommandListItem
        v-for="(item, idx) in props.items"
        :key="item.id"
        :command="item"
        :active="idx === props.activeIndex"
        @select="emit('select', idx)"
        @execute="emit('execute', idx)"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.command-list {
  min-height: 0;
  height: 100%;
}

.command-list__scroll {
  max-height: 360px;
  min-height: 160px;
  overflow-y: auto;
  padding: $space-xs;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
</style>
