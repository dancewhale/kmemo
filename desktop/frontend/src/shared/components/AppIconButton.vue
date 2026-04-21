<script setup lang="ts">
withDefaults(
  defineProps<{
    label: string
    active?: boolean
  }>(),
  {
    active: false,
  },
)
</script>

<template>
  <button type="button" class="app-icon-btn" :class="{ 'app-icon-btn--active': active }" :title="label">
    <span class="app-icon-btn__inner">
      <slot />
    </span>
  </button>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: $radius-sm;
  color: $color-text-secondary;
}

.app-icon-btn:hover {
  background: $color-hover;
  color: $color-text;
}

.app-icon-btn--active {
  background: $color-active;
  color: $color-primary;
}

.app-icon-btn__inner {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}
</style>
