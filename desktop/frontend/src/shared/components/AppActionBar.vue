<template>
  <div class="app-action-bar">
    <slot />
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-action-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: $space-sm;
  padding: $space-sm $space-md;
  border-bottom: 1px solid $color-border-subtle;
  background: $color-bg-pane;
}
</style>
