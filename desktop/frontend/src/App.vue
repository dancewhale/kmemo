<script setup lang="ts">
</script>

<template>
  <div class="dark">
    <router-view />
  </div>
</template>

<style scoped>
.dark {
  height: 100%;
  min-height: 0;
  color-scheme: dark;
}
</style>
