<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    kind?: 'neutral' | 'info' | 'success' | 'warning' | 'danger'
  }>(),
  {
    kind: 'neutral',
  },
)
</script>

<template>
  <span class="app-status-tag" :class="`app-status-tag--${props.kind}`">
    <slot />
  </span>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-status-tag {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: $font-size-caption;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  border: 1px solid $color-border;
  border-radius: $radius-sm;
  padding: 1px 6px;
  color: $color-text-secondary;
}

.app-status-tag--info {
  border-color: rgba(107, 155, 209, 0.5);
  color: $color-info;
}

.app-status-tag--success {
  border-color: rgba(122, 180, 132, 0.5);
  color: $color-success;
}

.app-status-tag--warning {
  border-color: rgba(201, 184, 122, 0.55);
  color: $color-warning;
}

.app-status-tag--danger {
  border-color: rgba(195, 120, 120, 0.55);
  color: $color-danger;
}
</style>
