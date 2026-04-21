<script setup lang="ts">
import { computed } from 'vue'
import {
  ArrowLeft,
  CircleCheck,
  Clock,
  Close,
  Collection,
  Delete,
  Document,
  DocumentCopy,
  EditPen,
  FolderOpened,
  Notebook,
  Plus,
  Postcard,
  Reading,
  Right,
  View,
} from '@element-plus/icons-vue'
import type { ContextActionItem } from '../types'

const props = defineProps<{
  item: ContextActionItem
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

const ICON_MAP: Record<string, object> = {
  ArrowLeft,
  CircleCheck,
  Clock,
  Close,
  Collection,
  Delete,
  Document,
  DocumentCopy,
  EditPen,
  FolderOpened,
  Notebook,
  Plus,
  Postcard,
  Reading,
  Right,
  View,
}

const iconComp = computed(() => {
  const name = props.item.icon
  return name ? (ICON_MAP[name] ?? null) : null
})
</script>

<template>
  <div v-if="item.dividerBefore" class="context-menu-item__divider" role="separator" />
  <button
    type="button"
    class="context-menu-item"
    :class="{ 'context-menu-item--danger': item.danger }"
    :disabled="item.disabled"
    @click="emit('select', item.id)"
  >
    <span v-if="iconComp" class="context-menu-item__ico">
      <el-icon><component :is="iconComp" /></el-icon>
    </span>
    <span class="context-menu-item__label">{{ item.label }}</span>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.context-menu-item__divider {
  height: 1px;
  margin: $space-xs 0;
  background: $color-border-subtle;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: $space-xs;
  width: 100%;
  padding: 6px $space-sm;
  border: none;
  border-radius: $radius-sm;
  background: transparent;
  color: $color-text;
  font-size: $font-size-sm;
  text-align: left;
  cursor: default;
}

.context-menu-item:hover:not(:disabled) {
  background: $color-hover;
}

.context-menu-item:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.context-menu-item--danger {
  color: #e57373;
}

.context-menu-item__ico {
  flex: 0 0 auto;
  display: flex;
  font-size: 14px;
  color: $color-text-secondary;
}

.context-menu-item__label {
  flex: 1 1 auto;
  min-width: 0;
}
</style>
