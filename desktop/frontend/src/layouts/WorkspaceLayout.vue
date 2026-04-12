<template>
  <AppShell>
    <section class="workspace-layout">
      <aside class="pane pane-sidebar">
        <slot name="sidebar" />
      </aside>
      <section class="pane pane-main">
        <slot name="main" />
      </section>
      <aside class="pane pane-detail">
        <slot name="detail" />
      </aside>
    </section>
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from "./AppShell.vue";
</script>
