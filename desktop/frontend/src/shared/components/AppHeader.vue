<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <header class="app-header">
    <div class="app-header__titles">
      <span class="app-header__title">{{ title }}</span>
      <span v-if="subtitle" class="app-header__subtitle">{{ subtitle }}</span>
    </div>
    <div v-if="$slots.actions" class="app-header__actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $space-md;
  height: $header-height;
  padding: 0 $space-lg;
}

.app-header__titles {
  display: flex;
  align-items: baseline;
  gap: $space-md;
  min-width: 0;
}

.app-header__title {
  font-size: $font-size-lg;
  font-weight: 600;
  color: $color-text;
  white-space: nowrap;
}

.app-header__subtitle {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-header__actions {
  display: flex;
  align-items: center;
  gap: $space-xs;
  flex-shrink: 0;
}
</style>
