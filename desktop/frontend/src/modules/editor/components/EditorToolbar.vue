<script setup lang="ts">
import AppActionBar from '@/shared/components/AppActionBar.vue'
import { useToast } from '@/shared/composables/useToast'

const tools = ['Bold', 'Italic', 'Highlight', 'Quote', 'Extract'] as const
const toast = useToast()

function onToolClick(name: string) {
  toast.info(`${name} command is a mock action`)
}
</script>

<template>
  <AppActionBar class="editor-toolbar">
    <button
      v-for="tool in tools"
      :key="tool"
      type="button"
      class="editor-toolbar__btn"
      @click="onToolClick(tool)"
    >
      {{ tool }}
    </button>
  </AppActionBar>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.editor-toolbar__btn {
  padding: 2px $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  font-size: $font-size-xs;
  color: $color-text-secondary;
  background: $color-bg-pane-elevated;
}

.editor-toolbar__btn:hover {
  color: $color-text;
  border-color: $color-active-border;
}
</style>
