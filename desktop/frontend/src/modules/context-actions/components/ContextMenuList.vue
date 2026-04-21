<script setup lang="ts">
import type { ContextActionItem } from '../types'
import ContextMenuItem from './ContextMenuItem.vue'

defineProps<{
  items: ContextActionItem[]
}>()

const emit = defineEmits<{
  select: [id: string]
}>()
</script>

<template>
  <div class="context-menu-list" role="menu">
    <ContextMenuItem
      v-for="item in items"
      :key="item.id"
      :item="item"
      @select="emit('select', $event)"
    />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.context-menu-list {
  min-width: 200px;
  max-width: 280px;
  padding: $space-xs;
}
</style>
