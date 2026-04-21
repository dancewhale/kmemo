<script setup lang="ts">
const props = defineProps<{
  keys: string[]
}>()
</script>

<template>
  <span class="shortcut-hint">
    <kbd v-for="(key, idx) in props.keys" :key="`${key}-${idx}`" class="shortcut-hint__key">
      {{ key }}
    </kbd>
  </span>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.shortcut-hint {
  display: inline-flex;
  gap: 4px;
  align-items: center;
}

.shortcut-hint__key {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: $radius-sm;
  border: 1px solid $color-border;
  background: $color-bg-pane;
  color: $color-text-secondary;
  font-family: $font-mono;
  font-size: 10px;
  line-height: 1;
}
</style>
