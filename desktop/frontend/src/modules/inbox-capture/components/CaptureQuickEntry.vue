<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { useCaptureStore } from '../stores/capture.store'
import type { CaptureType } from '../types'

const props = withDefaults(
  defineProps<{
    defaultType?: CaptureType
    compact?: boolean
  }>(),
  {
    defaultType: 'blank',
    compact: false,
  },
)

const capture = useCaptureStore()

function onClick() {
  capture.openDialog(props.defaultType)
}
</script>

<template>
  <button
    type="button"
    class="capture-quick-entry"
    :class="{ 'capture-quick-entry--compact': compact }"
    title="Add to Inbox"
    @click="onClick"
  >
    <el-icon v-if="compact" class="capture-quick-entry__ico"><Plus /></el-icon>
    <span v-if="!compact" class="capture-quick-entry__label">Add to Inbox</span>
    <span v-else class="capture-quick-entry__sr">Add to Inbox</span>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.capture-quick-entry {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: $space-xs;
  padding: 5px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-active-border;
  font-size: $font-size-xs;
  font-weight: 600;
  color: $color-text;
  background: $color-bg-pane-elevated;
  cursor: default;
}

.capture-quick-entry:hover {
  background: $color-hover;
}

.capture-quick-entry--compact {
  padding: 4px 8px;
  min-width: 30px;
}

.capture-quick-entry__ico {
  font-size: 15px;
}

.capture-quick-entry__label {
  white-space: nowrap;
}

.capture-quick-entry__sr {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
