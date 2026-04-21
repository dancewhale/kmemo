<script setup lang="ts">
import AppHeader from '@/shared/components/AppHeader.vue'
</script>

<template>
  <div class="main-layout">
    <div class="main-layout__chrome">
      <AppHeader title="kmemo" subtitle="Progressive reading shell" />
    </div>
    <div class="main-layout__body">
      <router-view />
    </div>
  </div>
</template>
