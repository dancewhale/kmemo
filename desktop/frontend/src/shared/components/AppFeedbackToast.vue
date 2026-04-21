<script setup lang="ts">
import { useToast } from '@/shared/composables/useToast'

const { toasts, removeToast } = useToast()
</script>

<template>
  <Teleport to="body">
    <div class="app-feedback-toast">
      <button
        v-for="toast in toasts"
        :key="toast.id"
        type="button"
        class="app-feedback-toast__item"
        :class="`app-feedback-toast__item--${toast.kind}`"
        @click="removeToast(toast.id)"
      >
        {{ toast.message }}
      </button>
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.app-feedback-toast {
  position: fixed;
  right: $space-lg;
  bottom: $space-lg;
  display: flex;
  flex-direction: column;
  gap: $space-sm;
  z-index: 1500;
}

.app-feedback-toast__item {
  min-width: 240px;
  max-width: 440px;
  text-align: left;
  border-radius: $radius-md;
  border: 1px solid $color-border;
  background: $color-bg-pane-elevated;
  color: $color-text;
  padding: $space-sm $space-md;
  font-size: $font-size-sm;
  box-shadow: $shadow-overlay;
}

.app-feedback-toast__item--success {
  border-color: color-mix(in srgb, $color-success 48%, transparent);
}

.app-feedback-toast__item--info {
  border-color: color-mix(in srgb, $color-info 48%, transparent);
}

.app-feedback-toast__item--warning {
  border-color: color-mix(in srgb, $color-warning 48%, transparent);
}

.app-feedback-toast__item--error {
  border-color: color-mix(in srgb, $color-danger 52%, transparent);
}
</style>
