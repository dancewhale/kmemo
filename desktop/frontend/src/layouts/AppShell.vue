<template>
  <div class="app-shell">
    <header class="topbar">
      <div class="brand">kmemo</div>
      <nav class="topnav">
        <RouterLink to="/workspace">Workspace</RouterLink>
        <RouterLink to="/review">Review</RouterLink>
        <RouterLink to="/settings">Settings</RouterLink>
      </nav>
    </header>
    <main class="shell-content">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
</script>
