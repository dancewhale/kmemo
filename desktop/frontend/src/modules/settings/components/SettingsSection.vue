<script setup lang="ts">
const props = defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <section class="settings-section">
    <header class="settings-section__head">
      <h2 class="settings-section__title">{{ props.title }}</h2>
      <p v-if="props.description" class="settings-section__desc">{{ props.description }}</p>
    </header>
    <div class="settings-section__body">
      <slot />
    </div>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.settings-section {
  border: 1px solid $color-border-subtle;
  border-radius: $radius-md;
  background: $color-bg-pane;
}

.settings-section__head {
  padding: $space-md;
  border-bottom: 1px solid $color-border-subtle;
}

.settings-section__title {
  margin: 0;
  font-size: $font-size-sm;
  font-weight: 600;
}

.settings-section__desc {
  margin: $space-xs 0 0;
  font-size: $font-size-xs;
  line-height: $line-normal;
  color: $color-text-secondary;
}

.settings-section__body {
  padding: $space-md;
  display: grid;
  gap: $space-md;
}
</style>
