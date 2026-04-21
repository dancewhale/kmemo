<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { useObjectCreationStore } from '../stores/object-creation.store'
import type { CreateObjectType } from '../types'

const props = withDefaults(
  defineProps<{
    defaultType?: CreateObjectType
    compact?: boolean
  }>(),
  {
    defaultType: 'node',
    compact: false,
  },
)

const creation = useObjectCreationStore()

function onClick() {
  creation.openDialog(props.defaultType)
}
</script>

<template>
  <button
    type="button"
    class="create-quick-entry"
    :class="{ 'create-quick-entry--compact': compact }"
    title="Create knowledge object"
    @click="onClick"
  >
    <el-icon v-if="compact" class="create-quick-entry__ico"><Plus /></el-icon>
    <span v-if="!compact" class="create-quick-entry__label">New</span>
    <span v-else class="create-quick-entry__sr">New</span>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.create-quick-entry {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: $space-xs;
  padding: 5px $space-md;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  font-size: $font-size-xs;
  font-weight: 600;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
  cursor: default;
}

.create-quick-entry:hover {
  border-color: $color-active-border;
  color: $color-text;
}

.create-quick-entry--compact {
  padding: 4px 8px;
  min-width: 30px;
}

.create-quick-entry__ico {
  font-size: 15px;
}

.create-quick-entry__label {
  white-space: nowrap;
}

.create-quick-entry__sr {
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
