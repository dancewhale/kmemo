<script setup lang="ts">
defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  createExtract: []
}>()
</script>

<template>
  <div v-if="visible" class="extract-toolbar">
    <button type="button" class="extract-toolbar__btn extract-toolbar__btn--primary" @click="emit('createExtract')">
      Create Extract
    </button>
    <button type="button" class="extract-toolbar__btn" disabled>Highlight</button>
    <button type="button" class="extract-toolbar__btn" disabled>Note</button>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.extract-toolbar {
  position: absolute;
  top: $space-sm;
  right: $space-sm;
  z-index: 5;
  display: inline-flex;
  align-items: center;
  gap: $space-xs;
  padding: $space-xs;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  background: $color-bg-pane-elevated;
}

.extract-toolbar__btn {
  font-size: $font-size-xs;
  padding: 2px $space-sm;
  border-radius: $radius-sm;
  border: 1px solid $color-border-subtle;
  color: $color-text-secondary;
  background: $color-bg-pane;
}

.extract-toolbar__btn--primary {
  color: $color-text;
  border-color: $color-active-border;
}

.extract-toolbar__btn:disabled {
  opacity: 0.45;
}
</style>
