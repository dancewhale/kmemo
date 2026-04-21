<script setup lang="ts">
const props = defineProps<{
  label: string
  description?: string
}>()
</script>

<template>
  <div class="preference-field">
    <div class="preference-field__meta">
      <p class="preference-field__label">{{ props.label }}</p>
      <p v-if="props.description" class="preference-field__desc">{{ props.description }}</p>
    </div>
    <div class="preference-field__control">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.preference-field {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: $space-md;
  align-items: center;
}

.preference-field__meta {
  min-width: 0;
}

.preference-field__label {
  margin: 0;
  font-size: $font-size-sm;
  color: $color-text;
}

.preference-field__desc {
  margin: 2px 0 0;
  font-size: $font-size-xs;
  line-height: $line-normal;
  color: $color-text-secondary;
}

.preference-field__control {
  display: inline-flex;
  align-items: center;
  gap: $space-sm;
}
</style>
