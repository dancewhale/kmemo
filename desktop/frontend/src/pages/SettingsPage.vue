<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppPane from '@/shared/components/AppPane.vue'
import { ROUTE_PATHS } from '@/shared/constants/routes'
</script>

<template>
  <AppPane title="Settings" class="settings-page">
    <p class="settings-page__p">Local preferences and host integration will live here.</p>
    <RouterLink class="settings-page__back" :to="ROUTE_PATHS.reading">← Back to workspace</RouterLink>
  </AppPane>
</template>

<style scoped lang="scss">
@use '@/app/styles/variables.scss' as *;

.settings-page {
  max-width: 560px;
}

.settings-page__p {
  margin: 0 0 $space-lg;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: $line-normal;
}

.settings-page__back {
  font-size: $font-size-sm;
  color: $color-primary;
}

.settings-page__back:hover {
  text-decoration: underline;
}
</style>
