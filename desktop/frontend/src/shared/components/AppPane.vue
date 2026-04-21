<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    title?: string
    bordered?: boolean
    scrollable?: boolean
    padded?: boolean | 'none' | 'sm' | 'md'
    compact?: boolean
  }>(),
  {
    bordered: true,
    scrollable: true,
    padded: 'md',
    compact: false,
  },
)

const padClass = computed(() => {
  if (props.padded === false || props.padded === 'none') {
    return 'app-pane__body--pad-none'
  }
  if (props.padded === 'sm') {
    return 'app-pane__body--pad-sm'
  }
  return 'app-pane__body--pad-md'
})
</script>

<template>
  <section
    class="app-pane"
    :class="{
      'app-pane--bordered': bordered,
      'app-pane--scroll': scrollable,
      'app-pane--compact': compact,
    }"
  >
    <header v-if="title || $slots.header" class="app-pane__header">
      <slot name="header">
        <span v-if="title" class="app-pane__title">{{ title }}</span>
      </slot>
    </header>
    <div class="app-pane__body" :class="padClass">
      <slot />
    </div>
    <footer v-if="$slots.footer" class="app-pane__footer">
      <slot name="footer" />
    </footer>
  </section>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-pane {
  display: flex;
  flex-direction: column;
  min-height: 0;
  min-width: 0;
  background: $color-bg-pane;
  border-radius: $radius-sm;
}

.app-pane--bordered {
  border: 1px solid $color-border-subtle;
}

.app-pane--scroll .app-pane__body {
  overflow: auto;
}

.app-pane__header {
  flex: 0 0 auto;
  padding: $space-sm $space-md;
  border-bottom: 1px solid $color-border-subtle;
  font-size: $font-size-sm;
  font-weight: 600;
  color: $color-text-secondary;
  letter-spacing: 0.02em;
}

.app-pane--compact .app-pane__header {
  padding: $space-xs $space-sm;
}

.app-pane__title {
  color: $color-text;
}

.app-pane__body {
  flex: 1 1 auto;
  min-height: 0;
}

.app-pane__body--pad-none {
  padding: 0;
}

.app-pane__body--pad-sm {
  padding: $space-sm;
}

.app-pane__body--pad-md {
  padding: $space-md;
}

.app-pane__footer {
  flex: 0 0 auto;
  padding: $space-sm $space-md;
  border-top: 1px solid $color-border-subtle;
}
</style>
