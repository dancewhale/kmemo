<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <h3 class="app-section-title">{{ title }}</h3>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-section-title {
  margin: 0 0 $space-xs;
  font-size: $font-size-caption;
  color: $color-text-muted;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}
</style>
