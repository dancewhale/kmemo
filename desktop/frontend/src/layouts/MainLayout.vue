<script setup lang="ts">
import { onMounted } from 'vue'
import AppHeader from '@/shared/components/AppHeader.vue'
import AppFeedbackToast from '@/shared/components/AppFeedbackToast.vue'
import { setupShortcuts } from '@/app/providers/shortcuts'
import CommandPalette from '@/modules/command-center/components/CommandPalette.vue'
import ContextMenu from '@/modules/context-actions/components/ContextMenu.vue'
import CaptureDialog from '@/modules/inbox-capture/components/CaptureDialog.vue'
import CreateObjectDialog from '@/modules/object-creation/components/CreateObjectDialog.vue'

onMounted(() => {
  setupShortcuts()
})
</script>

<template>
  <div class="main-layout">
    <div class="main-layout__chrome">
      <AppHeader title="kmemo" subtitle="Progressive reading shell" />
    </div>
    <div class="main-layout__body">
      <router-view />
    </div>
    <CommandPalette />
    <ContextMenu />
    <CaptureDialog />
    <CreateObjectDialog />
    <AppFeedbackToast />
  </div>
</template>
