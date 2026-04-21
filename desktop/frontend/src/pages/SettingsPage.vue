<script setup lang="ts">
import { onMounted } from 'vue'
import SettingsPanel from '@/modules/settings/components/SettingsPanel.vue'
import { useSettingsStore } from '@/modules/settings/stores/settings.store'

const settings = useSettingsStore()

onMounted(() => {
  settings.initialize()
})
</script>

<template>
  <SettingsPanel />
</template>
