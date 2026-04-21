<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    title?: string
    description?: string
    message?: string
    type?: 'default' | 'search' | 'list'
  }>(),
  {
    title: '',
    description: '',
    message: 'Nothing here yet',
    type: 'default',
  },
)
</script>

<template>
  <div class="app-empty" :class="`app-empty--${props.type}`">
    <div v-if="$slots.icon" class="app-empty__icon">
      <slot name="icon" />
    </div>
    <p class="app-empty__title">{{ props.title || props.message }}</p>
    <p v-if="props.description" class="app-empty__desc">{{ props.description }}</p>
    <div v-if="$slots.default" class="app-empty__extra">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-empty {
  padding: $space-xl $space-md;
  text-align: center;
  color: $color-text-muted;
  font-size: $font-size-xs;
  display: grid;
  gap: $space-xs;
  justify-items: center;
}

.app-empty__icon {
  color: $color-text-secondary;
}

.app-empty__title {
  margin: 0;
  color: $color-text-secondary;
  font-size: $font-size-sm;
}

.app-empty__desc {
  margin: 0;
  max-width: 480px;
  line-height: $line-normal;
}

.app-empty__extra {
  margin-top: $space-md;
}
</style>
