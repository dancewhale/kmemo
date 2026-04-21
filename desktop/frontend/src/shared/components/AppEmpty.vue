<script setup lang="ts">
withDefaults(
  defineProps<{
    message?: string
  }>(),
  {
    message: 'Nothing here yet',
  },
)
</script>

<template>
  <div class="app-empty">
    <p class="app-empty__msg">{{ message }}</p>
    <div v-if="$slots.default" class="app-empty__extra">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-empty {
  padding: $space-xl;
  text-align: center;
  color: $color-text-secondary;
  font-size: $font-size-sm;
}

.app-empty__msg {
  margin: 0;
}

.app-empty__extra {
  margin-top: $space-md;
}
</style>
