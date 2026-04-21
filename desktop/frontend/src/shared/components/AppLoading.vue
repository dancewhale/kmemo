<script setup lang="ts">
withDefaults(
  defineProps<{
    label?: string
  }>(),
  {
    label: 'Loading…',
  },
)
</script>

<template>
  <div class="app-loading" role="status" :aria-label="label">
    <span class="app-loading__dot" />
    <span class="app-loading__dot app-loading__dot--2" />
    <span class="app-loading__dot app-loading__dot--3" />
    <span class="app-loading__label">{{ label }}</span>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-loading {
  display: inline-flex;
  align-items: center;
  gap: $space-sm;
  color: $color-text-secondary;
  font-size: $font-size-sm;
}

.app-loading__dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: $color-primary;
  opacity: 0.35;
  animation: app-loading-pulse 0.9s ease-in-out infinite;
}

.app-loading__dot--2 {
  animation-delay: 0.15s;
}

.app-loading__dot--3 {
  animation-delay: 0.3s;
}

.app-loading__label {
  margin-left: $space-xs;
}

@keyframes app-loading-pulse {
  0%,
  100% {
    opacity: 0.25;
    transform: scale(0.9);
  }
  50% {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
